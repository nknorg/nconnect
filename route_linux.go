package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("ip", "route", "add", dest.String(), "via", gateway).Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("ip", "route", "delete", dest.String(), "via", gateway).Output()
}
