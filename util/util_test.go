package util

import (
	"log"
	"testing"

	ts "github.com/nknorg/nkn-tuna-session"
)

// go test -v -run=TestGetFreePort
func TestGetFreePort(t *testing.T) {
	port, err := ts.GetFreePort(0)
	if err != nil {
		log.Println(err)
	}
	log.Println(port)

	port, err = ts.GetFreePort(1080)
	if err != nil {
		log.Println(err)
	}
	log.Println(port)
}
