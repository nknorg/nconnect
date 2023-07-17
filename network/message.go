package network

import (
	"encoding/json"
	"time"

	"github.com/nknorg/nconnect/admin"
)

// msgType constants
const (
	MT_NONE = iota
	JOIN_NETWORK
	UPDATE_MY_INFO
	GET_MY_INFO
	UPDATE_SERVER_ADDRESS
	GET_NODES_I_ACCEPT
	GET_NODES_I_CAN_ACCESS
	LEAVE_NETWORK
	NKN_PING
	NKN_PONG

	NOTI_AUTHORIZED
	NOTI_NEW_MEMBER
	NOTI_UPD_I_CAN_ACCESS
	NOTI_UPD_I_ACCEPT
	NOTI_MEMBER_ONLINE
	NOTI_LEAVE_NETWORK
)

type NodeInfo struct {
	IP            string    `json:"ip"`
	Netmask       string    `json:"netmask"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`       // client address
	ServerAddress string    `json:"serverAddress"` // nconnect server listen address
	LastSeen      time.Time `json:"lastSeen"`
	Server        bool      `json:"server"`
	Balance       string    `json:"balance"`
}

type networkInfo struct {
	Domain  string `json:"domain"`
	Gateway string `json:"gateway"`
	DNS     string `json:"dns"`
}

type memberToManager struct {
	MsgType       int    `json:"msgType"`
	Name          string `json:"name"`
	ServerAddress string `json:"serverAddress"`
}

type managerToMember struct {
	MsgType     int          `json:"msgType"`
	Err         string       `json:"err"`
	NetworkInfo *networkInfo `json:"networkInfo"`
	NodeInfo    []*NodeInfo  `json:"nodeInfo"`
}

func SendMsg(mc *admin.Client, address string, msg interface{}, waitResponse bool) (*managerToMember, error) {
	reply, err := mc.SendMsg(address, msg, waitResponse)
	if err != nil || !waitResponse {
		return nil, err
	}

	var respMsg managerToMember
	if err = json.Unmarshal(reply.Data, &respMsg); err != nil {
		return nil, err
	}
	return &respMsg, nil
}
