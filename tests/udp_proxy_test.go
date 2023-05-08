package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/txthinking/socks5"
)

// go test -v -run=TestUDPByProxy
func TestUDPByProxy(t *testing.T) {
	tunaNode, err := getTunaNode()
	if err != nil {
		fmt.Printf("getTunaNode err %v\n", err)
		return
	}

	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("server.json", tuna, udp, tun, tunaNode)
		if err != nil {
			fmt.Printf("start nconnect server err: %v\n", err)
			return
		}
	}()

	time.Sleep(15 * time.Second)

	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()

	time.Sleep(15 * time.Second)

	go StartUDPClient()

	waitFor(ch, exited)
}

func StartUdpServer() error {
	a, err := net.ResolveUDPAddr("udp", udpServerAddr)
	if err != nil {
		return err
	}
	udpServer, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}

	fmt.Printf("UDP server is listening at %v\n", udpServerAddr)

	b := make([]byte, 1024)
	for {
		n, addr, err := udpServer.ReadFromUDP(b)
		if err != nil {
			fmt.Printf("StartUdpServer.ReadFromUDP err: %v\n", err)
			break
		}
		fmt.Printf("UDP Server got: %v\n", string(b[:n]))

		time.Sleep(100 * time.Millisecond)
		_, _, err = udpServer.WriteMsgUDP(b[:n], nil, addr)
		if err != nil {
			fmt.Printf("StartUdpServer.WriteMsgUDP err: %v\n", err)
			break
		}

		if strings.Contains(string(b[:n]), udpClientExited) {
			break
		}
	}

	ch <- udpServerExited
	return nil
}

func StartUDPClient() error {
	s5c, err := socks5.NewClient(proxyAddr, "", "", 0, 60)
	if err != nil {
		ch <- udpClientExited
		return err
	}
	uc, err := s5c.Dial("udp", udpServerAddr)
	if err != nil {
		fmt.Println("StartUDPClient.s5c.Dial err: ", err)
		ch <- udpClientExited
		return err
	}
	defer uc.Close()

	user := &Person{Name: "udp_boy", Age: 0}
	for i := 0; i < rounds; i++ {
		user.Age++
		send, _ := json.Marshal(user)
		if _, err := uc.Write(send); err != nil {
			fmt.Println("StartUDPClient.Write err ", err)
			break
		}

		recv := make([]byte, 512)
		n, err := uc.Read(recv)
		if err != nil {
			fmt.Println("StartUDPClient.Read err ", err)
			break
		}
		if !bytes.Equal(recv[:n], send) {
			fmt.Printf("StartUDPClient.recv %v is not as same as sent %v\n", string(recv[:n]), string(send))
			break
		}
	}

	uc.Write([]byte(udpClientExited))
	time.Sleep(time.Second)

	ch <- udpClientExited

	return nil
}
