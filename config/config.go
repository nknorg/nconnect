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

	Tuna            bool     `json:"tuna,omitempty" long:"tuna" description:"enable tuna sessions"`
	TunaMaxPrice    string   `json:"tunaMaxPrice,omitempty" long:"tuna-max-price" description:"Tuna max price in unit of NKN/MB"`
	TunaCountry     []string `json:"tunaCountry,omitempty" long:"tuna-country" description:"Tuna service node allowed country code, e.g. US. All countries will be allowed if not provided"`
	TunaServiceName string   `json:"tunaServiceName,omitempty" long:"tuna-service-name" description:"Tuna reverse service name"`

	Password string `json:"password,omitempty" long:"password" description:"Socks proxy password"`

	AdminHTTPAddr   string `json:"adminHttpAddr,omitempty" long:"admin-http" description:"Admin web GUI listen address (e.g. 127.0.0.1:8000)"`
	AdminIdentifier string `json:"adminIdentifier,omitempty" long:"admin-identifier" description:"Admin NKN client identifier prefix"`

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
