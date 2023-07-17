package network

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nkn-sdk-go"
)

const (
	defaultDomain  = "nconnect.nkn"
	defaultIpStart = "10.0.86.2"
	defaultIpEnd   = "10.0.86.254"
	defaultNetmask = "255.255.255.0"
	defaultGateway = "10.0.86.1"
	defaultDNS     = "1.1.1.1"

	networkDataFile = "network.json"

	AllMembers = "allMembers"
)

type networkData struct {
	NetworkInfo *networkInfo `json:"networkInfo"`
	IpStart     string       `json:"ipStart"` // start ip of the network
	IpEnd       string       `json:"ipEnd"`   // end ip of the network
	Netmask     string       `json:"netmask"` // mask of the network
	NextIp      string       `json:"nextIp"`  // next available ip

	Waiting       map[string]*NodeInfo `json:"waiting"`       // nodes waiting for authorization, map address to node info
	Member        map[string]*NodeInfo `json:"member"`        // authorized member list, map address to node info
	AcceptAddress map[string][]string  `json:"acceptAddress"` // accept address data for each member. map address to list of accepted address
	NameToAddress map[string]string    `json:"nameToAddress"` // map name to address

	ManagerBalance string `json:"managerBalance"` // manager's NKN balance
}

type Manager struct {
	opts *config.Opts
	c    *admin.Client

	sync.RWMutex
	networkData *networkData // persisted data that will be saved to disk
}

var manager *Manager

func NewManager(account *nkn.Account, clientConfig *nkn.ClientConfig, opts *config.Opts) (*Manager, error) {
	if manager != nil {
		return manager, nil
	}

	manager = &Manager{opts: opts}
	err := manager.loadNetworkData()
	if err != nil {
		return nil, err
	}

	if len(opts.Identifier) == 0 {
		return nil, errors.New("network manager's identifier should not be empty")
	}

	manager.c, err = admin.NewClient(account, clientConfig, opts.Identifier)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *Manager) StartManager() error {
	log.Println("nConnect manager is listening at:", m.c.MultiClient.Address())

	for {
		msg := <-m.c.MultiClient.OnMessage.C
		resp, err := m.handleRequest(msg)
		if err != nil {
			log.Println("nConnect manager handle request error", err)
			continue
		}

		b, err := json.Marshal(resp)
		if err != nil {
			log.Println("nConnect manager json.Marshal resp error", err)
			continue
		}

		err = msg.Reply(b)
		if err != nil {
			log.Println("nConnect manager msg.Reply error", err)
			continue
		}
	}
}

func (m *Manager) handleRequest(msg *nkn.Message) (*managerToMember, error) {
	req := &memberToManager{}
	err := json.Unmarshal(msg.Data, req)
	if err != nil {
		return nil, err
	}

	var node *NodeInfo
	resp := &managerToMember{}
	resp.MsgType = req.MsgType

	switch req.MsgType {
	case JOIN_NETWORK:
		resp.NetworkInfo = m.networkData.NetworkInfo
		node, err = m.JoinNetwork(msg.Src, req.Name, req.ServerAddress)
		if node != nil {
			resp.NodeInfo = append(resp.NodeInfo, node)
		}

	case LEAVE_NETWORK:
		err = m.LeaveNetwork(msg.Src, req.Name)

	case GET_MY_INFO:
		resp.NetworkInfo = m.networkData.NetworkInfo
		if n := m.GetNodeInfo(msg.Src); n != nil {
			resp.NodeInfo = append(resp.NodeInfo, n)
		} else {
			err = errors.New(errNodeNotFound)
		}

	case GET_NODES_I_ACCEPT:
		list := m.GetAcceptNodes(msg.Src)
		resp.NodeInfo = list

	case GET_NODES_I_CAN_ACCESS:
		list := m.GetNodesICanAccess(msg.Src)
		resp.NodeInfo = list

	case UPDATE_SERVER_ADDRESS:
		err = m.SetNodeServerAddress(msg.Src, req.ServerAddress)

	case NKN_PING:
		fmt.Println("Got ping from", msg.Src)
		resp.MsgType = NKN_PONG

	case NKN_PONG:
		fmt.Println("Got pong from", msg.Src)

	default:
		return nil, fmt.Errorf("nConnect manager got unknown message type: %v", req.MsgType)
	}

	if err != nil {
		resp.Err = err.Error()
	}

	return resp, nil
}

