package nconnect

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/imdario/mergo"
	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/arch"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/network"
	"github.com/nknorg/nconnect/ss"
	"github.com/nknorg/nconnect/util"
	"github.com/nknorg/ncp-go"
	"github.com/nknorg/nkn-sdk-go"
	ts "github.com/nknorg/nkn-tuna-session"
	tunnel "github.com/nknorg/nkn-tunnel"
	"github.com/nknorg/nkn/v2/common"
	"github.com/nknorg/nkn/v2/util/address"
	"github.com/nknorg/nkngomobile"
	"github.com/nknorg/tuna/filter"
	"github.com/nknorg/tuna/geo"
	"github.com/nknorg/tuna/types"
	"gopkg.in/natefinch/lumberjack.v2"
)

type nconnect struct {
	opts    *config.Opts
	account *nkn.Account

	walletConfig   *nkn.WalletConfig
	clientConfig   *nkn.ClientConfig
	tunnelConfig   *tunnel.Config
	ssClientConfig *ss.Config
	ssServerConfig *ss.Config
	persistConf    *config.Config

	adminClientCache   *admin.Client
	sync.RWMutex                                     // lock for maps
	remoteInfoCache    map[string]*admin.GetInfoJSON // map remote admin address to remote info
	remoteInfoByTunnel map[string]*admin.GetInfoJSON // map tunnel address to remote info

	clientTunnels []*tunnel.Tunnel // tunnels for client mode
	serverTunnel  *tunnel.Tunnel   // tunnel for server mode
	serverReady   chan struct{}    // channel to notify server is ready

	tunaNode *types.Node // It is used to connect specified tuna node, mainly is for testing.

	networkMember  *network.Member
	networkTunnels map[string]*tunnel.Tunnel // tunnels for network nodes
	routeCIDRs     []*net.IPNet              // CIDRs for routing traffic through network nodes
}

