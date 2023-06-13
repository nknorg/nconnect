package nconnect

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/imdario/mergo"
	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/arch"
	"github.com/nknorg/nconnect/config"
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

const (
	mtu = 1500
)

type nconnect struct {
	opts    *config.Opts
	account *nkn.Account

	walletConfig *nkn.WalletConfig
	clientConfig *nkn.ClientConfig
	tunnelConfig *tunnel.Config
	ssConfig     *ss.Config
	persistConf  *config.Config

	adminClientCache   *admin.Client
	remoteInfoCache    map[string]*admin.GetInfoJSON // map remote admin address to remote info
	remoteInfoByTunnel map[string]*admin.GetInfoJSON // map tunnel address to remote info

	tunnels  []*tunnel.Tunnel
	tunaNode *types.Node // It is used to connect specified tuna node, mainly is for testing.
}

func NewNconnect(opts *config.Opts) (*nconnect, error) {
	err := (&opts.Config).SetPlatformSpecificDefaultValues()
	if err != nil {
		return nil, err
	}

	if opts.Client == opts.Server {
		log.Fatal("Exactly one mode (client or server) should be selected.")
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
		TunaMaxPrice:           opts.TunaMaxPrice,
		TunaMinNanoPayFee:      opts.TunaMinFee,
		TunaNanoPayFeeRatio:    opts.TunaFeeRatio,
		TunaIPFilter:           &geo.IPFilter{Allow: allowedIP, Disallow: disallowedIP},
		TunaNknFilter:          &filter.NknFilter{Allow: allowedNknAddrs, Disallow: disallowedNknAddrs},
		TunaServiceName:        opts.TunaServiceName,
		TunaDownloadGeoDB:      !opts.TunaDisableDownloadGeoDB,
		TunaGeoDBPath:          opts.TunaGeoDBPath,
		TunaMeasureBandwidth:   !opts.TunaDisableMeasureBandwidth,
		TunaMeasureStoragePath: opts.TunaMeasureStoragePath,
		TunaMinBalance:         opts.TunaMinBalance,
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

	ssConfig := &ss.Config{
		TCP:      true,
		Cipher:   opts.Cipher,
		Password: opts.Password,

		Verbose:    opts.Verbose,
		UDPTimeout: config.DefaultUDPTimeout,
		UDP:        opts.UDP,

		TargetToClient: make(map[string]string),
	}

	if opts.UDP && opts.Client {
		ssConfig.UDPSocks = true
	}

	nc := &nconnect{
		opts:         opts,
		account:      account,
		clientConfig: clientConfig,
		tunnelConfig: tunnelConfig,
		ssConfig:     ssConfig,
		walletConfig: walletConfig,
		persistConf:  persistConf,

		remoteInfoCache:    make(map[string]*admin.GetInfoJSON),
		remoteInfoByTunnel: make(map[string]*admin.GetInfoJSON),
	}

	return nc, nil
}

// Lazy create admin client to avoid unnecessary client creation.
func (nc *nconnect) getAdminClient() (*admin.Client, error) {
	if nc.adminClientCache != nil {
		return nc.adminClientCache, nil
	}
	c, err := admin.NewClient(nc.account, nc.clientConfig)
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
		return nil, fmt.Errorf("get remote server info error: %v. make sure server is online and accepting connections from this client address", err)
	}

	nc.remoteInfoCache[remoteAdminAddr] = remoteInfoCache
	nc.remoteInfoByTunnel[remoteInfoCache.Addr] = remoteInfoCache

	return remoteInfoCache, nil
}

