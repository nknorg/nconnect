// +build !linux,!darwin

package ss

import (
	"errors"
	"net"
	"time"
)

func redirLocal(addr, server string, shadow func(net.Conn) net.Conn) error {
	return errors.New("TCP redirect not supported")
}

func redir6Local(addr, server string, shadow func(net.Conn) net.Conn) error {
	return errors.New("TCP6 redirect not supported")
}

func timedCork(c *net.TCPConn, d time.Duration) error { return nil }
