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
	"github.com/nknorg/tuna"
)

var opts struct {
	config.Config
	ConfigFile string `short:"c" long:"config-file" default:"config.json" description:"Config file path"`
	Version    bool   `short:"v" long:"version" description:"Print version"`
}

var (
	Version string
)

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
		fmt.Println(Version)
		os.Exit(0)
	}

	conf, err := config.LoadOrNewConfig(opts.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	err = mergo.Merge(&opts.Config, conf, mergo.WithOverride)
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
		conf.Seed = hex.EncodeToString(account.Seed())
		shouldSave = true
	}

	if len(opts.Identifier) == 0 {
		conf.Identifier = config.RandomIdentifier()
		opts.Identifier = conf.Identifier
		shouldSave = true
	}

	if shouldSave {
		err = conf.Save()
		if err != nil {
			log.Fatal(err)
		}
	}

	var seedRPCServerAddr *nkn.StringArray
	if len(opts.SeedRPCServerAddr) > 0 {
		seedRPCServerAddr = nkn.NewStringArray(opts.SeedRPCServerAddr...)
	}

	tunaMaxPrice := opts.TunaMaxPrice
	if len(tunaMaxPrice) == 0 {
		tunaMaxPrice = config.DefaultTunaMaxPrice
	}

	locations := make([]tuna.Location, len(opts.TunaCountry))
	for i := range opts.TunaCountry {
		locations[i].CountryCode = strings.TrimSpace(opts.TunaCountry[i])
	}

	clientConfig := &nkn.ClientConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	walletConfig := &nkn.WalletConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	tsConfig := &ts.Config{
		TunaMaxPrice:    tunaMaxPrice,
		TunaIPFilter:    &tuna.IPFilter{Allow: locations},
		TunaServiceName: opts.TunaServiceName,
	}
	tunnelConfig := &tunnel.Config{
		AcceptAddrs:       nkn.NewStringArray(conf.AcceptAddrs...),
		ClientConfig:      clientConfig,
		WalletConfig:      walletConfig,
		TunaSessionConfig: tsConfig,
	}

	port, err := util.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	ssAddr := "127.0.0.1:" + strconv.Itoa(port)

	ssConfig := &ss.Config{
		TCP:      true,
		UDP:      false,
		Cipher:   "AEAD_CHACHA20_POLY1305",
		Password: opts.Password,
		Server:   ssAddr,
	}

	go func() {
		err := ss.Start(ssConfig)
		if err != nil {
			log.Println(err)
		}
		os.Exit(0)
	}()

	tun, err := tunnel.NewTunnel(account, opts.Identifier, "", ssAddr, opts.Tuna, tunnelConfig)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := tun.Start()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	log.Println("NKN tunnel listening address:", tun.FromAddr())

	if len(opts.AdminIdentifier) > 0 {
		go func() {
			identifier := opts.AdminIdentifier
			if len(opts.Identifier) > 0 {
				identifier += "." + opts.Identifier
			}
			err := admin.StartClient(account, identifier, clientConfig, tun, conf)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		log.Println("Admin client listening address:", opts.AdminIdentifier+"."+tun.FromAddr())
	}

	if len(opts.AdminHTTPAddr) > 0 {
		go func() {
			err := admin.StartWeb(opts.AdminHTTPAddr, tun, conf)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		log.Println("Admin web dashboard listening address:", opts.AdminHTTPAddr)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
