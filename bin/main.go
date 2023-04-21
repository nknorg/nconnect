package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
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

	nc, err := nconnect.NewNconnect(opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Client {
		err = nc.StartClient()
		if err != nil {
			log.Fatal(err)
		}
	}
	if opts.Server {
		err = nc.StartServer()
		if err != nil {
			log.Fatal(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
