package tests

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nknorg/tuna/types"
)

var remoteTuna = flag.Bool("remoteTuna", false, "use remote tuna nodes")
var tun = flag.Bool("tun", false, "use tun device")

func TestMain(m *testing.M) {
	flag.Parse()
	if *remoteTuna {
		fmt.Println("We are using remote tuna node")
	} else {
		fmt.Println("Using local tuna node. If want to use remote tuna nodes, please run: go test -v -remoteTuna .")
	}

	go func() {
		err := StartTcpServer()
		if err != nil {
			log.Fatalf("StartTcpServer err %v", err)
			return
		}
	}()
	go func() {
		err := StartWebServer()
		if err != nil {
			log.Fatalf("StartWebServer err %v", err)
			return
		}
	}()
	go func() {
		err := StartUdpServer()
		if err != nil {
			log.Fatalf("StartUdpServer err %v", err)
			return
		}
	}()

	var tunaNode *types.Node
	var err error
	if !(*remoteTuna) {
		tunaNode, err = getTunaNode()
		if err != nil {
			log.Fatalf("getTunaNode err %v", err)
			return
		}
	}

	err = startNconnect("server.json", true, true, false, tunaNode)
	if err != nil {
		log.Fatalf("start nconnect server err: %v", err)
		return
	}

	time.Sleep(10 * time.Second)

	exitVal := m.Run()
	os.Exit(exitVal)
}
