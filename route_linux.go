package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	b, err := exec.Command("route", "-n", "add", dest.String(), "gw", gateway).Output()
	if err != nil {
		return exec.Command("route", "-n", "change", dest.String(), "gw", gateway).Output()
	}
	return b, nil
}

func deleteRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	return exec.Command("route", "-n", "del", dest.String(), "gw", gateway).Output()
}
