package main

import (
	"flag"
	"log"

	"github.com/nknorg/nconnect/tests"
)

const (
	port = ":12345"
)

func main() {
	var udp = flag.Bool("udp", false, "udp mode")
	var server = flag.Bool("server", false, "run as server")
	var serverAddr = flag.String("serverAddr", "127.0.0.1", "server's ip")
	flag.Parse()

	if *server { // server, both tcp and udp
		go func() {
			err := tests.StartTCPServer(port)
			if err != nil {
				log.Printf("StartTCPServer err: %v\n", err)
			}
		}()

		err := tests.StartUDPServer(port)
		if err != nil {
			log.Printf("StartUDPServer err: %v\n", err)
		}
	} else { // client
		if *udp { // udp client
			err := tests.StartUDPTunClient(*serverAddr + port)
			if err != nil {
				log.Printf("StartUDPTunClient err: %v\n", err)
			}
		} else { // tcp client
			err := tests.StartTCPTunClient(*serverAddr + port)
			if err != nil {
				log.Printf("StartTCPTunClient err: %v\n", err)
			}
		}
	}
}
