package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	Cli_RPC = "127.0.0.1:10032"
)

const (
	Cli_Status = iota // Get my network status
	Cli_List          // List all nodes I can access, and all nodes I accept
	Cli_Join          // Join a network
	Cli_Leave         // Leave a network
)

type CliMsgReq struct {
	MsgType int `json:"msgType"`
}

type CliMsgResp struct {
	MsgType       int          `json:"msgType"`
	Err           string       `json:"err"`
	NetworkInfo   *networkInfo `json:"networkInfo"`
	NodeInfo      *NodeInfo    `json:"nodeInfo"`
	NodeICanAcces []*NodeInfo  `json:"nodeICanAccess"`
	NodeIAccept   []*NodeInfo  `json:"nodeIAccept"`
}

func (m *Member) StartCliService() error {
	a, err := net.ResolveUDPAddr("udp", Cli_RPC)
	if err != nil {
		return err
	}
	udpServer, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}

	defer udpServer.Close()

	b := make([]byte, 1024)
	var req CliMsgReq
	var resp CliMsgResp
	for {
		n, addr, err := udpServer.ReadFromUDP(b)
		if err != nil {
			log.Printf("StartCliService.ReadFromUDP err: %v", err)
			time.Sleep(time.Second)
			continue
		}

		err = json.Unmarshal(b[:n], &req)
		if err != nil {
			log.Printf("StartCliService.Unmarshal err: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		resp.MsgType = req.MsgType
		switch req.MsgType {
		case Cli_Status:
			resp.NetworkInfo = m.networkData.NetworkInfo
			resp.NodeInfo = m.networkData.NodeInfo
		case Cli_List:
			m.GetNodeIAccept()
			m.GetNodeICanAccess()
			resp.NetworkInfo = m.networkData.NetworkInfo
			resp.NodeIAccept = m.networkData.NodesIAccept
			resp.NodeICanAcces = m.networkData.NodesICanAccess
		case Cli_Join:
			m.JoinNetwork(m.serverAddress)
			resp.NetworkInfo = m.networkData.NetworkInfo
			resp.NodeInfo = m.networkData.NodeInfo
		case Cli_Leave:
			m.LeaveNetwork()
			resp.NetworkInfo = m.networkData.NetworkInfo

		default:
			resp.Err = "Unknown msgType"
		}

		buf, err := json.Marshal(resp)
		if err != nil {
			log.Printf("StartCliService.Marshal err: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		_, _, err = udpServer.WriteMsgUDP(buf, nil, addr)
		if err != nil {
			log.Printf("StartCliService.WriteMsgUDP err: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func CliStatus() {
	resp, err := CliRequest(Cli_Status)
	if err != nil {
		fmt.Println("CliStatus err: ", err)
		return
	}
	fmt.Println("\nNetwork Domain: ", resp.NetworkInfo.Domain)
	if resp.NodeInfo != nil && resp.NodeInfo.IP != "" {
		fmt.Println("Ip:", resp.NodeInfo.IP, "\tMask:", resp.NodeInfo.Netmask, "\tNode Name:", resp.NodeInfo.Name)
	} else {
		fmt.Println("You don't join the network yet")
	}
}

func CliList() {
	resp, err := CliRequest(Cli_List)
	if err != nil {
		fmt.Println("CliList err: ", err)
		return
	}
	fmt.Println("\nNodes I accept:")
	for _, node := range resp.NodeIAccept {
		fmt.Println("IP:", node.IP, "\tMask:", node.Netmask, "\tNode Name:", node.Name)
	}
	fmt.Println("\nNodes I can access:")
	for _, node := range resp.NodeICanAcces {
		fmt.Println("IP:", node.IP, "\tMask:", node.Netmask, "\tNode Name:", node.Name)
	}
}

func CliJoin() {
	resp, err := CliRequest(Cli_Join)
	if err != nil {
		fmt.Println("CliJoin err: ", err)
		return
	}
	fmt.Println("\nNetwork Domain: ", resp.NetworkInfo.Domain)
	if resp.NodeInfo != nil && resp.NodeInfo.IP != "" {
		fmt.Println("You have joined the network")
		fmt.Println("Ip:", resp.NodeInfo.IP, "\tMask:", resp.NodeInfo.Netmask, "\tNode Name:", resp.NodeInfo.Name)
	} else {
		fmt.Println("Join network request is sent, please wait for the manager to authorize it")
	}
}

func CliLeave() {
	resp, err := CliRequest(Cli_Leave)
	if err != nil {
		fmt.Println("CliLeave err: ", err)
		return
	}
	if resp.NetworkInfo == nil {
		fmt.Println("You have left the network")
	} else {
		fmt.Println("Leave network failed, please try again")
	}
}

func CliRequest(msgType int) (resp CliMsgResp, err error) {
	req := CliMsgReq{
		MsgType: msgType,
	}

	uc, err := net.Dial("udp", Cli_RPC)
	if err != nil {
		log.Println("CliRequest net.Dial err: ", err)
		return
	}
	defer uc.Close()

	send, _ := json.Marshal(req)
	if _, err = uc.Write(send); err != nil {
		log.Println("CliRequest.Write err ", err)
		return
	}

	b := make([]byte, 65535)
	n, err := uc.Read(b)
	if err != nil {
		log.Println("CliRequest.Read err ", err)
		return
	}

	err = json.Unmarshal(b[:n], &resp)
	if err != nil {
		log.Printf("CliRequest.Unmarshal err: %v\n", err)
		return
	}

	return
}