func NewNconnect(opts *config.Opts) (*nconnect, error) {
	err := (&opts.Config).SetPlatformSpecificDefaultValues()
	if err != nil {
		return nil, err
	}

	if !(opts.NetworkMember || opts.NetworkManager) {
		if opts.Client == opts.Server {
			log.Fatal("Exactly one mode (client or server) should be selected if not join a network.")
		}
	}

	persistConf, err := config.LoadOrNewConfig(opts.ConfigFile)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(&opts.Config, persistConf)
	if err != nil {
		return nil, err
	}

	if len(opts.LogFileName) > 0 {
		log.SetOutput(&lumberjack.Logger{
			Filename:   opts.LogFileName,
			MaxSize:    opts.LogMaxSize,
			MaxBackups: opts.LogMaxBackups,
		})
	}

	seed, err := hex.DecodeString(opts.Seed)
	if err != nil {
		return nil, err
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		return nil, err
	}

	shouldSave := false
	if len(opts.Seed) == 0 {
		persistConf.Seed = hex.EncodeToString(account.Seed())
		opts.Seed = persistConf.Seed
		shouldSave = true
	}

	if len(opts.Identifier) == 0 {
		persistConf.Identifier = config.RandomIdentifier()
		opts.Identifier = persistConf.Identifier
		shouldSave = true
	}

	if shouldSave {
		err = persistConf.Save()
		if err != nil {
			return nil, err
		}
	}

	if opts.Address {
		addr := address.MakeAddressString(account.PubKey(), opts.Identifier)
		if opts.Server && len(opts.AdminIdentifier) > 0 {
			addr = opts.AdminIdentifier + "." + addr
		}
		fmt.Println(addr)
		os.Exit(0)
	}

	if opts.WalletAddress {
		if opts.Server {
			fmt.Println(account.WalletAddress())
		} else {
			fmt.Println("Wallet address will not be shown in client mode")
		}
		os.Exit(0)
	}

	var seedRPCServerAddr *nkngomobile.StringArray
	if len(opts.SeedRPCServerAddr) > 0 {
		seedRPCServerAddr = nkn.NewStringArray(opts.SeedRPCServerAddr...)
	}

	locations := make([]geo.Location, 0, len(opts.TunaCountry))
	for i := range opts.TunaCountry {
		countries := strings.Split(opts.TunaCountry[i], ",")
		l := make([]geo.Location, len(countries))
		for i := range countries {
			l[i].CountryCode = strings.TrimSpace(countries[i])
		}
		locations = append(locations, l...)
	}

	allowIps := make([]geo.Location, len(opts.TunaAllowIp))
	for i := range opts.TunaAllowIp {
		ips := strings.Split(opts.TunaAllowIp[i], ",")
		l := make([]geo.Location, len(ips))
		for i := range ips {
			l[i].IP = strings.TrimSpace(ips[i])
		}
		allowIps = append(allowIps, l...)
	}
	var allowedIP = append(locations, allowIps...)

	disallowedIP := make([]geo.Location, len(opts.TunaDisallowIp))
	for i := range opts.TunaDisallowIp {
		ips := strings.Split(opts.TunaDisallowIp[i], ",")
		l := make([]geo.Location, len(ips))
		for i := range ips {
			l[i].IP = strings.TrimSpace(ips[i])
		}
		disallowedIP = append(disallowedIP, l...)
	}

	allowedNknAddrs := make([]filter.NknClient, len(opts.TunaAllowNknAddr))
	for i := range opts.TunaAllowNknAddr {
		addrs := strings.Split(opts.TunaAllowNknAddr[i], ",")
		l := make([]filter.NknClient, len(addrs))
		for i := range addrs {
			l[i].Address = strings.TrimSpace(addrs[i])
		}
		allowedNknAddrs = append(allowedNknAddrs, l...)
	}

	disallowedNknAddrs := make([]filter.NknClient, len(opts.TunaDisallowNknAddr))
	for i := range opts.TunaDisallowNknAddr {
		addrs := strings.Split(opts.TunaDisallowNknAddr[i], ",")
		l := make([]filter.NknClient, len(addrs))
		for i := range addrs {
			l[i].Address = strings.TrimSpace(addrs[i])
		}
		disallowedNknAddrs = append(disallowedNknAddrs, l...)
	}

	clientConfig := &nkn.ClientConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	walletConfig := &nkn.WalletConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	dialConfig := &nkn.DialConfig{
		DialTimeout: opts.DialTimeout,
	}

	if util.IsValidUrl(opts.TunaMaxPrice) {
		price, err := util.GetRemotePrice(opts.TunaMaxPrice)
		if err != nil {
			log.Printf("Get remote price error: %v", err)
			price = config.FallbackTunaMaxPrice
		}
		log.Printf("Set dynamic price to %s", price)
		opts.TunaMaxPrice = price
	}

	tsConfig := &ts.Config{
		TunaMaxPrice:                 opts.TunaMaxPrice,
		TunaMinNanoPayFee:            opts.TunaMinFee,
		TunaNanoPayFeeRatio:          opts.TunaFeeRatio,
		TunaIPFilter:                 &geo.IPFilter{Allow: allowedIP, Disallow: disallowedIP},
		TunaNknFilter:                &filter.NknFilter{Allow: allowedNknAddrs, Disallow: disallowedNknAddrs},
		TunaServiceName:              opts.TunaServiceName,
		TunaDownloadGeoDB:            !opts.TunaDisableDownloadGeoDB,
		TunaGeoDBPath:                opts.TunaGeoDBPath,
		TunaMeasureBandwidth:         !opts.TunaDisableMeasureBandwidth,
		TunaMeasureStoragePath:       opts.TunaMeasureStoragePath,
		TunaMeasurementBytesDownLink: opts.TunaMeasureBandwidthBytes,
		TunaMinBalance:               opts.TunaMinBalance,
		Verbose:                      opts.Verbose,
	}
	if opts.Verbose {
		tsConfig.NumTunaListeners = 1
	}

	if opts.SessionWindowSize > 0 {
		clientConfig.SessionConfig = &ncp.Config{SessionWindowSize: opts.SessionWindowSize}
		tsConfig.SessionConfig = &ncp.Config{SessionWindowSize: opts.SessionWindowSize}
	}

	tunnelConfig := &tunnel.Config{
		AcceptAddrs:       nkn.NewStringArray(persistConf.AcceptAddrs...),
		ClientConfig:      clientConfig,
		WalletConfig:      walletConfig,
		DialConfig:        dialConfig,
		TunaSessionConfig: tsConfig,
		Verbose:           opts.Verbose,
		UDP:               opts.UDP,
		UDPIdleTime:       opts.UDPIdleTime,
	}

	ssClientConfig := ss.Config{
		TCP:      true,
		Cipher:   opts.Cipher,
		Password: opts.Password,

		Verbose:    opts.Verbose,
		UDPTimeout: config.DefaultUDPTimeout,
		UDP:        opts.UDP,

		TargetToClient: make(map[string]string),
	}

	if opts.UDP && opts.Client {
		ssClientConfig.UDPSocks = true
	}
	ssServerConfig := ssClientConfig

	nc := &nconnect{
		opts:           opts,
		account:        account,
		clientConfig:   clientConfig,
		tunnelConfig:   tunnelConfig,
		ssClientConfig: &ssClientConfig,
		ssServerConfig: &ssServerConfig,
		walletConfig:   walletConfig,
		persistConf:    persistConf,

		remoteInfoCache:    make(map[string]*admin.GetInfoJSON),
		remoteInfoByTunnel: make(map[string]*admin.GetInfoJSON),
		networkTunnels:     make(map[string]*tunnel.Tunnel),
		serverReady:        make(chan struct{}, 1),
	}

	return nc, nil
}

