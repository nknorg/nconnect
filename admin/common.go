package admin

import (
	"errors"
	"io/ioutil"
	"net"

	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/util"
	"github.com/nknorg/nkn-sdk-go"
	ts "github.com/nknorg/nkn-tuna-session"
	tunnel "github.com/nknorg/nkn-tunnel"
	"github.com/nknorg/tuna/filter"
	"github.com/nknorg/tuna/geo"
)

type permission uint8

const (
	rpcPermissionAcceptClient permission = 1 << iota
	rpcPermissionAdminClient
	rpcPermissionWeb
)

var (
	errUnknownMethod    = errors.New("unknown method")
	errPermissionDenied = errors.New("permission denied")
	resultSuccess       = "success"
)

var (
	rpcPermissions = map[string]permission{
		"getAdminToken":   rpcPermissionAdminClient | rpcPermissionWeb,
		"getAddrs":        rpcPermissionAdminClient | rpcPermissionWeb,
		"setAddrs":        rpcPermissionAdminClient | rpcPermissionWeb,
		"addAddrs":        rpcPermissionAdminClient | rpcPermissionWeb,
		"removeAddrs":     rpcPermissionAdminClient | rpcPermissionWeb,
		"getLocalIP":      rpcPermissionAcceptClient | rpcPermissionAdminClient | rpcPermissionWeb,
		"getInfo":         rpcPermissionAcceptClient | rpcPermissionAdminClient | rpcPermissionWeb,
		"getBalance":      rpcPermissionAcceptClient | rpcPermissionAdminClient | rpcPermissionWeb,
		"setAdminHttpApi": rpcPermissionAdminClient | rpcPermissionWeb,
		"getSeed":         rpcPermissionAdminClient | rpcPermissionWeb,
		"setSeed":         rpcPermissionAdminClient | rpcPermissionWeb,
		"setTunaConfig":   rpcPermissionAdminClient | rpcPermissionWeb,
		"getLog":          rpcPermissionAdminClient | rpcPermissionWeb,
	}
)

type RpcReq struct {
	ID      string                 `json:"id"`
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Token   string                 `json:"token"`
}

type RpcResp struct {
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

type GetInfoJSON struct {
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
	ServiceName     string   `json:"serviceName"`
	Country         []string `json:"country"`
	AllowNknAddr    []string `json:"allowNknAddr"`
	DisallowNknAddr []string `json:"disallowNknAddr"`
	AllowIp         []string `json:"allowIp"`
	DisallowIp      []string `json:"disallowIp"`
}

type getLogJSON struct {
	MaxSize int `json:"maxSize"`
}

func handleRequest(req *RpcReq, persistConf, mergedConf *config.Config, tun *tunnel.Tunnel, rpcPerm permission) *RpcResp {
	resp := &RpcResp{}

	if rpcPermissions[req.Method]&rpcPerm == 0 {
		resp.Error = errPermissionDenied.Error()
		return resp
	}

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
	case "getLog":
		params := &getLogJSON{}
		err := util.JSONConvert(req.Params, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		logContent, err := getLog(mergedConf, params)
		if err != nil {
			resp.Error = err.Error()
			break
		}
		resp.Result = logContent
	default:
		resp.Error = errUnknownMethod.Error()
	}
	return resp
}

func getAdminToken() *adminTokenJSON {
	if len(serverAdminAddr) == 0 {
		return nil
	}
	return &adminTokenJSON{
		Addr:  serverAdminAddr,
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

func getInfo(conf *config.Config, tun *tunnel.Tunnel) (*GetInfoJSON, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return nil, err
	}
	info := &GetInfoJSON{
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
		info.InPrice = make([]string, 0, len(tunaPubAddrs.Addrs))
		info.OutPrice = make([]string, 0, len(tunaPubAddrs.Addrs))
		for _, addr := range tunaPubAddrs.Addrs {
			if len(addr.IP) > 0 {
				info.InPrice = append(info.InPrice, addr.InPrice)
				info.OutPrice = append(info.OutPrice, addr.OutPrice)
			}
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
	err := persistConf.SetTunaConfig(params.ServiceName, params.Country, params.AllowNknAddr, params.DisallowNknAddr, params.AllowIp, params.DisallowIp)
	if err != nil {
		return err
	}
	err = mergedConf.SetTunaConfig(params.ServiceName, params.Country, params.AllowNknAddr, params.DisallowNknAddr, params.AllowIp, params.DisallowIp)
	if err != nil {
		return err
	}
	tsClient := tun.TunaSessionClient()
	if tsClient != nil {
		locations := make([]geo.Location, len(params.Country))
		for i := range params.Country {
			locations[i].CountryCode = params.Country[i]
		}
		allowIps := make([]geo.Location, len(params.AllowIp))
		for i := range params.AllowIp {
			allowIps[i].IP = params.AllowIp[i]
		}
		var allowed = append(locations, allowIps...)

		disallowed := make([]geo.Location, len(params.DisallowIp))
		for i := range params.DisallowIp {
			disallowed[i].IP = params.DisallowIp[i]
		}

		allowNknAddrs := make([]filter.NknClient, len(params.AllowNknAddr))
		for i := range params.AllowNknAddr {
			allowNknAddrs[i].Address = params.AllowNknAddr[i]
		}

		disallowNknAddrs := make([]filter.NknClient, len(params.DisallowNknAddr))
		for i := range params.DisallowNknAddr {
			disallowNknAddrs[i].Address = params.DisallowNknAddr[i]
		}

		err = tsClient.SetConfig(&ts.Config{
			TunaIPFilter:    &geo.IPFilter{Allow: allowed, Disallow: disallowed},
			TunaNknFilter:   &filter.NknFilter{Allow: allowNknAddrs, Disallow: disallowNknAddrs},
			TunaServiceName: params.ServiceName,
		})
		if err != nil {
			return err
		}
		go tsClient.RotateAll()
	}
	return nil
}

func getLog(conf *config.Config, params *getLogJSON) (string, error) {
	if len(conf.LogFileName) == 0 {
		return "", nil
	}
	b, err := ioutil.ReadFile(conf.LogFileName)
	if err != nil {
		return "", err
	}
	if conf.LogAPIResponseSize > 0 && len(b) > conf.LogAPIResponseSize {
		b = b[len(b)-conf.LogAPIResponseSize:]
	}
	if params.MaxSize > 0 && len(b) > params.MaxSize {
		b = b[len(b)-params.MaxSize:]
	}
	return string(b), nil
}
