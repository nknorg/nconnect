package admin

import (
	"errors"
	"net"

	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/util"
	"github.com/nknorg/nkn-sdk-go"
	ts "github.com/nknorg/nkn-tuna-session"
	tunnel "github.com/nknorg/nkn-tunnel"
	"github.com/nknorg/tuna/geo"
)

var (
	errUnknownMethod = errors.New("unknown method")
	resultSuccess    = "success"
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
	Addr                 string       `json:"addr"`
	LocalIP              *localIPJSON `json:"localIP"`
	AdminHTTPAPIDisabled bool         `json:"adminHttpApiDisabled"`
	Version              string       `json:"version"`
	Tuna                 bool         `json:"tuna"`
	TunaServiceName      string       `json:"tunaServiceName,omitempty"`
	TunaCountry          []string     `json:"tunaCountry,omitempty"`
	InPrice              []string     `json:"inPrice,omitempty"`
	OutPrice             []string     `json:"outPrice,omitempty"`
	Tags                 []string     `json:"tags,omitempty"`
}

type setSeedJSON struct {
	Seed string `json:"seed"`
}

type adminHTTPAPIJSON struct {
	Disable bool `json:"disable"`
}

type tunaConfigJSON struct {
	ServiceName string   `json:"serviceName"`
	Country     []string `json:"country"`
}

func handleRequest(req *rpcReq, persistConf, mergedConf *config.Config, tun *tunnel.Tunnel) *rpcResp {
	resp := &rpcResp{}
	switch req.Method {
	case "getAdminToken":
		resp.Result = getAdminToken()
	case "getAddrs":
		resp.Result = getAddrs(persistConf)
	case "setAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = setAddrs(persistConf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(persistConf)
	case "addAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = addAddrs(persistConf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(persistConf)
	case "removeAddrs":
		addrs := &addrsJSON{}
		err := util.JSONConvert(req.Params, addrs)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = removeAddrs(persistConf, addrs, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = getAddrs(persistConf)
	case "getLocalIP":
		localIP, err := getLocalIP()
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = localIP
	case "getInfo":
		info, err := getInfo(mergedConf, tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = info
	case "getBalance":
		balance, err := getBalance(tun)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = balance
	case "setAdminHttpApi":
		params := &adminHTTPAPIJSON{}
		err := util.JSONConvert(req.Params, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = setAdminHTTPAPI(persistConf, mergedConf, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = resultSuccess
	case "getSeed":
		resp.Result = mergedConf.Seed
	case "setSeed":
		params := &setSeedJSON{}
		err := util.JSONConvert(req.Params, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = persistConf.SetSeed(params.Seed)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = resultSuccess
	case "setTunaConfig":
		params := &tunaConfigJSON{}
		err := util.JSONConvert(req.Params, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		err = setTunaConfig(tun, persistConf, mergedConf, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = resultSuccess
	default:
		resp.Error = errUnknownMethod.Error()
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

func getInfo(conf *config.Config, tun *tunnel.Tunnel) (*getInfoJSON, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return nil, err
	}
	info := &getInfoJSON{
		Addr:                 tun.FromAddr(),
		LocalIP:              localIP,
		AdminHTTPAPIDisabled: conf.DisableAdminHTTPAPI,
		Tuna:                 conf.Tuna,
		TunaServiceName:      conf.TunaServiceName,
		TunaCountry:          conf.TunaCountry,
		Version:              config.Version,
	}
	tunaPubAddrs := tun.TunaPubAddrs()
	if tunaPubAddrs != nil {
		info.InPrice = make([]string, len(tunaPubAddrs.Addrs))
		info.OutPrice = make([]string, len(tunaPubAddrs.Addrs))
		for i := range tunaPubAddrs.Addrs {
			info.InPrice[i] = tunaPubAddrs.Addrs[i].InPrice
			info.OutPrice[i] = tunaPubAddrs.Addrs[i].OutPrice
		}
	}
	if len(conf.Tags) > 0 {
		info.Tags = conf.Tags
	}
	return info, nil
}

func getBalance(tun *tunnel.Tunnel) (string, error) {
	balance, err := tun.MultiClient().Balance()
	if err != nil {
		return "", err
	}
	return balance.String(), nil
}

func setAdminHTTPAPI(persistConf, mergedConf *config.Config, params *adminHTTPAPIJSON) error {
	err := persistConf.SetAdminHTTPAPI(params.Disable)
	if err != nil {
		return err
	}
	return mergedConf.SetAdminHTTPAPI(params.Disable)
}

func setTunaConfig(tun *tunnel.Tunnel, persistConf, mergedConf *config.Config, params *tunaConfigJSON) error {
	err := persistConf.SetTunaConfig(params.ServiceName, params.Country)
	if err != nil {
		return err
	}
	err = mergedConf.SetTunaConfig(params.ServiceName, params.Country)
	if err != nil {
		return err
	}
	tsClient := tun.TunaSessionClient()
	if tsClient != nil {
		locations := make([]geo.Location, len(params.Country))
		for i := range params.Country {
			locations[i].CountryCode = params.Country[i]
		}
		err = tsClient.SetConfig(&ts.Config{
			TunaIPFilter:    &geo.IPFilter{Allow: locations},
			TunaServiceName: params.ServiceName,
		})
		if err != nil {
			return err
		}
		go tsClient.RotateAll()
	}
	return nil
}
