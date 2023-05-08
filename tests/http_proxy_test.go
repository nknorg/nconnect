package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// go test -v -run=TestHttpByProxy
func TestHttpByProxy(t *testing.T) {
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

	go StartWebClient()

	waitFor(ch, exited)
}

func StartWebServer() error {
	http.HandleFunc("/httpEcho", httpEcho)
	fmt.Println("WEB server is serving at ", httpServerAddr)
	if err := http.ListenAndServe(httpServerAddr, nil); err != nil {
		ch <- webServerExited
		log.Fatal(err)
	}

	return nil
}

func httpEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &Person{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("web server got user:", user)
	b, _ := json.Marshal(user)
	w.Write(b)
}

func StartWebClient() error {
	socksProxy := "socks5://" + proxyAddr
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(socksProxy)
	}

	httpTransport := &http.Transport{
		Proxy: proxy,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}

	user := &Person{Name: "http_boy", Age: 0}
	b := new(bytes.Buffer)

	for i := 0; i < rounds; i++ {
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
