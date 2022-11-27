package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
	"log"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %+v", r)
		}
	}()

	var opts config.NConfig
	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	nconnect.Run(&opts)
}
