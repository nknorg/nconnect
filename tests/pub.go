package tests

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn/v2/vault"
	"github.com/nknorg/tuna"
	"github.com/nknorg/tuna/pb"
	"github.com/nknorg/tuna/types"
	"github.com/nknorg/tuna/util"
)

var ch chan string = make(chan string)

func startNconnect(configFile string, n *types.Node) error {
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

	nc, _ := nconnect.NewNconnect(opts)
	if opts.Server {
		nc.SetTunaNode(n)
		nc.StartServer()
	} else {
		nc.StartClient()
	}

	return nil
}

func startTunaNode() (*types.Node, error) {
	tunaSeed, _ := hex.DecodeString(seedHex)
	acc, err := nkn.NewAccount(tunaSeed)
	if err != nil {
		return nil, err
	}

	go runReverseEntry(tunaSeed)

	n := &types.Node{
		Delay:     0,
		Bandwidth: 0,
		Metadata: &pb.ServiceMetadata{
			Ip:              "127.0.0.1",
			TcpPort:         30020,
			UdpPort:         30021,
			ServiceId:       0,
			Price:           "0.0",
			BeneficiaryAddr: "",
		},
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
