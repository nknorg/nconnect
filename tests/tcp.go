package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func StartTCPServer(port string) error {
	tcpServer, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	log.Printf("TCP server is listening at %v\n", port)

	for {
		c, err := tcpServer.Accept()
		if err != nil {
			return err
		}

		log.Printf("TCP server get a connection from %v\n", c.RemoteAddr())

		go func(conn net.Conn) {
			defer conn.Close()
			b := make([]byte, 1024)
			for {
				n, err := conn.Read(b)
				if err != nil {
					if strings.Contains(err.Error(), "closed") {
						log.Printf("client connection closed\n")
					} else {
						log.Printf("StartTcpServer, Read err %v\n", err)
					}
					break
				}

				log.Printf("TCP Server got: %v\n", string(b[:n]))
				_, err = conn.Write(b[:n])
				if err != nil {
					log.Printf("StartTcpServer, write err %v\n", err)
					break
				}
			}
		}(c)
	}
}

func StartTCPClient(serverAddr string) error {
	auth := proxy.Auth{User: "", Password: ""}
	proxyAddr := fmt.Sprintf("127.0.0.1:%v", port)
	dailer, err := proxy.SOCKS5("tcp", proxyAddr, &auth, &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		log.Printf("StartTCPClient, proxy.SOCKS5 err: %v\n", err)
		return err
	}

	conn, err := dailer.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("StartTCPClient, dailer.Dial err: %v\n", err)
		return err
	}

	defer conn.Close()
	log.Printf("StartTCPClient, dail to %v success\n", serverAddr)

	user := &Person{Name: "tcp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
		user.Age++
		b1, _ := json.Marshal(user)
		_, err = conn.Write(b1)
		if err != nil {
			log.Printf("StartTCPClient, conn.Write err: %v\n", err)
			return err
		}
		log.Printf("StartTCPClient, conn.Write %+v\n", user)

		b2 := make([]byte, 1024)
		n, err := conn.Read(b2)
		if err != nil {
			log.Printf("StartTCPClient, conn.Read err: %v\n", err)
			return err
		}
		respUser := &Person{}
		err = json.Unmarshal(b2[:n], respUser)
		if err != nil {
			log.Printf("StartTCPClient, json.Unmarshal err: %v\n", err)
			return err
		}

		if respUser.Age != user.Age {
			return fmt.Errorf("StartTCPClient, got wrong response, sent %+v, recv %+v", user, respUser)
		}
		log.Printf("Got echo %+v\n", respUser)
	}

	return nil
}

func StartTCPTunClient(serverAddr string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("StartTCPClient, dailer.Dial err: %v\n", err)
		return err
	}

	user := &Person{Name: "tcp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
		user.Age++
		b1, _ := json.Marshal(user)
		_, err = conn.Write(b1)
		if err != nil {
			log.Printf("StartTCPClient, conn.Write err: %v\n", err)
			return err
		}

		b2 := make([]byte, 1024)
		n, err := conn.Read(b2)
		if err != nil {
			log.Printf("StartTCPClient, conn.Read err: %v\n", err)
			return err
		}
		respUser := &Person{}
		err = json.Unmarshal(b2[:n], respUser)
		if err != nil {
			log.Printf("StartTCPClient, json.Unmarshal err: %v\n", err)
			return err
		}
		log.Printf("TCP client got echo: %+v\n", respUser)

		if respUser.Age != user.Age {
			return fmt.Errorf("StartTCPClient, got wrong response, sent %+v, recv %+v", user, respUser)
		}
	}

	return nil
}