func (m *Manager) JoinNetwork(address, name, serverAddr string) (*NodeInfo, error) {
	if m.opts.Verbose {
		log.Println("A new member is joining network:", name, address)
	}

	m.RLock()
	node, ok := m.networkData.Member[address]
	m.RUnlock()
	if ok {
		node.LastSeen = time.Now()
		if name != "" {
			node.Name = name
		}
		if serverAddr != "" {
			node.ServerAddress = serverAddr
			node.Server = true
		}

		m.Lock()
		m.networkData.Member[address] = node
		m.Unlock()

		if err := m.saveNetworkData(); err != nil {
			return nil, err
		}

		notification := &managerToMember{
			MsgType:  NOTI_MEMBER_ONLINE,
			NodeInfo: []*NodeInfo{node},
		}
		// broadcast member online event to related members
		m.NotifyIAccept(address, notification)

		log.Printf("The member '%v' is online, its IP is %v\n", node.Name, node.IP)

		return node, nil
	}

	m.Lock()
	defer m.Unlock()
	if node, ok := m.networkData.Waiting[address]; ok {
		changed := false
		if name != "" && name != node.Name {
			node.Name = name
			changed = true
		}
		if serverAddr != "" && serverAddr != node.ServerAddress {
			node.ServerAddress = serverAddr
			changed = true
		}
		if changed {
			m.networkData.Waiting[address] = node
			if err := m.saveNetworkData(); err != nil {
				return nil, err
			}
		}

		return nil, errors.New(errWaitForAuth)
	}

	if _, nameExists := m.networkData.NameToAddress[name]; nameExists {
		return nil, errors.New(errNameExist)
	}

	if len(name) > 0 {
		m.networkData.NameToAddress[name] = address
	}
	m.networkData.Waiting[address] = &NodeInfo{Name: name, Address: address, ServerAddress: address, LastSeen: time.Now()}
	if err := m.saveNetworkData(); err != nil {
		return nil, err
	}

	log.Println("A new member is waiting for authorization:", name, address)

	return nil, errors.New(errWaitForAuth)
}

func (m *Manager) LeaveNetwork(address, name string) error {
	m.Lock()
	delete(m.networkData.NameToAddress, name)
	_, ok := m.networkData.Member[address]
	m.Unlock()

	if ok {
		m.Lock()
		delete(m.networkData.Member, address)
		delete(m.networkData.AcceptAddress, address)

		for _, n := range m.networkData.Member {
			acceptAddrs := m.networkData.AcceptAddress[n.Address]
			for i, addr := range acceptAddrs {
				if addr == address {
					acceptAddrs = append(acceptAddrs[:i], acceptAddrs[i+1:]...)
					m.networkData.AcceptAddress[n.Address] = acceptAddrs
					break
				}
			}
		}
		m.Unlock()

		notification := &managerToMember{
			MsgType:  NOTI_LEAVE_NETWORK,
			NodeInfo: []*NodeInfo{{Name: name, Address: address}},
		}
		m.NotifyIAccept(address, notification)
		m.NotifyICanAccess(address, notification)

	} else {
		m.Lock()
		delete(m.networkData.Waiting, address)
		m.Unlock()
	}

	log.Printf("The node %v left network, its address is %v\n", name, address)

	return m.saveNetworkData()
}

func (m *Manager) AuthorizeMemeber(address string) error {
	m.RLock()
	nw, ok := m.networkData.Waiting[address]
	m.RUnlock()

	if !ok {
		return errors.New(errNodeNotFound)
	}

	ip, err := m.GetAvailableIp()
	if err != nil {
		return err
	}
	nw.IP = ip
	nw.Netmask = m.networkData.Netmask

	m.Lock()
	m.networkData.Member[address] = nw
	delete(m.networkData.Waiting, address)
	m.Unlock()

	if err = m.saveNetworkData(); err != nil {
		return err
	}

	notification := &managerToMember{
		MsgType:     NOTI_AUTHORIZED,
		NetworkInfo: m.networkData.NetworkInfo,
		NodeInfo:    []*NodeInfo{nw},
	}

	if _, err = SendMsg(m.c, address, notification, false); err != nil {
		return err
	}

	m.NotifyICanAccess(address, &managerToMember{MsgType: NOTI_NEW_MEMBER})

	log.Println("You just authorized a new member:", nw.Name, nw.IP)

	return nil
}

func (m *Manager) DeleteWaiting(address string) error {
	m.Lock()
	delete(m.networkData.Waiting, address)
	m.Unlock()

	return m.saveNetworkData()
}

func (m *Manager) RemoveMember(address string) error {
	m.RLock()
	nw, ok := m.networkData.Member[address]
	m.RUnlock()

	if ok {
		m.Lock()
		m.networkData.Waiting[address] = nw
		delete(m.networkData.Member, address)
		delete(m.networkData.AcceptAddress, address)
		m.Unlock()

		err := m.saveNetworkData()
		if err != nil {
			return err
		}

		notification := &managerToMember{
			MsgType:  NOTI_LEAVE_NETWORK,
			NodeInfo: []*NodeInfo{nw},
		}
		m.NotifyIAccept(address, notification)
		m.NotifyICanAccess(address, notification)
	}

	log.Println("You just removed a member:", nw.Name, nw.IP)

	return nil
}

