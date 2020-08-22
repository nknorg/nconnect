package admin

import (
	"net"

	"github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/util"
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

type localIPJSON struct {
	Ipv4 []string `json:"ipv4"`
}

type getInfoJSON struct {
	Addr    string       `json:"addr"`
	LocalIP *localIPJSON `json:"localIP"`
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
	case "getLocalIP":
		localIP, err := getLocalIP()
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = localIP
	case "getInfo":
		info, err := getInfo(tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = info
	default:
		resp.Error = "unknown method"
	}
	return resp
}

func getAdminToken() *adminTokenJSON {
	if len(clientAddr) == 0 {
		return nil
	}
	return &adminTokenJSON{
		Addr:  clientAddr,
		Token: tokenStore.GetCurrentToken(),
	}
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

func getLocalIP() (*localIPJSON, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ipv4 := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ipv4 = append(ipv4, ip.String())
		}
	}
	return &localIPJSON{Ipv4: ipv4}, nil
}

func getInfo(tun *tunnel.Tunnel) (*getInfoJSON, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return nil, err
	}
	info := &getInfoJSON{
		Addr:    tun.FromAddr(),
		LocalIP: localIP,
	}
	return info, nil
}
