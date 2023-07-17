package ss

import (
	"fmt"
	"net"
	"time"

	"sync"

	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type mode int

const (
	remoteServer mode = iota
	relayClient
	socksClient
)

const udpBufSize = 64 * 1024

// Listen on laddr for UDP packets, encrypt and send to server to reach target.
func udpLocal(laddr, server, target string, shadow func(net.PacketConn) net.PacketConn) error {
	server = getClient(target)
	if server == "" {
		return nil // fmt.Errorf("UDP target address error: invalid target address: %q", target)
	}
	srvAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return fmt.Errorf("UDP server address error: %v", err)
	}

	tgt := socks.ParseAddr(target)
	if tgt == nil {
		return fmt.Errorf("UDP target address error: invalid target address: %q", target)
	}

	c, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return fmt.Errorf("UDP local listen error: %v", err)
	}
	defer c.Close()

	nm := newNATmap(config.UDPTimeout)
	buf := make([]byte, udpBufSize)
	copy(buf, tgt)

	logf("UDP tunnel %s <-> %s <-> %s", laddr, server, target)
	for {
		n, raddr, err := c.ReadFrom(buf[len(tgt):])
		if err != nil {
			logf("UDP local read error: %v", err)
			continue
		}

		pc := nm.Get(raddr.String())
		if pc == nil {
			pc, err = net.ListenPacket("udp", "")
			if err != nil {
				logf("UDP local listen error: %v", err)
				continue
			}

			pc = shadow(pc)
			nm.Add(raddr, c, pc, relayClient)
		}

		_, err = pc.WriteTo(buf[:len(tgt)+n], srvAddr)
		if err != nil {
			logf("UDP local write error: %v", err)
			continue
		}
	}
}

// Listen on laddr for Socks5 UDP packets, encrypt and send to server to reach target.
func udpSocksLocal(laddr, server string, shadow func(net.PacketConn) net.PacketConn) error {
	c, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return fmt.Errorf("UDP Socks local listen error: %v", err)
	}
	defer c.Close()

	nm := newNATmap(config.UDPTimeout)
	buf := make([]byte, udpBufSize)

	for {
		n, raddr, err := c.ReadFrom(buf)
		if err != nil {
			logf("UDP local read error: %v", err)
			continue
		}

		pc := nm.Get(raddr.String())
		if pc == nil {
			pc, err = net.ListenPacket("udp", "")
			if err != nil {
				logf("UDP local listen error: %v", err)
				continue
			}
			// logf("UDP socks tunnel %s <-> %s <-> %s", laddr, server, socks.Addr(buf[3:]))
			pc = shadow(pc)
			nm.Add(raddr, c, pc, socksClient)
		}

		dest := socks.Addr(buf[3:])
		server = getClient(dest.String())
		if server == "" {
			// logf("UDP target address error: invalid target address: %q", dest)
			continue
		}
		srvAddr, err := net.ResolveUDPAddr("udp", server)
		if err != nil {
			return fmt.Errorf("UDP server address error: %v", err)
		}

		_, err = pc.WriteTo(buf[3:n], srvAddr)
		if err != nil {
			logf("UDP local write error: %v", err)
			continue
		}
	}
}

// Listen on addr for encrypted packets and basically do UDP NAT.
func udpRemote(addr string, shadow func(net.PacketConn) net.PacketConn) error {
	c, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("UDP remote listen error: %v", err)
	}
	defer c.Close()
	c = shadow(c)

	nm := newNATmap(config.UDPTimeout)
	buf := make([]byte, udpBufSize)

	logf("listening UDP on %s", addr)
	for {
		n, raddr, err := c.ReadFrom(buf)
		if err != nil {
			logf("UDP remote read error: %v", err)
			continue
		}

		tgtAddr := socks.SplitAddr(buf[:n])
		if tgtAddr == nil {
			logf("failed to split target address from packet: %q", buf[:n])
			continue
		}

		tgtUDPAddr, err := net.ResolveUDPAddr("udp", tgtAddr.String())
		if err != nil {
			logf("failed to resolve target UDP address: %v", err)
			continue
		}

		payload := buf[len(tgtAddr):n]

		pc := nm.Get(raddr.String())
		if pc == nil {
			pc, err = net.ListenPacket("udp", "")
			if err != nil {
				logf("UDP remote listen error: %v", err)
				continue
			}

			nm.Add(raddr, c, pc, remoteServer)
		}

		_, err = pc.WriteTo(payload, tgtUDPAddr) // accept only UDPAddr despite the signature
		if err != nil {
			logf("UDP remote write error: %v", err)
			continue
		}
	}
}

// Packet NAT table
type natmap struct {
	sync.RWMutex
	m       map[string]net.PacketConn
	timeout time.Duration
}

func newNATmap(timeout time.Duration) *natmap {
	m := &natmap{}
	m.m = make(map[string]net.PacketConn)
	m.timeout = timeout
	return m
}

func (m *natmap) Get(key string) net.PacketConn {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *natmap) Set(key string, pc net.PacketConn) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = pc
}

func (m *natmap) Del(key string) net.PacketConn {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return nil
}

func (m *natmap) Add(peer net.Addr, dst, src net.PacketConn, role mode) {
	m.Set(peer.String(), src)

	go func() {
		timedCopy(dst, peer, src, m.timeout, role)
		if pc := m.Del(peer.String()); pc != nil {
			pc.Close()
		}
	}()
}

// copy from src to dst at target with read timeout
func timedCopy(dst net.PacketConn, target net.Addr, src net.PacketConn, timeout time.Duration, role mode) error {
	buf := make([]byte, udpBufSize)

	for {
		src.SetReadDeadline(time.Now().Add(timeout))
		n, raddr, err := src.ReadFrom(buf)
		if err != nil {
			return err
		}

		switch role {
		case remoteServer: // server -> client: add original packet source
			srcAddr := socks.ParseAddr(raddr.String())
			copy(buf[len(srcAddr):], buf[:n])
			copy(buf, srcAddr)
			_, err = dst.WriteTo(buf[:len(srcAddr)+n], target)
		case relayClient: // client -> user: strip original packet source
			srcAddr := socks.SplitAddr(buf[:n])
			_, err = dst.WriteTo(buf[len(srcAddr):n], target)
		case socksClient: // client -> socks5 program: just set RSV and FRAG = 0
			_, err = dst.WriteTo(append([]byte{0, 0, 0}, buf[:n]...), target)
		}

		if err != nil {
			return err
		}
	}
}
