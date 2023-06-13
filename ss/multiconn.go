package ss

import (
	"strings"
	"sync"
)

var routes struct {
	sync.RWMutex
	TargetToClient map[string]string // map target ip to local tunnel port
	DefaultClient  string            // the default client for the targets are not in TargetToClient map
}

func getClient(target string) string {
	tgtIp := strings.Split(target, ":")

	routes.RLock()
	defer routes.RUnlock()
	server, ok := routes.TargetToClient[tgtIp[0]]

	if ok {
		return server
	}
	return routes.DefaultClient
}
