package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nknorg/nconnect"
	"github.com/nknorg/nconnect/config"
	"github.com/txthinking/brook"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

const (
	SocksProxy = "socks5://127.0.0.1:1080"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello")
}

func runServer() {
	var serverConfig config.NConfig
	b, err := os.ReadFile("server.json")
	if err != nil {
		log.Fatalf("read config file err: %v", err)
		return
	}
	err = json.Unmarshal(b, &serverConfig)
	if err != nil {
		log.Fatalf("parse config err: %v", err)
		return
	}
	nconnect.Run(&serverConfig)
}

func runClient() {
	var clientConfig config.NConfig
	b, err := os.ReadFile("client.json")
	if err != nil {
		log.Fatalf("read config file err: %v", err)
		return
	}
	err = json.Unmarshal(b, &clientConfig)
	if err != nil {
		log.Fatalf("parse config err: %v", err)
		return
	}
	nconnect.Run(&clientConfig)
}

func TestMain(m *testing.M) {
	go runServer()
	time.Sleep(100 * time.Second)
	go runClient()
	time.Sleep(20 * time.Second)
	os.Exit(m.Run())
}

func TestTCPSocks5Proxy(t *testing.T) {
	http.HandleFunc("/hello", hello)
	go http.ListenAndServe(":22222", nil)
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(SocksProxy)
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

func TestUDPSocks5Proxy(t *testing.T) {
	err := brook.Socks5Test("127.0.0.1:1080", "", "", "http3.ooo", "137.184.237.95", "8.8.8.8:53")
	if err != nil {
		t.Fatal(err)
	}
}
