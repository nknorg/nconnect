package arch

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"

	"github.com/eycorsican/go-tun2socks/tun"
	"github.com/songgao/water"
)

var tundev *water.Interface

func openTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	rwc, err := tun.OpenTunDevice(name, addr, gw, mask, dnsServers, persist)
	if err == nil {
		tundev = rwc.(*water.Interface)
	}
	return rwc, err
}

func SetTunIp(tunName, addr, mask, gw string) error {

	var params string
	params = fmt.Sprintf("lo0 alias %v %v ", addr, mask)

	out, err := exec.Command("ifconfig", strings.Split(params, " ")...).Output()
	if err != nil {
		if len(out) != 0 {
			return errors.New(fmt.Sprintf("%v, output: %s", err, out))
		}
		return err
	}

	return nil
}

func SetTunIp_old(tunName, addr, mask, gw string) error {
	ip := net.ParseIP(addr)
	if ip == nil {
		return errors.New("invalid IP address")
	}

	if tundev == nil {
		return errors.New("tun device is not open")
	}

	tunName = tundev.Name()

	var params string
	params = fmt.Sprintf("%s inet %s netmask %s %s", tunName, addr, mask, gw)

	out, err := exec.Command("ifconfig", strings.Split(params, " ")...).Output()
	if err != nil {
		if len(out) != 0 {
			return errors.New(fmt.Sprintf("%v, output: %s", err, out))
		}
		return err
	}

	return nil
}
