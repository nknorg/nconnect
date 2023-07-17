package arch

import (
	"io"
	"log"
	"os/exec"

	"github.com/eycorsican/go-tun2socks/tun"
)

const (
	tunComponentID = "tap0901"
)

var wintapdev io.ReadWriteCloser

func openTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	var err error
	wintapdev, err = tun.OpenTunDevice(name, addr, gw, mask, dnsServers, persist)
	return wintapdev, err
}

func SetTunIp(name, addr, mask, gw string) error {
	// if wintapdev != nil {
	// 	wintapdev.Close()
	// 	time.Sleep(2 * time.Second)
	// 	wintapdev = nil
	// }
	out, err := exec.Command("netsh", "interface", "ip", "add", "address", name, addr, mask).Output()
	log.Printf("SetTunIp: ip %s, mask %v, result: %s\n", addr, mask, string(out))
	// var err error
	// wintapdev, err = tun.OpenTunDevice(name, addr, gw, mask, []string{}, false)
	return err
}
