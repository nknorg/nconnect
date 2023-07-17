package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/arch"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nkn-sdk-go"
	tunnel "github.com/nknorg/nkn-tunnel"
)

const (
	memberFile = "member.json"
)

var (
	errNoDataInFile = "no data in file"
	errWaitForAuth  = "wait for authorization"
	errNameExist    = "network node name already exists"
	errNodeNotFound = "node not found"
)

type memberNetworkData struct {
	NetworkInfo     *networkInfo `json:"networkInfo"`
	NodeInfo        *NodeInfo    `json:"NodeInfo"`
	NodesIAccept    []*NodeInfo  `json:"nodesIAccept"`    // node info of nodes that this node accepts
	NodesICanAccess []*NodeInfo  `json:"nodesICanAccess"` // node info of nodes that this node can access
}

type callbackNodeICanAccessUpdated func(nodes []*NodeInfo) error

type Member struct {
	opts                    *config.Opts
	c                       *admin.Client
	networkData             memberNetworkData // node info of this node
	serverAddress           string            // nconnect server tunnel address
	serverTunnel            *tunnel.Tunnel
	joinedNetwork           bool
	CbNodeICanAccessUpdated callbackNodeICanAccessUpdated
	openTunOnce             sync.Once // only open tun device once
}

func NewMember(opts *config.Opts, c *admin.Client) *Member {
	return &Member{opts: opts, c: c, networkData: memberNetworkData{NetworkInfo: &networkInfo{}, NodeInfo: &NodeInfo{}}}
}

func (m *Member) StartMember(serverAddress string) error {
	m.serverAddress = serverAddress

	err := m.loadMemberData()
	if err != nil && err.Error() != errNoDataInFile {
		return err
	}

	if err = m.JoinNetwork(serverAddress); err != nil {
		return err
	}

	if m.joinedNetwork {
		if err = m.GetNodeICanAccess(); err != nil {
			return err
		}
		if err = m.GetNodeIAccept(); err != nil {
			return err
		}
	}

	log.Println("nConnect Network member is listening at:", m.c.Address())
	for {
		msg := <-m.c.OnMessage.C

		req := &managerToMember{}
		err := json.Unmarshal(msg.Data, req)
		if err != nil {
			log.Println("Network member, received multiclient msg, unmarshal msg.Data error: ", err)
			continue
		}
		if m.opts.Verbose {
			log.Printf("Network member, received multiclient msg: %+v\n", req)
		}

		go func() {
			err = m.handleNknMsg(req)
			if err != nil {
				log.Println(err)
			}
			if req.MsgType == NKN_PING {
				resp := req
				resp.MsgType = NKN_PONG
				b, err := json.Marshal(resp)
				if err != nil {
					log.Println("Network member, marshal pong error: ", err)
					return
				}

				err = msg.Reply(b)
				if err != nil {
					log.Println("Network member, reply pong error: ", err)
				}
			}
		}()
	}
}

// handle notification from manager
func (m *Member) handleNknMsg(notification *managerToMember) error {
	switch notification.MsgType {
	case NOTI_AUTHORIZED: // I was authorized by manager
		if len(notification.NodeInfo) > 0 {
			m.networkData.NodeInfo = notification.NodeInfo[0]
			if m.serverAddress != "" && m.networkData.NodeInfo.ServerAddress != m.serverAddress {
				err := m.SetServerTunnel(m.serverTunnel)
				if err != nil {
					return err
				}
				m.networkData.NodeInfo.ServerAddress = m.serverAddress
			}
			if err := m.saveMemberData(); err != nil {
				return err
			}

			log.Printf("\n\nCongratulations!!! Your nConnect network member is authorized, IP: %v, mask: %v\n\n",
				m.networkData.NodeInfo.IP, m.networkData.NodeInfo.Netmask)

			m.OpenTunAndSetIp()
			m.GetNodeICanAccess()
		}

	case NOTI_NEW_MEMBER: // new member is authorized and joined the network
		if len(notification.NodeInfo) > 0 {
			m.networkData.NodesIAccept = append(m.networkData.NodesIAccept, notification.NodeInfo...)
			if err := m.saveMemberData(); err != nil {
				return err
			}
			m.UpdMyAccept(notification.NodeInfo)
		}

	case NOTI_UPD_I_ACCEPT:
		m.GetNodeIAccept()

	case NOTI_MEMBER_ONLINE:
		m.GetNodeICanAccess()
		if m.CbNodeICanAccessUpdated != nil {
			m.CbNodeICanAccessUpdated(m.networkData.NodesICanAccess)
		}
		m.UpdMyAccept(notification.NodeInfo)

	case NOTI_UPD_I_CAN_ACCESS:
		m.GetNodeICanAccess()
		if m.CbNodeICanAccessUpdated != nil {
			m.CbNodeICanAccessUpdated(m.networkData.NodesICanAccess)
		}

	case NKN_PING:
		log.Println("Network member, received ping from manager, send pong back")

	default:
		return fmt.Errorf("nConnect member got unknown notification type: %v", notification.MsgType)
	}

	return nil
}

