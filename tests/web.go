package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func StartWebServer() error {
	http.HandleFunc("/httpEcho", httpEcho)
	fmt.Println("WEB server is serving at ", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
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

	b, _ := json.Marshal(user)
	w.Write(b)
}

func StartWebClient(httpServUrl string) error {
	fmt.Printf("http request to: %v\n", httpServUrl)

	proxyAddr := fmt.Sprintf("127.0.0.1:%v", port)
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

	for i := 0; i < numMsgs; i++ {
		user.Age++
		err := json.NewEncoder(b).Encode(user)
		if err != nil {
			fmt.Printf("StartWebClient.Encode err: %v\n", err)
			return err
		}
		req, err := http.NewRequest(http.MethodPost, httpServUrl, b)
		req.Header.Set("Content-type", "application/json")

		if err != nil {
			fmt.Printf("StartWebClient.http.NewRequest err: %v\n", err)
			return err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("StartWebClient.http.Do err: %v\n", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("StartWebClient.io.ReadAll err: %v\n", err)
			return err
		}

		respUser := &Person{}
		err = json.Unmarshal(body, respUser)
		if err != nil {
			fmt.Printf("StartWebClient.json.Unmarshal %v err: %v\n", string(body), err)
			return err
		}

		if respUser.Age != user.Age {
			return fmt.Errorf("StartWebClient got wrong response, sent %+v, recv %+v", user, respUser)
		} else {
			fmt.Printf("StartWebClient got echo: %+v\n", respUser)
		}
	}

	time.Sleep(time.Second)

	return nil
}

func StartTunWebClient(httpServUrl string) error {
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
			return err
		}
		req, err := http.NewRequest(http.MethodPost, httpServUrl, b)
		req.Header.Set("Content-type", "application/json")

		if err != nil {
			fmt.Printf("StartWebClient.http.NewRequest err: %v\n", err)
			return err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("StartWebClient.http.Do err: %v\n", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("StartWebClient.io.ReadAll err: %v\n", err)
			return err
		}

		respUser := &Person{}
		err = json.Unmarshal(body, respUser)
		if err != nil {
			fmt.Printf("StartWebClient.json.Unmarshal err: %v\n", err)
			return err
		}

		if respUser.Age != user.Age {
			return fmt.Errorf("StartWebClient got wrong response, sent %+v, recv %+v", user, respUser)
		}
	}

	return nil
}
