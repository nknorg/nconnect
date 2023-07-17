package arch

import "testing"

func TestOpenTunDevice(t *testing.T) {
	name := "tap0901"
	addr := "192.168.0.2"
	gw := "192.168.0.1"
	mask := "255.255.255.0"
	dnsServers := []string{"192.168.0.1"}
	persist := false

	openTunDevice(name, addr, gw, mask, dnsServers, persist)
}
