package main

import (
	"net"
	"os/exec"
)

func addRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	out, err := exec.Command("ip", "route", "add", dest.String(), "via", gateway, "dev", devName).Output()
	if err == nil {
		return out, nil
	}
	out, err = exec.Command("ip", "route", "change", dest.String(), "via", gateway, "dev", devName).Output()
	if err == nil {
		return out, nil
	}
	out, err = exec.Command("route", "-n", "add", dest.String(), "gw", gateway).Output()
	if err == nil {
		return out, nil
	}
	return exec.Command("route", "-n", "change", dest.String(), "gw", gateway).Output()
}

func deleteRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	out, err := exec.Command("ip", "route", "del", dest.String(), "via", gateway, "dev", devName).Output()
	if err == nil {
		return out, nil
	}
	return exec.Command("route", "-n", "del", dest.String(), "gw", gateway).Output()
}
