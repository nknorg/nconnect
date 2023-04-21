package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/txthinking/brook"
)

func TestMain(m *testing.M) {
	n, err := startTunaNode()
	if err != nil {
		fmt.Printf("start tuna node err: %v\n", err)
		return
	}

	go func() {
		err = startNconnect("server.json", n)
		if err != nil {
			fmt.Printf("start nconnect server err: %v\n", err)
			return
		}
	}()
	time.Sleep(20 * time.Second)

	go func() {
		err = startNconnect("client.json", nil)
		if err != nil {
			fmt.Printf("start nconnect client err: %v\n", err)
			return
		}
	}()
	time.Sleep(20 * time.Second)

	os.Exit(m.Run())
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello")
}

// go test -v -run=TestTCPSocks5Proxy
func TestTCPSocks5Proxy(t *testing.T) {
	http.HandleFunc("/hello", hello)
	go http.ListenAndServe(":22222", nil)
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

	req, err := http.NewRequest("GET", "http://127.0.0.1:22222/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(body, []byte("hello")) {
		t.Fatal("bytes not equal")
	}
}

// go test -v -run=TestUDPSocks5Proxy
func TestUDPSocks5Proxy(t *testing.T) {
	for i := 1; i <= 5; i++ {
		err := brook.Socks5Test("127.0.0.1:1080", "", "", "http3.ooo", "137.184.237.95", "8.8.8.8:53")
		if err != nil {
			fmt.Printf("TestUDPSocks5Proxy try %v err: %v\n", i, err)
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			break
		}
	}
}
