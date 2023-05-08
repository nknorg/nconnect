package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/txthinking/brook"
)

// go test -v -run=TestDNSByTun
func TestDNSByTun(t *testing.T) {
	go StartNconnectServerWithTunaNode(true, true, false)
	time.Sleep(15 * time.Second)

	go DnsByTun()
	waitFor(ch, exited)
}

func DnsByTun() {
	tuna, udp, tun := true, true, true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(15 * time.Second)

	for i := 1; i <= 10; i++ {
		err := brook.Socks5Test(proxyAddr, "", "", "http3.ooo", "137.184.237.95", "8.8.8.8:53")
		if err != nil {
			fmt.Printf("TestDNSProxy try %v err: %v\n", i, err)
			time.Sleep(time.Duration(i) * time.Second)
			break
		}
	}

	ch <- udpClientExited
}
