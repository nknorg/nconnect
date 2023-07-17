package network

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nknorg/nconnect/admin"
	"github.com/nknorg/nconnect/util"
)

const (
	success          = "success"
	defaultAdminAddr = "127.0.0.1:8000"
)

type addressData struct {
	Address string `json:"address"`
}

type addresses struct {
	Address         string   `json:"address"`
	AcceptAddresses []string `json:"acceptAddresses"`
}

type sendTokenData struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

func (m *Manager) StartWebServer() error {
	if m.opts.AdminHTTPAddr == "" {
		m.opts.AdminHTTPAddr = defaultAdminAddr
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New() // gin.Default()

	// This is for development, when start web page with "yarn dev" at ../web/src
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	r.POST("/rpc/network", func(c *gin.Context) {
		req := &admin.RpcReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := m.handleWebRequest(req)
		if m.opts.Verbose {
			log.Printf("Web request %v, response %+v\n", req.Method, resp)
		}

		c.JSON(http.StatusOK, resp)
	})

	r.StaticFile("/network", path.Join(m.opts.WebRootPath, "network.html"))
	r.StaticFile("/favicon.ico", path.Join(m.opts.WebRootPath, "favicon.ico"))
	r.StaticFile("/sw.js", path.Join(m.opts.WebRootPath, "sw.js"))
	r.Static("/static", path.Join(m.opts.WebRootPath, "static"))
	r.Static("/_nuxt", path.Join(m.opts.WebRootPath, "_nuxt"))
	r.Static("/img", path.Join(m.opts.WebRootPath, "img"))
	r.Static("/zh", path.Join(m.opts.WebRootPath, "zh"))
	r.Static("/zh-TW", path.Join(m.opts.WebRootPath, "zh-TW"))

	log.Println("Network manager web serve at ", "http://"+m.opts.AdminHTTPAddr+"/network")
	return r.Run(m.opts.AdminHTTPAddr)
}

func (m *Manager) handleWebRequest(req *admin.RpcReq) *admin.RpcResp {
	resp := &admin.RpcResp{}
	var err error

	switch req.Method {
	case "getNetworkConfig":
		resp.Result = m.GetNetworkConfig()

	case "setNetworkConfig":
		params := &networkData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		err = m.SetNetworkConfig(params)
		resp.Result = success

	case "authorizeMember":
		params := &addressData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		m.AuthorizeMemeber(params.Address)

		resp.Result = success

	case "removeMember":
		params := &addressData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		m.RemoveMember(params.Address)

		resp.Result = success

	case "deleteWaiting":
		params := &addressData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		m.RemoveMember(params.Address)

		resp.Result = success

	case "setAcceptAddress":
		params := &addresses{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}

		m.SetAcceptAddress(params.Address, params.AcceptAddresses)
		resp.Result = success

	case "sendToken":
		params := &sendTokenData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		if err = m.SendToken(params.Address, params.Amount); err != nil {
			break
		}
		resp.Result = success

	case "nknPing":
		fmt.Println("got network webservice nknPing")
		params := &addressData{}
		if err = util.JSONConvert(req.Params, params); err != nil {
			break
		}
		ms, err := m.NknPing(params.Address)
		if err != nil {
			break
		}
		resp.Result = fmt.Sprintf("%s, RTT time = %v ms", success, ms)

	default:
		resp.Error = "nConnect manager webservice got unknown method"
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp
}
