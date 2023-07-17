package admin

import (
	"encoding/json"
	"log"

	"github.com/nknorg/nconnect/config"
	"github.com/nknorg/nconnect/util"
	"github.com/nknorg/nkn-sdk-go"
	tunnel "github.com/nknorg/nkn-tunnel"
)

func StartNKNServer(account *nkn.Account, identifier string, clientConfig *nkn.ClientConfig, tun *tunnel.Tunnel, persistConf, mergedConf *config.Config) error {
	m, err := nkn.NewMultiClient(account, identifier, 4, false, clientConfig)
	if err != nil {
		return err
	}

	<-m.OnConnect.C

	serverAdminAddr = m.Address()

	for {
		msg := <-m.OnMessage.C

		req := &RpcReq{}
		err := json.Unmarshal(msg.Data, req)
		if err != nil {
			log.Println("Unmarshal client request error:", err)
			continue
		}

		isAcceptAddr := util.MatchRegex(persistConf.GetAcceptAddrs(), msg.Src)
		isAdminAddr := util.MatchRegex(persistConf.GetAdminAddrs(), msg.Src)

		if !isAdminAddr && tokenStore.IsValid(req.Token) {
			isAdminAddr = true
		}

		if !isAcceptAddr && !isAdminAddr {
			log.Println("Ignore authorized message from", msg.Src)
			continue
		}

		var perm permission
		if isAcceptAddr {
			perm |= rpcPermissionAcceptClient
		}
		if isAdminAddr {
			perm |= rpcPermissionAdminClient
		}

		resp := handleRequest(req, persistConf, mergedConf, tun, perm)

		b, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			continue
		}

		err = msg.Reply(string(b))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
