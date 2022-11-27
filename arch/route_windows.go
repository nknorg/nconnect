package arch

import (
	"net"
	"os/exec"
)

func AddRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	out, err := exec.Command("netsh", "interface", "ipv4", "add", "route", dest.String(), "nexthop="+gateway, "interface="+devName, "metric=0", "store=active").Output()
	if err == nil {
		return out, nil
	}
	return exec.Command("netsh", "interface", "ipv4", "set", "route", dest.String(), "nexthop="+gateway, "interface="+devName, "metric=0", "store=active").Output()
}

func DeleteRouteCmd(dest *net.IPNet, gateway, devName string) ([]byte, error) {
	return exec.Command("netsh", "interface", "ipv4", "delete", "route", dest.String(), "interface="+devName).Output()
}
