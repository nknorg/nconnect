package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/network"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %+v", r)
		}
	}()

	var opts = &config.Opts{}
	_, err := flags.Parse(opts)
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

	if opts.Info != "" {
		cli(opts.Info)
		os.Exit(0)
	}

	nc, err := nconnect.NewNconnect(opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.NetworkManager {
		err = nc.StartNetworkManager()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if opts.NetworkMember {
		err = nc.StartNetworkMember()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if opts.Client {
		err = nc.StartClient()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if opts.Server {
		err = nc.StartServer()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

const help = `
nConnect -i <cmd>, to get nConnect information. The cmd can be:
help: 	this help
join: 	join network
leave: 	leave network
status: get network status
list: 	list nodes I can access and nodes which can access me
`

func cli(cmd string) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	switch cmd {
	case "help":
		fmt.Print(help)
	case "join":
		network.CliJoin()
	case "leave":
		network.CliLeave()
	case "status":
		network.CliStatus()
	case "list":
		network.CliList()
	default:
		fmt.Print("Unknown command: ", cmd, "\n", help)
	}
}
