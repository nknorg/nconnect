package tests

import (
	"fmt"
	"testing"
	"time"
)

// go test -v -run=TestUDPByTun
func TestUDPByTun(t *testing.T) {
	tunaNode, err := getTunaNode()
	if err != nil {
		fmt.Printf("getTunaNode err %v\n", err)
		return
	}

	go StartUdpServer()

	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("server.json", tuna, udp, tun, tunaNode)
		if err != nil {
			fmt.Printf("start nconnect server err: %v\n", err)
			return
		}
	}()

	time.Sleep(15 * time.Second)

	tun = true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()

	time.Sleep(15 * time.Second)

	go StartUDPClient()

	waitFor(ch, udpClientExited)
}
