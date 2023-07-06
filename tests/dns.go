package tests

import (
	"fmt"

	"github.com/txthinking/brook"
)

func dnsQuery() error {
	proxyAddr := fmt.Sprintf("127.0.0.1:%v", port)
	for i := 1; i <= numMsgs; i++ {
		err := brook.Socks5Test(proxyAddr, "", "", "http3.ooo", "137.184.237.95", "8.8.8.8:53")
		if err != nil {
			fmt.Printf("TestDNSProxy try %v err: %v\n", i, err)
			return err
		}
	}
	return nil
}
