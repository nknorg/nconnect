package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

func StartTcpServer() error {
	tcpServer, err := net.Listen("tcp", tcpPort)
	if err != nil {
		return err
	}
	fmt.Println("TCP Server is listening at ", tcpPort)

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
				break
			}
			fmt.Printf("TCP Server got: %v\n", string(b[:n]))
			_, err = conn.Write(b[:n])
			if err != nil {
				fmt.Printf("StartTcpServer, write err %v\n", err)
				break
			}
		}
	}
}

func StartTCPClient(serverAddr string) error {
	auth := proxy.Auth{User: "", Password: ""}
	dailer, err := proxy.SOCKS5("tcp", proxyAddr, &auth, &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		fmt.Printf("StartTCPClient, proxy.SOCKS5 err: %v\n", err)
		return err
	}

	conn, err := dailer.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("StartTCPClient, dailer.Dial err: %v\n", err)
		return err
	}
	fmt.Printf("StartTCPClient, dail to %v  success\n", serverAddr)

	user := &Person{Name: "tcp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
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

		if respUser.Age != user.Age {
			fmt.Printf("StartTCPClient, got wrong response, sent %+v, recv %+v\n", user, respUser)
			break
		}
		fmt.Printf("Got echo %+v\n", respUser)
	}

	return nil
}

func StartTCPTunClient(serverAddr string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("StartTCPClient, dailer.Dial err: %v\n", err)
		return err
	}

	user := &Person{Name: "tcp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
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

		if respUser.Age != user.Age {
			fmt.Printf("StartTCPClient, got wrong response, sent %+v, recv %+v\n", user, respUser)
			break
		}
	}

	return nil
}
