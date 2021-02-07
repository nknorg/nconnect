package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eycorsican/go-tun2socks/tun"
	"github.com/nknorg/nconnect/util"
)

func OpenTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	tunDev, err := tun.OpenTunDevice(name, addr, gw, mask, dnsServers, persist)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	var params string
	if ip.To4() != nil {
		params = fmt.Sprintf("%s inet %s netmask %s up", name, addr, mask)
	} else {
		prefixlen, err := strconv.Atoi(mask)
		if err != nil {
			return nil, fmt.Errorf("parse IPv6 prefixlen failed: %v", err)
		}
		params = fmt.Sprintf("%s inet6 %s/%d up", name, addr, prefixlen)
	}

	out, err := exec.Command("ifconfig", strings.Split(params, " ")...).Output()
	if err != nil {
		if len(out) > 0 {
			log.Print(string(out))
		}
		return nil, util.ParseExecError(err)
	}

	return tunDev, nil
}
