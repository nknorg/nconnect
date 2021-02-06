package main

import (
	"net"
	"os/exec"
	"strconv"
	"strings"
)

func maskToString(m net.IPMask) string {
	s := make([]string, len(m))
	for i, b := range m {
		s[i] = strconv.Itoa(int(b))
	}
	return strings.Join(s, ".")
}

func addRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "add", dest.IP.String(), "MASK", maskToString(dest.Mask), gateway).Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "delete", dest.IP.String()).Output()
}
