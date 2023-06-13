package tests

import (
	"fmt"
	"testing"
	"time"
)

// go test -v -run=TestTun
func TestTun(t *testing.T) {
	tuna, udp, tun := true, true, true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(20 * time.Second)

	dnsQuery()
	for _, server := range servers {
		StartTunWebClient("http://" + server + httpPort + "/httpEcho")
		StartTCPClient(server + tcpPort)
		StartUDPClient(server + udpPort)
	}
}
