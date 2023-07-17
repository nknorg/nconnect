package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

func StartUDPServer(port string) error {
	a, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return err
	}
	udpServer, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}

	defer udpServer.Close()
	log.Printf("UDP server is listening at %v\n", port)

	b := make([]byte, 1024)
	for {
		n, addr, err := udpServer.ReadFromUDP(b)
		if err != nil {
			log.Printf("StartUdpServer.ReadFromUDP err: %v\n", err)
			return err
		}
		log.Printf("UDP Server got: %v from %v\n", string(b[:n]), addr.String())

		time.Sleep(100 * time.Millisecond)
		_, _, err = udpServer.WriteMsgUDP(b[:n], nil, addr)
		if err != nil {
			log.Printf("StartUdpServer.WriteMsgUDP err: %v\n", err)
			return err
		}
	}
}

func StartUDPClient(serverAddr string) error {
	proxyAddr := fmt.Sprintf("127.0.0.1:%v", port)
	s5c, err := socks5.NewClient(proxyAddr, "", "", 0, 60)
	if err != nil {
		return err
	}
	uc, err := s5c.Dial("udp", serverAddr)
	if err != nil {
		log.Println("StartUDPClient.s5c.Dial err: ", err)
		return err
	}
	defer uc.Close()

	user := &Person{Name: "udp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
		user.Age++
		send, _ := json.Marshal(user)
		if _, err := uc.Write(send); err != nil {
			log.Println("StartUDPClient.Write err ", err)
			return err
		}

		recv := make([]byte, 512)
		n, err := uc.Read(recv)
		if err != nil {
			log.Println("StartUDPClient.Read err ", err)
			return err
		}
		if !bytes.Equal(recv[:n], send) {
			return fmt.Errorf("StartUDPClient.recv %v is not as same as sent %v", string(recv[:n]), string(send))
		} else {
			log.Printf("StartUDPClient got echo: %v\n", string(recv[:n]))
		}
	}

	return nil
}

func StartUDPTunClient(serverAddr string) error {
	uc, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Println("StartUDPClient dial err: ", err)
		return err
	}
	defer uc.Close()

	user := &Person{Name: "udp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
		user.Age++
		send, _ := json.Marshal(user)
		if _, err := uc.Write(send); err != nil {
			log.Println("UDP client Write err ", err)
			return err
		}

		recv := make([]byte, 512)
		n, err := uc.Read(recv)
		if err != nil {
			log.Println("UDP client Read err ", err)
			return err
		}
		if !bytes.Equal(recv[:n], send) {
			return fmt.Errorf("UDP client recv %v is not as same as sent %v", string(recv[:n]), string(send))
		} else {
			log.Printf("UDP client got echo: %v\n", string(recv[:n]))
		}
	}

	return nil
}
