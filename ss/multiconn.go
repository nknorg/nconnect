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
	if server, ok := routes.TargetToClient[tgtIp[0]]; ok {
		return server
	}
	return routes.DefaultClient
}

func UpdateTargetToClient(targetToClient map[string]string) {
	routes.Lock()
	defer routes.Unlock()
	routes.TargetToClient = targetToClient
}