func (m *Manager) SetNodeServerAddress(address, serverAddress string) error {
	m.Lock()
	defer m.Unlock()

	if serverAddress == "" {
		return nil
	}

	if node, ok := m.networkData.Member[address]; ok && node.ServerAddress != serverAddress {
		node.ServerAddress = serverAddress
		m.networkData.Member[address] = node
		return m.saveNetworkData()
	}

	if node, ok := m.networkData.Waiting[address]; ok && node.ServerAddress != serverAddress {
		node.ServerAddress = serverAddress
		m.networkData.Waiting[address] = node
		return m.saveNetworkData()
	}

	return nil
}

func (m *Manager) GetAcceptAddress(address string) []string {
	m.RLock()
	defer m.RUnlock()
	list := m.networkData.AcceptAddress[address]
	return list
}

func (m *Manager) GetAcceptNodes(address string) []*NodeInfo {
	m.RLock()
	defer m.RUnlock()

	addressList := m.networkData.AcceptAddress[address]
	var list []*NodeInfo
	if len(addressList) > 0 && addressList[0] == AllMembers {
		for _, n := range m.networkData.Member {
			if n.Address != address {
				list = append(list, n)
			}
		}
	} else {
		for _, a := range addressList {
			if n, ok := m.networkData.Member[a]; ok {
				list = append(list, n)
			}
		}
	}

	return list
}

func (m *Manager) GetNodesICanAccess(address string) []*NodeInfo {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.networkData.Member[address]; !ok { // not a member
		return nil
	}

	var list []*NodeInfo
	for addr, acceptAddress := range m.networkData.AcceptAddress {
		if addr == address {
			continue
		}

		if len(acceptAddress) > 0 && acceptAddress[0] == AllMembers {
			if n, ok := m.networkData.Member[addr]; ok {
				list = append(list, n)
			}
			continue
		}

		for _, a := range acceptAddress {
			if a == address {
				if n, ok := m.networkData.Member[addr]; ok {
					list = append(list, n)
				}
			}
		}
	}

	return list
}

func (m *Manager) SetAcceptAddress(address string, acceptAddress []string) error {
	m.Lock()
	m.networkData.AcceptAddress[address] = acceptAddress
	m.Unlock()

	if err := m.saveNetworkData(); err != nil {
		return err
	}
	notification := &managerToMember{MsgType: NOTI_UPD_I_ACCEPT}
	if _, err := SendMsg(m.c, address, notification, false); err != nil {
		return err
	}

	m.RLock()
	n := m.networkData.Member[address]
	m.RUnlock()

	notification = &managerToMember{MsgType: NOTI_UPD_I_CAN_ACCESS, NodeInfo: []*NodeInfo{n}}
	m.NotifyIAccept(address, notification)

	return nil
}

func (m *Manager) SendToken(address, amount string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	walletAddr, err := nkn.ClientAddrToWalletAddr(address)
	if err != nil {
		log.Printf("ClientAddrToWalletAddr client address %v error: %v", address, err)
		return err
	}

	_, err = nkn.TransferContext(ctx, m.c.MultiClient, walletAddr, amount, nkn.GetDefaultTransactionConfig())
	if err != nil {
		log.Printf("Manager sent %v NKN to %v err %v\n", amount, address, err)
		return err
	}
	log.Printf("Manager sent %v NKN to %v successfully\n", amount, address)

	return nil
}

func (m *Manager) NknPing(address string) (int, error) {
	msg := &managerToMember{MsgType: NKN_PING}
	start := time.Now()
	_, err := SendMsg(m.c, address, msg, true)
	if err != nil {
		return 0, err
	}
	rtt := time.Since(start)

	return int(rtt.Milliseconds()), nil
}

// Send notification to all the nodes which I(initiatorAddr) accept
func (m *Manager) NotifyIAccept(initiatorAddr string, notification *managerToMember) error {
	m.RLock()
	acceptAddress := m.networkData.AcceptAddress[initiatorAddr]
	m.RUnlock()

	if len(acceptAddress) > 0 && acceptAddress[0] == AllMembers {
		for _, n := range m.networkData.Member {
			if n.Address != initiatorAddr {
				if _, err := SendMsg(m.c, n.Address, notification, false); err != nil {
					log.Printf("Send msg type %v to %v error %v\n", notification.MsgType, n.Address, err)
				}
			}
		}
	} else {
		for _, addr := range acceptAddress {
			if _, err := SendMsg(m.c, addr, notification, false); err != nil {
				log.Printf("Send msg type %v to %v error %v\n", notification.MsgType, addr, err)
			}
		}
	}

	return nil
}

