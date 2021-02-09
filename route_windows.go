package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	b, err := exec.Command("netsh", "interface", "ipv4", "add", "route", dest.String(), "nexthop="+gateway, "interface="+devName, "metric=0", "store=active").Output()
	if err == nil {
		return b, nil
	}
	return exec.Command("netsh", "interface", "ipv4", "set", "route", dest.String(), "nexthop="+gateway, "interface="+devName, "metric=0", "store=active").Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	return exec.Command("netsh", "interface", "ipv4", "delete", "route", dest.String(), "interface="+devName).Output()
}
