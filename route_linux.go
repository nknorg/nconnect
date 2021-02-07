package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "-n", "add", dest.String(), "gw", gateway).Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "-n", "del", dest.String(), "gw", gateway).Output()
}
