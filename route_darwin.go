package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "-n", "add", "-net", dest.String(), gateway).Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway string) ([]byte, error) {
	return exec.Command("route", "-n", "delete", "-net", dest.String(), gateway).Output()
}
