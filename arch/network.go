package arch

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/nknorg/nconnect/util"
)

const (
	mtu = 1500
)

func OpenTun(tunName, ip, gateway, mask, dns, socksAddr string) error {
	tunDevice, err := openTunDevice(tunName, ip, gateway, mask, []string{dns}, false)
	if err != nil {
		return fmt.Errorf("failed to open TUN device: %v", err)
	}

	core.RegisterOutputFn(tunDevice.Write)

	proxyAddr, err := net.ResolveTCPAddr("tcp", socksAddr)
	if err != nil {
		return fmt.Errorf("invalid proxy server address %v err: %v", socksAddr, err)
	}
	proxyHost := proxyAddr.IP.String()
	proxyPort := uint16(proxyAddr.Port)

	core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort))
	core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30*time.Second))

	lwipWriter := core.NewLWIPStack()

	go func() {
		_, err := io.CopyBuffer(lwipWriter, tunDevice, make([]byte, mtu))
		if err != nil {
			log.Fatalf("Failed to write data to network stack: %v", err)
		}
	}()

	return nil
}

func SetVPNRoutes(tunName, gateway string, cidrs []*net.IPNet) ([]*net.IPNet, error) {
	for _, dest := range cidrs {
		log.Printf("Adding route %s by %s", dest, gateway)
		out, err := addRouteCmd(dest, gateway, tunName)
		if len(out) > 0 {
			os.Stdout.Write(out)
		}
		if err != nil {
			os.Stdout.Write([]byte(util.ParseExecError(err)))
			os.Exit(1)
		}
	}

	return cidrs, nil
}

func RemoveVPNRoutes(tunName, gateway string, cidrs []*net.IPNet) error {
	for _, dest := range cidrs {
		log.Printf("Deleting route %s", dest)
		out, err := deleteRouteCmd(dest, gateway, tunName)
		if len(out) > 0 {
			os.Stdout.Write(out)
		}
		if err != nil {
			os.Stdout.Write([]byte(util.ParseExecError(err)))
		}
	}
	return nil
}
