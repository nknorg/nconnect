package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/imdario/mergo"
	"github.com/jessevdk/go-flags"
	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/ss"
	"github.com/nknorg/nconnect/util"
	"github.com/nknorg/nkn-sdk-go"
	ts "github.com/nknorg/nkn-tuna-session"
	tunnel "github.com/nknorg/nkn-tunnel"
	"github.com/nknorg/tuna/geo"
)

var opts struct {
	Client bool `short:"c" long:"client" description:"Client mode"`
	Server bool `short:"s" long:"server" description:"Server mode"`

	config.Config
	ConfigFile string `short:"f" long:"config-file" default:"config.json" description:"Config file path"`

	Version bool `long:"version" description:"Print version"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %+v", r)
		}
	}()

	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	if opts.Version {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	if opts.Client == opts.Server {
		log.Fatal("Exactly one mode (client or server) should be selected.")
	}

	persistConf, err := config.LoadOrNewConfig(opts.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	err = mergo.Merge(&opts.Config, persistConf)
	if err != nil {
		log.Fatal(err)
	}

	seed, err := hex.DecodeString(opts.Seed)
	if err != nil {
		log.Fatal(err)
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
	}

	var seedRPCServerAddr *nkn.StringArray
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

	clientConfig := &nkn.ClientConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	walletConfig := &nkn.WalletConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	tsConfig := &ts.Config{
		TunaMaxPrice:         opts.TunaMaxPrice,
		TunaIPFilter:         &geo.IPFilter{Allow: locations},
		TunaServiceName:      opts.TunaServiceName,
		TunaDownloadGeoDB:    !opts.TunaDisableDownloadGeoDB,
		TunaGeoDBPath:        opts.TunaGeoDBPath,
		TunaMeasureBandwidth: !opts.TunaDisableMeasureBandwidth,
	}
	tunnelConfig := &tunnel.Config{
		AcceptAddrs:       nkn.NewStringArray(persistConf.AcceptAddrs...),
		ClientConfig:      clientConfig,
		WalletConfig:      walletConfig,
		TunaSessionConfig: tsConfig,
		Verbose:           opts.Verbose,
	}

	port, err := util.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	ssAddr := "127.0.0.1:" + strconv.Itoa(port)

	ssConfig := &ss.Config{
		TCP:      true,
		UDP:      false,
		UDPSocks: true,
		Cipher:   opts.Cipher,
		Password: opts.Password,
		Verbose:  opts.Verbose,
	}

	var tun *tunnel.Tunnel

	if opts.Client {
		if len(opts.RemoteAddr) == 0 {
			log.Fatal("Remote address should not be empty.")
		}

		ssConfig.Client = ssAddr
		ssConfig.Socks = opts.LocalAddr

		tun, err = tunnel.NewTunnel(account, opts.Identifier, ssAddr, opts.RemoteAddr, opts.Tuna, tunnelConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Client NKN address:", tun.Addr().String())
		log.Println("Client socks proxy listen address:", opts.LocalAddr)
	}

	if opts.Server {
		ssConfig.Server = ssAddr

		tun, err = tunnel.NewTunnel(account, opts.Identifier, "", ssAddr, opts.Tuna, tunnelConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Server listen address:", tun.FromAddr())

		if len(opts.AdminIdentifier) > 0 {
			go func() {
				identifier := opts.AdminIdentifier
				if len(opts.Identifier) > 0 {
					identifier += "." + opts.Identifier
				}
				err := admin.StartClient(account, identifier, clientConfig, tun, persistConf, &opts.Config)
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}()
			log.Println("Admin client listening address:", opts.AdminIdentifier+"."+tun.FromAddr())
		}

		if len(opts.AdminHTTPAddr) > 0 {
			go func() {
				err := admin.StartWeb(opts.AdminHTTPAddr, tun, persistConf, &opts.Config)
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}()
			log.Println("Admin web dashboard listening address:", opts.AdminHTTPAddr)
		}
	}

	go func() {
		err := ss.Start(ssConfig)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	go func() {
		err := tun.Start()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
