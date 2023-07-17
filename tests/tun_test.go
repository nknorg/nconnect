package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// var tun = flag.Bool("tun", false, "use tun device")

// go test -v -run=TestTun -tun
func TestTun(t *testing.T) {
	fmt.Println("Make sure run this case at root or administrator shell")

	if !(*tun) {
		t.Skip("Skip TestTun, if you want to run this test, please use: go test -v -tun .")
	}

	tuna, udp, tun := true, true, true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		require.NoError(t, err)
	}()

	time.Sleep(10 * time.Second)

	err := waitForSSProxReady()
	require.NoError(t, err)

	err = dnsQuery()
	require.NoError(t, err)
	for _, server := range servers {
		err := StartTunWebClient("http://" + server + httpPort + "/httpEcho")
		require.NoError(t, err)
		err = StartTCPClient(server + tcpPort)
		require.NoError(t, err)
		err = StartUDPClient(server + udpPort)
		require.NoError(t, err)
	}
}
