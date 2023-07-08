package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn/v2/vault"
	"github.com/nknorg/tuna"
	"github.com/nknorg/tuna/pb"
	"github.com/nknorg/tuna/types"
	"github.com/nknorg/tuna/util"
)

var ch chan string = make(chan string, 4)

func startNconnect(configFile string, tuna, udp, tun bool, n *types.Node) error {
	b, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("read config file %v err: %v", configFile, err)
		return err
	}
	var opts = &config.Opts{}
	err = json.Unmarshal(b, opts)
	if err != nil {
		log.Fatalf("parse config %v err: %v", configFile, err)
		return err
	}

	opts.Config.Tuna = tuna
	opts.Config.UDP = udp
	opts.Config.Tun = tun
	if tun {
		opts.Config.VPN = true
	}

	if opts.Client {
		port, err = getFreePort(port)
		if err != nil {
			return err
		}
		opts.LocalSocksAddr = fmt.Sprintf("127.0.0.1:%v", port)
	}
	fmt.Printf("opts.RemoteAdminAddr: %+v\n", opts.RemoteAdminAddr)

	nc, _ := nconnect.NewNconnect(opts)
	go func() {
		if opts.Server {
			nc.SetTunaNode(n)
			err = nc.StartServer()
			if err != nil {
				log.Fatalf("start nconnect server err: %v", err)
			}
		} else {
			err = nc.StartClient()
			if err != nil {
				log.Fatalf("start nconnect client err: %v", err)
			}
		}
	}()

	time.Sleep(5 * time.Second) // wait for nconnect to create tunnels

	tunnels := nc.GetTunnels()
	for _, tunnel := range tunnels {
		<-tunnel.TunaSessionClient().OnConnect()
	}

	return err
}

func getTunaNode() (*types.Node, error) {
	tunaSeed, _ := hex.DecodeString(seedHex)
	acc, err := nkn.NewAccount(tunaSeed)
	if err != nil {
		return nil, err
	}

	go runReverseEntry(tunaSeed)

	md := &pb.ServiceMetadata{
		Ip:              "127.0.0.1",
		TcpPort:         30020,
		UdpPort:         30021,
		ServiceId:       0,
		Price:           "0.0",
		BeneficiaryAddr: "",
	}
	n := &types.Node{
		Delay:       0,
		Bandwidth:   0,
		Metadata:    md,
		Address:     hex.EncodeToString(acc.PublicKey),
		MetadataRaw: "CgkxMjcuMC4wLjEQxOoBGMXqAToFMC4wMDE=",
	}

	return n, nil
}

func runReverseEntry(seed []byte) error {
	entryAccount, err := vault.NewAccountWithSeed(seed)
	if err != nil {
		return err
	}
	seedRPCServerAddr := nkn.NewStringArray(nkn.DefaultSeedRPCServerAddr...)

	walletConfig := &nkn.WalletConfig{
		SeedRPCServerAddr: seedRPCServerAddr,
	}
	entryWallet, err := nkn.NewWallet(&nkn.Account{Account: entryAccount}, walletConfig)
	if err != nil {
		return err
	}

	entryConfig := new(tuna.EntryConfiguration)
	err = util.ReadJSON("config.reverse.entry.json", entryConfig)
	if err != nil {
		return err
	}

	err = tuna.StartReverse(entryConfig, entryWallet)
	if err != nil {
		return err
	}

	ch <- tunaNodeStarted

	select {}
}

type Person struct {
	Name string
	Age  int
}

func getFreePort(p int) (int, error) {
	for i := 0; i < 100; i++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%v", p))
		if err != nil {
			return 0, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			p++
			continue
		}

		defer l.Close()

		return l.Addr().(*net.TCPAddr).Port, nil
	}
	return 0, fmt.Errorf("can't find free port")
}

func waitForSSProxReady() error {
	for i := 0; i < 100; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if conn != nil {
			conn.Close()
			return nil
		}
	}
	return fmt.Errorf("ss is not ready after 200 seconds, give up")
}
