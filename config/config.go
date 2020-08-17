package config

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/nknorg/nkn-socks/util"
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

	Identifier        string   `json:"identifier"`
	Password          string   `json:"password"`
	Seed              string   `json:"seed"`
	SeedRPCServerAddr []string `json:"seedRPCServerAddr,omitempty"`
	TunaMaxPrice      string   `json:"tunaMaxPrice,omitempty"`

	lock        sync.RWMutex
	AcceptAddrs []string `json:"acceptAddrs"`
	AdminAddrs  []string `json:"adminAddrs"`
}

func NewConfig() *Config {
	return &Config{
		Identifier:  randomIdentifier(6),
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

func (c *Config) Save() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.save()
}

func randomIdentifier(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = RandomIdentifierChars[rand.Intn(len(RandomIdentifierChars))]
	}
	return string(b)
}
