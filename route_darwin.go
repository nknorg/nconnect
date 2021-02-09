package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	b, err := exec.Command("route", "-n", "add", "-net", dest.String(), gateway).Output()
	if err != nil {
		return exec.Command("route", "-n", "change", "-net", dest.String(), gateway).Output()
	}
	return b, nil
}

func deleteRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	return exec.Command("route", "-n", "delete", "-net", dest.String(), gateway).Output()
}
