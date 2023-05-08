package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// go test -v -run=TestHttpByTun
func TestHttpByTun(t *testing.T) {
	tunaNode, err := getTunaNode()
	if err != nil {
		fmt.Printf("startTunaNode err %v\n", err)
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
	tun = true
	go func() {
		err := startNconnect("client.json", tuna, udp, tun, nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(15 * time.Second)

	go StartTunWebClient()

	waitFor(ch, exited)
}

func StartTunWebClient() error {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	user := &Person{Name: "http_tun_boy", Age: 0}
	b := new(bytes.Buffer)

	for i := 0; i < 10; i++ {
		user.Age++
		err := json.NewEncoder(b).Encode(user)
		if err != nil {
			fmt.Printf("StartWebClient.Encode err: %v\n", err)
			break
		}
		req, err := http.NewRequest(http.MethodPost, httpServiceUrl, b)
		req.Header.Set("Content-type", "application/json")

		if err != nil {
			fmt.Printf("StartWebClient.http.NewRequest err: %v\n", err)
			break
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("StartWebClient.http.Do err: %v\n", err)
			break
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("StartWebClient.io.ReadAll err: %v\n", err)
			break
		}

		respUser := &Person{}
		err = json.Unmarshal(body, respUser)
		if err != nil {
			fmt.Printf("StartWebClient.json.Unmarshal err: %v\n", err)
			break
		}

		fmt.Printf("respUser %+v\n", respUser)
		if respUser.Age != user.Age {
			fmt.Printf("StartWebClient got wrong response, sent %+v, recv %+v\n", user, respUser)
			break
		}
	}

	time.Sleep(time.Second)
	ch <- webClientExited

	return nil
}
