package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn-socks/admin"
	"github.com/nknorg/nkn-socks/config"
	"github.com/nknorg/nkn-socks/ss"
	"github.com/nknorg/nkn-socks/util"
	ts "github.com/nknorg/nkn-tuna-session"
	tunnel "github.com/nknorg/nkn-tunnel"
)

var opts struct {
	AdminHTTPAddr   string `long:"admin-http" description:"Admin web GUI listen address (e.g. 127.0.0.1:8000)"`
	AdminIdentifier string `long:"admin-identifier" description:"Admin NKN client identifier prefix"`
	ConfigFile      string `short:"c" long:"config-file" default:"config.json" description:"Config file path"`
	Version         bool   `short:"v" long:"version" description:"Print version"`
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

	seed, err := hex.DecodeString(conf.Seed)
	if err != nil {
		log.Fatal(err)
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Seed != hex.EncodeToString(account.Seed()) {
		conf.Seed = hex.EncodeToString(account.Seed())
		err = conf.Save()
		if err != nil {
			log.Fatal(err)
		}
	}

	var seedRPCServerAddr *nkn.StringArray
	if len(conf.SeedRPCServerAddr) > 0 {
		seedRPCServerAddr = nkn.NewStringArray(conf.SeedRPCServerAddr...)
	}

	tunaMaxPrice := conf.TunaMaxPrice
	if len(tunaMaxPrice) == 0 {
		tunaMaxPrice = config.DefaultTunaMaxPrice
	}

	clientConfig := &nkn.ClientConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	walletConfig := &nkn.WalletConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	tsConfig := &ts.Config{
		TunaMaxPrice: tunaMaxPrice,
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
		Password: conf.Password,
		Server:   ssAddr,
	}

	go func() {
		err := ss.Start(ssConfig)
		if err != nil {
			log.Println(err)
		}
		os.Exit(0)
	}()

	tun, err := tunnel.NewTunnel(account, conf.Identifier, "", ssAddr, true, tunnelConfig)
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
			if len(conf.Identifier) > 0 {
				identifier += "." + conf.Identifier
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