func (m *Member) JoinNetwork(serverAddr string) error {
	if serverAddr == "" {
		serverAddr = m.networkData.NodeInfo.ServerAddress
	} else {
		if m.networkData.NodeInfo.ServerAddress != serverAddr {
			m.networkData.NodeInfo.ServerAddress = serverAddr
			m.saveMemberData()
		}
	}

	msg := memberToManager{MsgType: JOIN_NETWORK, Name: m.opts.NodeName, ServerAddress: serverAddr}
	resp, err := SendMsg(m.c, m.opts.ManagerAddress, &msg, true)
	if err != nil {
		return err
	}

	if resp.Err == errWaitForAuth {
		m.networkData.NetworkInfo = resp.NetworkInfo
		m.saveMemberData()

		log.Println("You sent a join the network request to the manager, wait for the manager to authorize.")

		return nil
	} else if resp.Err == errNameExist {
		log.Println("You network node name is used by other node, please config another name")
		return errors.New(errNameExist)
	}

	if resp.Err != "" {
		return errors.New(resp.Err)
	}

	if len(resp.NodeInfo) > 0 {
		m.networkData.NodeInfo = resp.NodeInfo[0]
		m.networkData.NetworkInfo = resp.NetworkInfo
		m.saveMemberData()
		if m.networkData.NodeInfo.IP != "" {
			m.joinedNetwork = true
			m.OpenTunAndSetIp()

			log.Printf("\n\nCongratulations!!! Your nConnect network member IP is: %v, mask is: %v\n\n",
				m.networkData.NodeInfo.IP, m.networkData.NodeInfo.Netmask)
		}
	} else {
		log.Println("You sent a join the network request to the manager, wait for the manager to authorize.")
	}

	return nil
}

func (m *Member) LeaveNetwork() error {
	msg := memberToManager{MsgType: LEAVE_NETWORK, Name: m.opts.NodeName}
	resp, err := SendMsg(m.c, m.opts.ManagerAddress, &msg, true)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}

	m.networkData = memberNetworkData{}
	return m.saveMemberData()
}

func (m *Member) SetServerTunnel(t *tunnel.Tunnel) error {
	m.serverTunnel = t
	serverAddress := t.FromAddr()
	m.serverAddress = serverAddress
	err := m.GetNodeIAccept()
	if err != nil {
		return err
	}
	if m.networkData.NodeInfo != nil && serverAddress == m.networkData.NodeInfo.ServerAddress {
		return nil
	}

	msg := memberToManager{MsgType: UPDATE_SERVER_ADDRESS, ServerAddress: serverAddress}
	_, err = SendMsg(m.c, m.opts.ManagerAddress, &msg, false)
	if err != nil {
		return err
	}

	if m.networkData.NodeInfo == nil {
		return nil
	}

	m.networkData.NodeInfo.ServerAddress = serverAddress
	return m.saveMemberData()
}

func (m *Member) GetNodeIAccept() error {
	msg := memberToManager{MsgType: GET_NODES_I_ACCEPT}
	resp, err := SendMsg(m.c, m.opts.ManagerAddress, &msg, true)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}

	if len(resp.NodeInfo) == 0 {
		return nil
	}

	m.networkData.NodesIAccept = resp.NodeInfo
	if err = m.saveMemberData(); err != nil {
		return err
	}

	m.UpdMyAccept(m.networkData.NodesIAccept)

	return nil
}

