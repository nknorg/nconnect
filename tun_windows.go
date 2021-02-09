package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

const (
	tunComponentID = "tap0901"
)

func OpenTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	cmd := exec.Command("netsh", "interface", "ip", "set", "address", name, "static", addr, mask)
	cmd.Run()
	cmd = exec.Command("netsh", "interface", "ip", "set", "dns", name, "static", dnsServers[0])
	cmd.Run()
	if len(dnsServers) >= 2 {
		cmd = exec.Command("netsh", "interface", "ip", "add", "dns", name, dnsServers[1], "index=2")
		cmd.Run()
	}
	prefix, _ := net.IPMask(net.ParseIP(mask).To4()).Size()
	network := fmt.Sprintf("%v/%v", addr, prefix)
	return water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			ComponentID:   tunComponentID,
			InterfaceName: name,
			Network:       network,
		},
	})
}
