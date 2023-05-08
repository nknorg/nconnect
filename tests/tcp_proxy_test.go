package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

// go test -v -run=TestTCPByProxy
func TestTCPByProxy(t *testing.T) {

	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("server.json", tuna, udp, tun, nil)
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

	go StartTCPClient()

	waitFor(ch, exited)
}

func StartTcpServer() error {
	tcpServer, err := net.Listen("tcp", tcpServerAddr)
	if err != nil {
		return err
	}
	fmt.Println("TCP Server is listening at ", tcpServerAddr)
	ch <- tcpServerIsReady

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			return err
		}
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if err != nil {
				fmt.Printf("StartTcpServer, Read err %v\n", err)
				return err
			}
			fmt.Printf("StartTcpServer, got: %v\n", string(b[:n]))
			_, err = conn.Write(b[:n])
			if err != nil {
				fmt.Printf("StartTcpServer, write err %v\n", err)
				return err
			}

			if strings.Contains(string(b[:n]), tcpClientExited) {
				ch <- tcpServerExited
				return nil
			}
		}
	}
}

func StartTCPClient() error {
	auth := proxy.Auth{User: "", Password: ""}
	dailer, err := proxy.SOCKS5("tcp", proxyAddr, &auth, &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		fmt.Printf("StartTCPClient, proxy.SOCKS5 err: %v\n", err)
		return err
	}

	conn, err := dailer.Dial("tcp", tcpServerAddr)
	if err != nil {
		fmt.Printf("StartTCPClient, dailer.Dial err: %v\n", err)
		return err
	}

	user := &Person{Name: "tcp_boy", Age: 0}
	for i := 0; i < rounds; i++ {
		user.Age++
		b1, _ := json.Marshal(user)
		_, err = conn.Write(b1)
		if err != nil {
			fmt.Printf("StartTCPClient, conn.Write err: %v\n", err)
			break
		}

		b2 := make([]byte, 1024)
		n, err := conn.Read(b2)
		if err != nil {
			fmt.Printf("StartTCPClient, conn.Read err: %v\n", err)
			break
		}
		respUser := &Person{}
		err = json.Unmarshal(b2[:n], respUser)
		if err != nil {
			fmt.Printf("StartTCPClient, json.Unmarshal err: %v\n", err)
			break
		}

		fmt.Printf("respUser %+v\n", respUser)
		if respUser.Age != user.Age {
			fmt.Printf("StartTCPClient, got wrong response, sent %+v, recv %+v\n", user, respUser)
			break
		}
	}
	conn.Write([]byte(tcpClientExited))
	time.Sleep(time.Second)
	ch <- tcpClientExited

	return nil
}
