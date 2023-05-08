package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/txthinking/brook"
)

// go test -v -run=TestDNSByProxy
func TestDNSByProxy(t *testing.T) {
	go StartNconnectServerWithTunaNode(true, true, false)
	time.Sleep(15 * time.Second)

	go DnsByProxy()
	waitFor(ch, exited)
}

func DnsByProxy() {
	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(20 * time.Second)

	for i := 1; i <= rounds; i++ {
		err := brook.Socks5Test(proxyAddr, "", "", "http3.ooo", "137.184.237.95", "8.8.8.8:53")
		if err != nil {
			fmt.Printf("TestDNSProxy try %v err: %v\n", i, err)
			time.Sleep(time.Duration(i) * time.Second)
			break
		}
	}

	ch <- udpClientExited
}