// Send notification to all the nodes which accept me(initiatorAddr)
func (m *Manager) NotifyICanAccess(initiatorAddr string, notification *managerToMember) error {
	for _, n := range m.networkData.Member {
		if n.Address == initiatorAddr {
			continue
		}

		acceptAddr := m.networkData.AcceptAddress[n.Address]
		if len(acceptAddr) > 0 && acceptAddr[0] == AllMembers {
			if _, err := SendMsg(m.c, n.Address, notification, false); err != nil {
				log.Printf("Send msg type %v to %v error %v\n", notification.MsgType, n.Address, err)
			}
		} else {
			for _, addr := range acceptAddr { // broadcast accept info to nodes
				if addr == initiatorAddr {
					if _, err := SendMsg(m.c, n.Address, notification, false); err != nil {
						log.Printf("Send msg type %v to %v error %v\n", notification.MsgType, addr, err)
					}
				}
			}
		}
	}

	return nil
}

func (m *Manager) GetNodeInfo(address string) *NodeInfo {
	m.RLock()
	defer m.RUnlock()
	if n, ok := m.networkData.Member[address]; ok {
		return n
	}
	return nil
}

type network struct {
	NetworkData    *networkData `json:"networkData"`    // network data
	ManagerAddress string       `json:"managerAddress"` // manager's NKN address
	ManagerBalance string       `json:"managerBalance"` // manager's NKN balance
}

func (m *Manager) GetNetworkConfig() *network {
	m.RLock()
	defer m.RUnlock()
	for _, n := range m.networkData.Member {
		if n.Server {
			n.Balance = getBalance(n.ServerAddress)
		}
	}

	managerBalance := getBalance(m.c.MultiClient.Address())

	return &network{NetworkData: m.networkData, ManagerAddress: m.c.MultiClient.Address(), ManagerBalance: managerBalance}
}

func (m *Manager) SetNetworkConfig(conf *networkData) error {
	m.Lock()
	defer m.Unlock()
	m.networkData.NetworkInfo = conf.NetworkInfo
	m.networkData.IpStart = conf.IpStart
	m.networkData.IpEnd = conf.IpEnd
	m.networkData.Netmask = conf.Netmask

	return m.saveNetworkData()
}

func (m *Manager) GetAvailableIp() (string, error) {
	m.Lock()
	defer m.Unlock()
	if m.networkData.NextIp == "" {
		return "", errors.New("nConnect manager has no available ip")
	}

	ip := m.networkData.NextIp
	m.networkData.NextIp = int2ip(ip2int(m.networkData.NextIp) + 1)
	if m.networkData.NextIp == m.networkData.IpEnd {
		m.networkData.NextIp = ""
	}

	return ip, nil
}

func (m *Manager) loadNetworkData() error {
	m.Lock()
	defer m.Unlock()

	nwData := &networkData{
		Waiting:       make(map[string]*NodeInfo),
		Member:        make(map[string]*NodeInfo),
		AcceptAddress: make(map[string][]string),
		NameToAddress: make(map[string]string),
	}
	m.networkData = nwData

	jsonFile, err := os.OpenFile(networkDataFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		return json.Unmarshal(b, nwData)
	} else { // set default value
		nwData.IpStart = defaultIpStart // IpStart is reserved for manager
		nwData.IpEnd = defaultIpEnd
		nwData.Netmask = defaultNetmask
		ipNext := int2ip(ip2int(defaultIpStart) + 1)
		nwData.NextIp = ipNext
		nwData.NetworkInfo = &networkInfo{Domain: defaultDomain, Gateway: defaultGateway, DNS: defaultDNS}
		return m.saveNetworkData()
	}
}

func (m *Manager) saveNetworkData() error {
	if m.networkData == nil {
		return errors.New("networkData is nil")
	}

	b, err := json.MarshalIndent(m.networkData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(networkDataFile, b, os.ModePerm)
}

func ip2int(ip string) uint32 {
	s := net.ParseIP(ip).To4()
	return binary.BigEndian.Uint32(s)
}

func int2ip(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

func getBalance(serverAddr string) string {
	walletAddr, err := nkn.ClientAddrToWalletAddr(serverAddr)
	if err != nil {
		log.Printf("ClientAddrToWalletAddr client address %v error: %v", serverAddr, err)
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	amount, err := nkn.GetBalanceContext(ctx, walletAddr, nkn.GetDefaultRPCConfig())
	if err != nil {
		log.Printf("GetBalanceContext of %v error: %v", walletAddr, err)
		return ""
	}

	return amount.String()
}
