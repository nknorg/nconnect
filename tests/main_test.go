package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	go StartTcpServer()
	go StartWebServer()
	go StartUdpServer()

	exitVal := m.Run()
	os.Exit(exitVal)
}
