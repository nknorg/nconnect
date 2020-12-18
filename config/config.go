package config

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/nknorg/nconnect/util"
)

const (
	DefaultTunaMaxPrice    = "0.01"
	DefaultCipher          = "dummy"
	DefaultClientLocalAddr = "127.0.0.1:1080"
	RandomIdentifierChars  = "abcdefghijklmnopqrstuvwxyz0123456789"
	RandomIdentifierLength = 6
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	path string

	Identifier        string   `json:"identifier" long:"identifier" description:"NKN client identifier. A random one will be generated and saved to config.json if not provided"`
	Seed              string   `json:"seed" long:"seed" description:"NKN client secret seed. A random one will be generated and saved to config.json if not provided"`
	SeedRPCServerAddr []string `json:"seedRPCServerAddr,omitempty" long:"rpc" description:"Seed RPC server address"`

	LocalAddr  string `json:"localAddr,omitempty" short:"l" long:"local-addr" description:"(client only) Local socks proxy listen address (e.g. 127.0.0.1:1080)"`
	RemoteAddr string `json:"remoteAddr,omitempty" short:"r" long:"remote-addr" description:"(client only) Remote server NKN address"`

	Cipher   string `json:"cipher,omitempty" long:"cipher" description:"Socks proxy cipher. By default dummy (no cipher) will be used since NKN tunnel is already doing end to end encryption." choice:"dummy" choice:"chacha20-ietf-poly1305" choice:"aes-128-gcm" choice:"aes-256-gcm"`
	Password string `json:"password,omitempty" long:"password" description:"Socks proxy password"`

	Tuna                 bool     `json:"tuna,omitempty" short:"t" long:"tuna" description:"Enable tuna sessions"`
	TunaMaxPrice         string   `json:"tunaMaxPrice,omitempty" long:"tuna-max-price" description:"(server only) Tuna max price in unit of NKN/MB"`
	TunaCountry          []string `json:"tunaCountry,omitempty" long:"tuna-country" description:"(server only) Tuna service node allowed country code, e.g. US. All countries will be allowed if not provided"`
	TunaServiceName      string   `json:"tunaServiceName,omitempty" long:"tuna-service-name" description:"(server only) Tuna reverse service name"`
	TunaDownloadGeoDB    bool     `json:"tunaDownloadGeoDB,omitempty" long:"tuna-download-geo-db" description:"(server only) Download Tuna geo db to disk"`
	TunaGeoDBPath        string   `json:"tunaGeoDBPath,omitempty" long:"tuna-geo-db-path" description:"(server only) Path to store Tuna geo db"`
	TunaMeasureBandwidth bool     `json:"tunaMeasureBandwidth,omitempty" long:"tuna-measure-bandwidth" description:"(server only) Let Tuna measure bandwidth and connect to service node with highest bandwidth"`

	AdminHTTPAddr   string `json:"adminHttpAddr,omitempty" long:"admin-http" description:"(server only) Admin web GUI listen address (e.g. 127.0.0.1:8000)"`
	AdminIdentifier string `json:"adminIdentifier,omitempty" long:"admin-identifier" description:"(server only) Admin NKN client identifier prefix"`

	Tags []string `json:"tags,omitempty" long:"tags" description:"(server only) Tags that will be included in get info api"`

	Verbose bool `json:"verbose" short:"v" long:"verbose" description:"Verbose mode, show logs on dialing/accepting connections"`

	lock        sync.RWMutex
	AcceptAddrs []string `json:"acceptAddrs"`
	AdminAddrs  []string `json:"adminAddrs"`
}

func NewConfig() *Config {
	return &Config{
		AcceptAddrs: make([]string, 0),
		AdminAddrs:  make([]string, 0),
	}
}

func LoadOrNewConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c := NewConfig()
			c.path = path
			c.save()
			return c, nil
		}
		return nil, err
	}

	c := &Config{
		path: path,
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) GetAcceptAddrs() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.AcceptAddrs
}

func (c *Config) SetAcceptAddrs(acceptAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AcceptAddrs = acceptAddrs
	return c.save()
}

func (c *Config) AddAcceptAddrs(acceptAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AcceptAddrs = util.MergeStrings(c.AcceptAddrs, acceptAddrs)
	return c.save()
}

func (c *Config) RemoveAcceptAddrs(acceptAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AcceptAddrs = util.RemoveStrings(c.AcceptAddrs, acceptAddrs)
	return c.save()
}

func (c *Config) GetAdminAddrs() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.AdminAddrs
}

func (c *Config) SetAdminAddrs(adminAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AdminAddrs = adminAddrs
	return c.save()
}

func (c *Config) AddAdminAddrs(adminAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AdminAddrs = util.MergeStrings(c.AdminAddrs, adminAddrs)
	return c.save()
}

func (c *Config) RemoveAdminAddrs(adminAddrs []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.AdminAddrs = util.RemoveStrings(c.AdminAddrs, adminAddrs)
	return c.save()
}

func (c *Config) Save() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.save()
}

func (c *Config) save() error {
	if len(c.path) == 0 {
		return nil
	}

	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.path, b, 0666)
	if err != nil {
		return err
	}

	return nil
}

func RandomIdentifier() string {
	b := make([]byte, RandomIdentifierLength)
	for i := range b {
		b[i] = RandomIdentifierChars[rand.Intn(len(RandomIdentifierChars))]
	}
	return string(b)
}
