package arch

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

func openTunDevice(name, addr, gw, mask string, dnsServers []string, persist bool) (io.ReadWriteCloser, error) {
	tunDev, err := tun.OpenTunDevice(name, addr, gw, mask, dnsServers, persist)
	if err != nil {
		return nil, err
	}

	err = SetTunIp(name, addr, mask, gw)

	return tunDev, err
}

func SetTunIp(tapName, ip, mask, gw string) error {
	out, err := func() ([]byte, error) {
		out, err := exec.Command("ip", "addr", "replace", ip+"/"+mask, "dev", tapName).Output()
		if err != nil {
			return out, err
		}
		return exec.Command("ip", "link", "set", "dev", tapName, "up").Output()
	}()
	if err != nil {
		if len(out) > 0 {
			log.Print(string(out))
		}
		log.Println(util.ParseExecError(err))

		ip := net.ParseIP(ip)
		if ip == nil {
			return errors.New("invalid IP address")
		}

		var params string
		if ip.To4() != nil {
			params = fmt.Sprintf("%s inet %s netmask %s up", tapName, ip, mask)
		} else {
			prefixlen, err := strconv.Atoi(mask)
			if err != nil {
				return fmt.Errorf("parse IPv6 prefixlen failed: %v", err)
			}
			params = fmt.Sprintf("%s inet6 %s/%d up", tapName, ip, prefixlen)
		}

		out, err := exec.Command("ifconfig", strings.Split(params, " ")...).Output()
		if err != nil {
			if len(out) > 0 {
				log.Print(string(out))
			}
			return errors.New(util.ParseExecError(err))
		}
	}
	return nil
}
