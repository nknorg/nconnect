package admin

import (
	"github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn-socks/config"
	"github.com/nknorg/nkn-socks/util"
	tunnel "github.com/nknorg/nkn-tunnel"
)

type rpcReq struct {
	ID      string                 `json:"id"`
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Token   string                 `json:"token"`
}

type rpcResp struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type addrsJSON struct {
	AcceptAddrs []string `json:"acceptAddrs"`
	AdminAddrs  []string `json:"adminAddrs"`
}

type adminTokenJSON struct {
	Addr  string `json:"addr"`
	Token *Token `json:"token"`
}

func handleRequest(req *rpcReq, conf *config.Config, tun *tunnel.Tunnel) *rpcResp {
	resp := &rpcResp{}
	switch req.Method {
	case "getAdminToken":
		resp.Result = getAdminToken()
	case "getAddrs":
		resp.Result = getAddrs(conf)
	case "setAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = setAddrs(conf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(conf)
	case "addAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = addAddrs(conf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(conf)
	case "removeAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = removeAddrs(conf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(conf)
	default:
		resp.Error = "unknown method"
	}
	return resp
}

func getAddrs(conf *config.Config) *addrsJSON {
	return &addrsJSON{
		AcceptAddrs: conf.GetAcceptAddrs(),
		AdminAddrs:  conf.GetAdminAddrs(),
	}
}

func setAddrs(conf *config.Config, addrs *addrsJSON, tun *tunnel.Tunnel) error {
	if addrs.AcceptAddrs != nil {
		conf.SetAcceptAddrs(addrs.AcceptAddrs)
	}
	if addrs.AdminAddrs != nil {
		conf.SetAdminAddrs(addrs.AdminAddrs)
	}
	return tun.SetAcceptAddrs(nkn.NewStringArray(conf.GetAcceptAddrs()...))
}

func addAddrs(conf *config.Config, addrs *addrsJSON, tun *tunnel.Tunnel) error {
	if addrs.AcceptAddrs != nil {
		conf.AddAcceptAddrs(addrs.AcceptAddrs)
	}
	if addrs.AdminAddrs != nil {
		conf.AddAdminAddrs(addrs.AdminAddrs)
	}
	return tun.SetAcceptAddrs(nkn.NewStringArray(conf.GetAcceptAddrs()...))
}

func removeAddrs(conf *config.Config, addrs *addrsJSON, tun *tunnel.Tunnel) error {
	if addrs.AcceptAddrs != nil {
		conf.RemoveAcceptAddrs(addrs.AcceptAddrs)
	}
	if addrs.AdminAddrs != nil {
		conf.RemoveAdminAddrs(addrs.AdminAddrs)
	}
	return tun.SetAcceptAddrs(nkn.NewStringArray(conf.GetAcceptAddrs()...))
}
