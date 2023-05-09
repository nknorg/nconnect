package tests

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go StartTcpServer()
	go StartWebServer()
	go StartUdpServer()

	go StartNconnectServerWithTunaNode(true, true, false)
	time.Sleep(15 * time.Second)

	exitVal := m.Run()
	os.Exit(exitVal)
}
