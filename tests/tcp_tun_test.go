package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

// go test -v -run=TestTCPByTun
func TestTCPByTun(t *testing.T) {
	go StartTcpServer()

	tuna, udp, tun := true, true, false
	go func() {
		err := startNconnect("server.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect server err: %v\n", err)
			return
		}
	}()

	time.Sleep(15 * time.Second)

	tun = true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(10 * time.Second)

	go StartTCPClient()

	waitFor(ch, exited)
}

func StartTCPTunClient() error {
	conn, err := net.Dial("tcp", serverAddr)
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