// Lazy create admin client to avoid unnecessary client creation.
func (nc *nconnect) getAdminClient() (*admin.Client, error) {
	if nc.adminClientCache != nil {
		return nc.adminClientCache, nil
	}

	identifier := ""
	if nc.opts.NetworkMember {
		if nc.opts.Identifier == "" {
			identifier = "member"
		} else {
			identifier = "member." + nc.opts.Identifier
		}
	}
	c, err := admin.NewClient(nc.account, nc.clientConfig, identifier)
	if err != nil {
		return nil, err
	}
	// Wait for more sub-clients to connect
	time.Sleep(time.Second)
	nc.adminClientCache = c

	return nc.adminClientCache, nil
}

// Lazy get remote info to avoid unnecessary rpc call.
func (nc *nconnect) getRemoteInfo(remoteAdminAddr string) (*admin.GetInfoJSON, error) {
	if info, ok := nc.remoteInfoCache[remoteAdminAddr]; ok {
		return info, nil
	}

	c, err := nc.getAdminClient()
	if err != nil {
		return nil, err
	}

	remoteInfoCache, err := c.GetInfo(remoteAdminAddr)
	if err != nil {
		return nil, fmt.Errorf("get remote server info error: %v. make sure server is online and accepting this client address", err)
	}

	nc.remoteInfoCache[remoteAdminAddr] = remoteInfoCache
	nc.remoteInfoByTunnel[remoteInfoCache.Addr] = remoteInfoCache

	return remoteInfoCache, nil
}

