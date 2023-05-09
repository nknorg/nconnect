package tests

import (
	"fmt"
	"testing"
	"time"
)

// go test -v -run=TestProxy
func TestProxy(t *testing.T) {
	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(15 * time.Second)

	dnsQuery()
	StartWebClient()
	StartTCPClient()
	StartUDPClient()
}
