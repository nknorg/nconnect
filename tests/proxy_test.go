package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// go test -v -run=TestProxy
func TestProxy(t *testing.T) {
	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		require.NoError(t, err)
	}()

	time.Sleep(5 * time.Second)

	err := waitSSAndTunaReady()
	require.NoError(t, err)

	err = dnsQuery()
	require.NoError(t, err)
	for _, server := range servers {
		err := StartWebClient("http://" + server + httpPort + "/httpEcho")
		require.NoError(t, err)
		err = StartTCPClient(server + tcpPort)
		require.NoError(t, err)
		err = StartUDPClient(server + udpPort)
		require.NoError(t, err)
	}
}