func (nc *nconnect) StartClient() error {
	if !nc.opts.NetworkMember {
		err := nc.opts.VerifyClient()
		if err != nil {
			return err
		}
	}

	remoteTunnelAddr := nc.opts.RemoteTunnelAddr
	if len(remoteTunnelAddr) == 0 {
		for _, remoteAdminAddr := range nc.opts.RemoteAdminAddr {
			remoteInfo, err := nc.getRemoteInfo(remoteAdminAddr)
			if err != nil {
				log.Printf("getRemoteInfo %v err: %v", remoteAdminAddr, err)
				continue
			}
			remoteTunnelAddr = append(remoteTunnelAddr, remoteInfo.Addr)
		}
	}
	if !nc.opts.NetworkMember && len(remoteTunnelAddr) == 0 {
		return fmt.Errorf("no remote tunnel address, start client fail")
	}

	vpnRoutes, err := nc.getRemoteRoutes()
	if err != nil {
		return err
	}

	if len(remoteTunnelAddr) > 0 {
		var from, to []string
		for _, remote := range remoteTunnelAddr {
			port, err := ts.GetFreePort(0)
			if err != nil {
				return err
			}

			ssAddr := "127.0.0.1:" + strconv.Itoa(port)
			from = append(from, ssAddr)
			to = append(to, remote)

			if remoteInfo, ok := nc.remoteInfoByTunnel[remote]; ok {
				for _, addr := range remoteInfo.LocalIP.Ipv4 {
					nc.ssClientConfig.TargetToClient[addr] = ssAddr
				}
			}
		}

		identifier := config.RandomIdentifier()
		tunnels, err := tunnel.NewTunnels(nc.account, identifier, from, to, nc.opts.Tuna, nc.tunnelConfig, nil)
		if err != nil {
			return err
		}
		nc.clientTunnels = tunnels

		nc.ssClientConfig.Client = from[0]
		nc.ssClientConfig.DefaultClient = from[0] // the first config is the default client
	} else {
		nc.ssClientConfig.Client = "127.0.0.1"
		nc.ssClientConfig.DefaultClient = ""
	}
	nc.ssClientConfig.Socks = nc.opts.LocalSocksAddr

	log.Println("nConnect socks proxy listen address:", nc.opts.LocalSocksAddr)

	if nc.opts.Tun || nc.opts.VPN {
		if !nc.opts.NetworkMember {
			err := arch.OpenTun(nc.opts.TunName, nc.opts.TunAddr, nc.opts.TunGateway, nc.opts.TunMask, nc.opts.TunDNS[0], nc.opts.LocalSocksAddr)
			if err != nil {
				log.Printf("OpenTun error: %v", err)
			} else {
				log.Println("Started tun2socks, interface:", nc.opts.TunName, "address:", nc.opts.TunAddr)
			}
		}

		if nc.opts.VPN {
			vpnCIDR, err := arch.SetVPNRoutes(nc.opts.TunName, nc.opts.TunGateway, vpnRoutes)
			if err != nil {
				return err
			}
			nc.routeCIDRs = vpnCIDR
			defer arch.RemoveVPNRoutes(nc.opts.TunName, nc.opts.TunGateway, nc.routeCIDRs)
		}
	}

	nc.startSSAndTunnel(true)
	nc.waitForSignal()

	return nil
}