func (nc *nconnect) StartClient() error {
	err := nc.opts.VerifyClient()
	if err != nil {
		return err
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
	if len(remoteTunnelAddr) == 0 {
		return fmt.Errorf("no remote tunnel address, start client fail")
	}

	var vpnCIDR []*net.IPNet
	if nc.opts.VPN {
		vpnRoutes := nc.opts.VPNRoute
		if len(vpnRoutes) == 0 {
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
		if len(vpnRoutes) > 0 {
			vpnCIDR = make([]*net.IPNet, len(vpnRoutes))
			for i, cidr := range vpnRoutes {
				_, cidr, err := net.ParseCIDR(cidr)
				if err != nil {
					return fmt.Errorf("parse CIDR %s error: %v", cidr, err)
				}
				vpnCIDR[i] = cidr
			}
		}
	}

	proxyAddr, err := net.ResolveTCPAddr("tcp", nc.opts.LocalSocksAddr)
	if err != nil {
		return fmt.Errorf("invalid proxy server address: %v", err)
	}
	proxyHost := proxyAddr.IP.String()
	proxyPort := uint16(proxyAddr.Port)

	var from, to []string
	for _, remote := range remoteTunnelAddr {
		port, err := util.GetFreePort()
		if err != nil {
			return err
		}

		ssAddr := "127.0.0.1:" + strconv.Itoa(port)
		from = append(from, ssAddr)
		to = append(to, remote)

		if remoteInfo, ok := nc.remoteInfoByTunnel[remote]; ok {
			for _, addr := range remoteInfo.LocalIP.Ipv4 {
				nc.ssConfig.TargetToClient[addr] = ssAddr
			}
		}
	}
	tunnels, err := tunnel.NewTunnels(nc.account, nc.opts.Identifier, from, to, nc.opts.Tuna, nc.tunnelConfig)
	if err != nil {
		return err
	}
	nc.tunnels = tunnels

	nc.ssConfig.Socks = nc.opts.LocalSocksAddr
	nc.ssConfig.Client = from[0]
	nc.ssConfig.DefaultClient = from[0] // the first config is the default client

	log.Println("Client socks proxy listen address:", nc.opts.LocalSocksAddr)

	if nc.opts.Tun || nc.opts.VPN {
		tunDevice, err := arch.OpenTunDevice(nc.opts.TunName, nc.opts.TunAddr, nc.opts.TunGateway, nc.opts.TunMask, nc.opts.TunDNS, true)
		if err != nil {
			return fmt.Errorf("failed to open TUN device: %v", err)
		}

		core.RegisterOutputFn(tunDevice.Write)

		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30*time.Second))

		lwipWriter := core.NewLWIPStack()

		go func() {
			_, err := io.CopyBuffer(lwipWriter, tunDevice, make([]byte, mtu))
			if err != nil {
				log.Fatalf("Failed to write data to network stack: %v", err)
			}
		}()

		log.Println("Started tun2socks")

		if nc.opts.VPN {
			for _, dest := range vpnCIDR {
				log.Printf("Adding route %s", dest)
				out, err := arch.AddRouteCmd(dest, nc.opts.TunGateway, nc.opts.TunName)
				if len(out) > 0 {
					os.Stdout.Write(out)
				}
				if err != nil {
					os.Stdout.Write([]byte(util.ParseExecError(err)))
					os.Exit(1)
				}
				defer func(dest *net.IPNet) {
					log.Printf("Deleting route %s", dest)
					out, err := arch.DeleteRouteCmd(dest, nc.opts.TunGateway, nc.opts.TunName)
					if len(out) > 0 {
						os.Stdout.Write(out)
					}
					if err != nil {
						os.Stdout.Write([]byte(util.ParseExecError(err)))
					}
				}(dest)
			}
		}
	}

	nc.startSSAndTunnel()
	nc.waitForSignal()

	return nil
}

func (nc *nconnect) StartServer() error {
	err := nc.opts.VerifyServer()
	if err != nil {
		return err
	}

	port, err := util.GetFreePort()
	if err != nil {
		return err
	}
	ssAddr := "127.0.0.1:" + strconv.Itoa(port)
	nc.ssConfig.Server = ssAddr

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
	t, err := tunnel.NewTunnel(nc.account, nc.opts.Identifier, "", ssAddr, nc.opts.Tuna, nc.tunnelConfig)
	if err != nil {
		return err
	}
	nc.tunnels = append(nc.tunnels, t)
	log.Println("Tunnel listen address:", t.FromAddr())

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
		log.Println("Admin listening address:", nc.opts.AdminIdentifier+"."+t.FromAddr())
	}

	if len(nc.opts.AdminHTTPAddr) > 0 {
		go func() {
			err := admin.StartWebServer(nc.opts.AdminHTTPAddr, t, nc.persistConf, &nc.opts.Config)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		log.Println("Admin web dashboard listening address:", nc.opts.AdminHTTPAddr)
	}

	nc.startSSAndTunnel()
	nc.waitForSignal()

	return nil
}

func (nc *nconnect) startSSAndTunnel() {
	go func() {
		err := ss.Start(nc.ssConfig)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	for _, t := range nc.tunnels {
		go func(t *tunnel.Tunnel) {
			err := t.Start()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}(t)
	}
}

func (nc *nconnect) waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func (nc *nconnect) SetTunaNode(node *types.Node) {
	nc.tunaNode = node
}
