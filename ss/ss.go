package ss

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type Config struct {
	Client     string
	Server     string
	Cipher     string
	Key        string
	Password   string
	Socks      string
	RedirTCP   string
	RedirTCP6  string
	TCPTun     string
	UDPTun     string
	UDPSocks   bool
	UDP        bool
	TCP        bool
	Plugin     string
	PluginOpts string
	Verbose    bool
	UDPTimeout time.Duration
	TCPCork    bool

	TargetToClient map[string]string // map target ip to local tunnel port
	DefaultClient  string            // the default client for the targets are not in Target2Client map
}

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
	TCPCork    bool
}

func Start(flags *Config) error {
	if flags.Client == "" && flags.Server == "" {
		return errors.New("at least one of client/server mode should be used")
	}

	config.Verbose = flags.Verbose
	config.UDPTimeout = flags.UDPTimeout
	config.TCPCork = flags.TCPCork

	routes.TargetToClient = flags.TargetToClient
	routes.DefaultClient = flags.DefaultClient

	var key []byte
	if flags.Key != "" {
		k, err := base64.URLEncoding.DecodeString(flags.Key)
		if err != nil {
			return err
		}
		key = k
	}

	errChan := make(chan error, 1)

	if flags.Client != "" { // client mode
		addr := flags.Client
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				return err
			}
		}

		udpAddr := addr

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			return err
		}

		if flags.Plugin != "" {
			addr, err = startPlugin(flags.Plugin, flags.PluginOpts, addr, false)
			if err != nil {
				return err
			}
		}

		if flags.UDPTun != "" {
			for _, tun := range strings.Split(flags.UDPTun, ",") {
				p := strings.Split(tun, "=")
				go func() {
					sendErr(udpLocal(p[0], udpAddr, p[1], ciph.PacketConn), errChan)
				}()
			}
		}

		if flags.TCPTun != "" {
			for _, tun := range strings.Split(flags.TCPTun, ",") {
				p := strings.Split(tun, "=")
				go func() {
					sendErr(tcpTun(p[0], addr, p[1], ciph.StreamConn), errChan)
				}()
			}
		}

		if flags.Socks != "" {
			socks.UDPEnabled = flags.UDPSocks
			go func() {
				sendErr(socksLocal(flags.Socks, addr, ciph.StreamConn), errChan)
			}()
			if flags.UDPSocks {
				go func() {
					sendErr(udpSocksLocal(flags.Socks, udpAddr, ciph.PacketConn), errChan)
				}()
			}
		}

		if flags.RedirTCP != "" {
			go func() {
				sendErr(redirLocal(flags.RedirTCP, addr, ciph.StreamConn), errChan)
			}()
		}

		if flags.RedirTCP6 != "" {
			go func() {
				sendErr(redir6Local(flags.RedirTCP6, addr, ciph.StreamConn), errChan)
			}()
		}
	}

	if flags.Server != "" { // server mode
		addr := flags.Server
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				return err
			}
		}

		udpAddr := addr

		if flags.Plugin != "" {
			addr, err = startPlugin(flags.Plugin, flags.PluginOpts, addr, true)
			if err != nil {
				return err
			}
		}

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			return err
		}

		if flags.UDP {
			go func() {
				sendErr(udpRemote(udpAddr, ciph.PacketConn), errChan)
			}()
		}
		if flags.TCP {
			go func() {
				sendErr(tcpRemote(addr, ciph.StreamConn), errChan)
			}()
		}
	}

	defer killPlugin()

	return <-errChan
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}

func sendErr(err error, errChan chan error) {
	select {
	case errChan <- err:
	default:
	}
}