func (nc *nconnect) StartServer() error {
	err := nc.opts.VerifyServer()
	if err != nil {
		return err
	}

	port, err := ts.GetFreePort(0)
	if err != nil {
		return err
	}
	ssAddr := "127.0.0.1:" + strconv.Itoa(port)
	nc.ssServerConfig.Server = ssAddr
	nc.ssServerConfig.Client = ""

	if nc.opts.Tuna {
		minBalance, err := common.StringToFixed64(nc.opts.TunaMinBalance)
		if err != nil {
			return err
		}

		if minBalance > 0 {
			w, err := nkn.NewWallet(nc.account, nc.walletConfig)
			if err != nil {
				return err
			}

			balance, err := w.Balance()
			if err != nil {
				log.Println("Fetch balance error:", err)
			} else if balance.ToFixed64() < minBalance {
				log.Printf("Wallet balance %s is less than minimal balance to enable tuna %s, tuna will not be enabled",
					balance.String(), nc.opts.TunaMinBalance)
				nc.opts.Tuna = false
			}
		}
	}

	if nc.tunaNode != nil {
		nc.tunnelConfig.TunaNode = nc.tunaNode
	}
	t, err := tunnel.NewTunnel(nc.account, nc.opts.Identifier, "", ssAddr, nc.opts.Tuna, nc.tunnelConfig, nil)
	if err != nil {
		return err
	}
	nc.serverTunnel = t
	log.Println("nConnect server tunnel listen address:", t.FromAddr())
	if nc.networkMember != nil {
		nc.networkMember.SetServerTunnel(t)
	}

	if len(nc.opts.AdminIdentifier) > 0 {
		go func() {
			identifier := nc.opts.AdminIdentifier
			if len(nc.opts.Identifier) > 0 {
				identifier += "." + nc.opts.Identifier
			}
			err := admin.StartNKNServer(nc.account, identifier, nc.clientConfig, t, nc.persistConf, &nc.opts.Config)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		log.Println("nConnect admin listening address:", nc.opts.AdminIdentifier+"."+t.FromAddr())
	}

	if len(nc.opts.AdminHTTPAddr) > 0 {
		go func() {
			err := admin.StartWebServer(nc.opts.AdminHTTPAddr, t, nc.persistConf, &nc.opts.Config)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		log.Println("nConnect admin web dashboard serve at:", nc.opts.AdminHTTPAddr)
	}

	nc.serverReady <- struct{}{}

	nc.startSSAndTunnel(false)
	nc.waitForSignal()

	return nil
}

func (nc *nconnect) startSSAndTunnel(client bool) {
	var ssConfig *ss.Config
	if client {
		ssConfig = nc.ssClientConfig
	} else {
		ssConfig = nc.ssServerConfig
	}
	go func() {
		err := ss.Start(ssConfig)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	if client {
		for _, t := range nc.clientTunnels {
			go func(t *tunnel.Tunnel) {
				err := t.Start()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}(t)
		}
	} else {
		go func() {
			err := nc.serverTunnel.Start()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
	}
}

func (nc *nconnect) waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	s := <-sigs
	log.Printf("Received signal '%v', exiting now...", s)
}

func (nc *nconnect) SetTunaNode(node *types.Node) {
	nc.tunaNode = node
}

func (nc *nconnect) GetClientTunnels() []*tunnel.Tunnel {
	nc.RLock()
	defer nc.RUnlock()
	iLen := len(nc.clientTunnels) + len(nc.networkTunnels)
	if iLen == 0 {
		return nil
	}
	tunnels := make([]*tunnel.Tunnel, 0, iLen)
	if len(nc.clientTunnels) > 0 {
		tunnels = append(tunnels, nc.clientTunnels...)
	}

	for _, t := range nc.networkTunnels {
		tunnels = append(tunnels, t)
	}
	return tunnels
}

func (nc *nconnect) StartNetworkManager() error {
	m, err := network.NewManager(nc.account, nc.clientConfig, nc.opts)
	if err != nil {
		return err
	}

	go func() {
		if err = m.StartManager(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	go func() {
		if err = m.StartWebServer(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	nc.waitForSignal()

	return nil
}

func (nc *nconnect) StartNetworkMember() error {
	if nc.opts.ManagerAddress == "" {
		return errors.New("network manager address is not specified")
	}

	mc, err := nc.getAdminClient()
	if err != nil {
		return err
	}
	nc.networkMember = network.NewMember(nc.opts, mc)

	if nc.opts.Server {
		go func() {
			err = nc.StartServer()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	serverAddr := ""
	if nc.opts.Server {
		<-nc.serverReady // wait server tunnel is ready
		if nc.serverTunnel != nil {
			serverAddr = nc.serverTunnel.FromAddr()
		}
	}

	nc.networkMember.CbNodeICanAccessUpdated = nc.setupNetworkTunnel
	go func() {
		err = nc.networkMember.StartMember(serverAddr)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if nc.opts.Client {
		go func() {
			err = nc.StartClient()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	// Start Cli Service
	go nc.networkMember.StartCliService()

	nc.waitForSignal()
	return nil
}

func (nc *nconnect) setupNetworkTunnel(nodes []*network.NodeInfo) error {
	if len(nodes) == 0 || !nc.opts.Client {
		return nil
	}

	oldTunnels := make(map[string]struct{})
	for addr := range nc.networkTunnels {
		oldTunnels[addr] = struct{}{}
	}

	var cidrs []*net.IPNet
	var from, to []string
	for _, node := range nodes {
		if node.ServerAddress == "" {
			continue
		}
		if _, ok := oldTunnels[node.ServerAddress]; ok {
			delete(oldTunnels, node.ServerAddress)
			continue
		}

		port, err := ts.GetFreePort(0)
		if err != nil {
			return err
		}
		ssAddr := "127.0.0.1:" + strconv.Itoa(port)

		toAddr := node.ServerAddress
		nc.ssClientConfig.TargetToClient[node.IP] = ssAddr

		from = append(from, ssAddr)
		to = append(to, toAddr)
		delete(oldTunnels, toAddr)

		_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/32", node.IP))
		if err != nil {
			continue
		}
		cidrs = append(cidrs, cidr)
	}

	var mc *nkn.MultiClient
	if len(nc.clientTunnels) > 0 {
		mc = nc.clientTunnels[0].MultiClient()
	}

	if len(from) > 0 {
		identifier := config.RandomIdentifier()

		tunnels, err := tunnel.NewTunnels(nc.account, identifier, from, to, nc.opts.Tuna, nc.tunnelConfig, mc)
		if err != nil {
			return err
		}

		if nc.ssClientConfig.DefaultClient == "" {
			nc.ssClientConfig.DefaultClient = from[0]
		}

		arch.SetVPNRoutes(nc.opts.TunName, nc.networkMember.GetNetworkInfo().Gateway, cidrs)
		ss.UpdateTargetToClient(nc.ssClientConfig.TargetToClient)

		for _, tunel := range tunnels {
			go func(t *tunnel.Tunnel) {
				log.Println("Connecting to tunnel:", t.ToAddr())
				err := t.Start()
				if err != nil {
					log.Printf("nconnect tunnel to %v start error: %v\n", t.ToAddr(), err)
				} else {
					if nc.opts.Verbose {
						log.Printf("nconnect tunnel to %v started\n", t.ToAddr())
					}
				}
			}(tunel)

			nc.Lock()
			nc.networkTunnels[tunel.ToAddr()] = tunel
			nc.Unlock()
		}
	}

	for addr := range oldTunnels {
		t := nc.networkTunnels[addr]
		t.Close()
		delete(nc.networkTunnels, addr)
	}

	return nil
}

func (nc *nconnect) getRemoteRoutes() ([]*net.IPNet, error) {
	vpnRoutes := nc.opts.VPNRoute
	if nc.opts.VPN && len(vpnRoutes) == 0 {
		for _, remoteAdminAddr := range nc.opts.RemoteAdminAddr {
			remoteInfo, err := nc.getRemoteInfo(remoteAdminAddr)
			if err != nil {
				log.Printf("getRemoteInfo %v err: %v", remoteAdminAddr, err)
				continue
			}
			if len(remoteInfo.LocalIP.Ipv4) > 0 {
				vpnRoutes = make([]string, 0, len(remoteInfo.LocalIP.Ipv4))
				for _, ip := range remoteInfo.LocalIP.Ipv4 {
					if ip == nc.opts.TunAddr || ip == nc.opts.TunGateway {
						log.Printf("Skipping server's local IP %s in routes", ip)
						continue
					}
					vpnRoutes = append(vpnRoutes, fmt.Sprintf("%s/32", ip))
				}
			}
		}
	}

	var routeCIDRs []*net.IPNet
	if len(vpnRoutes) > 0 {
		routeCIDRs = make([]*net.IPNet, len(vpnRoutes))
		for i, r := range vpnRoutes {
			_, cidr, err := net.ParseCIDR(r)
			if err != nil {
				return nil, fmt.Errorf("parse CIDR %s error: %v", r, err)
			}
			routeCIDRs[i] = cidr
		}
	}

	return routeCIDRs, nil
}