func (m *Member) UpdMyAccept(nodes []*NodeInfo) {
	if m.opts.Verbose {
		log.Printf("Network member, nodes I accept: %+v\n", nodes)
	}

	var addrs []string
	for _, node := range nodes {
		arr := strings.Split(node.Address, ".")
		addrs = append(addrs, arr[len(arr)-1]+"$")
	}

	if len(addrs) > 0 {
		err := m.opts.Config.AddAcceptAddrs(addrs)
		if err != nil {
			log.Println("Network member, opts.Config.AddAcceptAddrs error: ", err)
		}
		if m.serverTunnel != nil {
			err = m.serverTunnel.SetAcceptAddrs(nkn.NewStringArray(m.opts.Config.GetAcceptAddrs()...))
			if err != nil {
				log.Println("Network member, serverTunnel.SetAcceptAddrs error: ", err)
			}
		}
	}
}

func (m *Member) GetNodeICanAccess() error {
	msg := memberToManager{MsgType: GET_NODES_I_CAN_ACCESS, Name: m.opts.NodeName}
	resp, err := SendMsg(m.c, m.opts.ManagerAddress, &msg, true)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}

	if len(resp.NodeInfo) > 0 {
		m.networkData.NodesICanAccess = resp.NodeInfo
		if err = m.saveMemberData(); err != nil {
			return err
		}

		if m.CbNodeICanAccessUpdated != nil {
			m.CbNodeICanAccessUpdated(m.networkData.NodesICanAccess)
		}
	}

	return nil
}

func (m *Member) loadMemberData() error {
	jsonFile, err := os.OpenFile(memberFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	data := memberNetworkData{NetworkInfo: &networkInfo{}, NodeInfo: &NodeInfo{}}
	if len(b) == 0 {
		return errors.New(errNoDataInFile)
	}

	if err = json.Unmarshal(b, &data); err != nil {
		return err
	}
	if data.NetworkInfo != nil {
		m.networkData.NetworkInfo = data.NetworkInfo
	}
	if data.NodeInfo != nil {
		m.networkData.NodeInfo = data.NodeInfo
	}

	return nil
}

func (m *Member) saveMemberData() error {
	b, err := json.MarshalIndent(m.networkData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(memberFile, b, os.ModePerm)
}

func (m *Member) SetRoutes() error {
	routes := make([]string, 0, len(m.networkData.NodesIAccept))
	for _, n := range m.networkData.NodesIAccept {
		routes = append(routes, fmt.Sprintf("%s/32", n.IP))
	}

	ipNets := make([]*net.IPNet, len(routes))
	if len(routes) > 0 {
		for i, cidr := range routes {
			_, cidr, err := net.ParseCIDR(cidr)
			if err != nil {
				return fmt.Errorf("parse CIDR %s error: %v", cidr, err)
			}
			ipNets[i] = cidr
		}
	}
	arch.SetVPNRoutes(m.opts.TunName, m.networkData.NetworkInfo.Gateway, ipNets)

	return nil
}

func (m *Member) DeleteRoutes() error {
	routes := make([]string, 0, len(m.networkData.NodesIAccept))
	for _, n := range m.networkData.NodesIAccept {
		routes = append(routes, fmt.Sprintf("%s/32", n.IP))
	}

	ipNets := make([]*net.IPNet, len(routes))
	if len(routes) > 0 {
		for i, cidr := range routes {
			_, cidr, err := net.ParseCIDR(cidr)
			if err != nil {
				return fmt.Errorf("parse CIDR %s error: %v", cidr, err)
			}
			ipNets[i] = cidr
		}
	}

	arch.RemoveVPNRoutes(m.opts.TunName, m.networkData.NetworkInfo.Gateway, ipNets)

	return nil
}

func (m *Member) GetNodeInfo() *NodeInfo {
	return m.networkData.NodeInfo
}

func (m *Member) GetNetworkInfo() *networkInfo {
	return m.networkData.NetworkInfo
}

func (m *Member) OpenTunAndSetIp() {
	m.openTunOnce.Do(func() {
		err := arch.OpenTun(m.opts.TunName, m.networkData.NodeInfo.IP, m.networkData.NetworkInfo.Gateway, m.networkData.NodeInfo.Netmask, m.opts.TunDNS[0], m.opts.LocalSocksAddr)
		if err != nil {
			log.Printf("OpenTun error: %v", err)
		} else {
			log.Println("Started tun2socks, interface:", m.opts.TunName, "address:", m.networkData.NodeInfo.IP)
		}
	})
}
