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

	name = tunDev.Name()
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, errors.New("invalid IP address")
	}

	var params string
	if isIPv4(ip) {
		params = fmt.Sprintf("%s inet %s netmask %s up", name, addr, mask)
	} else if isIPv6(ip) {
		prefixlen, err := strconv.Atoi(mask)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("parse IPv6 prefixlen failed: %v", err))
		}
		params = fmt.Sprintf("%s inet6 %s/%d up", name, addr, prefixlen)
	} else {
		return nil, errors.New("invalid IP address")
	}

	out, err := exec.Command("ifconfig", strings.Split(params, " ")...).Output()
	if err != nil {
		if len(out) != 0 {
			log.Println(out)
		}
		return nil, util.ParseExecError(err)
	}

	return tunDev, nil
}
