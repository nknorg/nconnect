package main

import (
	"io"

	"github.com/eycorsican/go-tun2socks/tun"
)

const (
	tunComponentID = "tap0901"
)

func OpenTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	return tun.OpenTunDevice(name, addr, gw, mask, dnsServers, persist)
}
