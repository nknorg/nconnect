package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nknorg/nconnect/config"
	tunnel "github.com/nknorg/nkn-tunnel"
)

var (
	errAdminHTTPAPIDisabled = errors.New("Web API is disabled")
)

func StartWebServer(listenAddr string, tun *tunnel.Tunnel, persistConf, mergedConf *config.Config) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/rpc/admin", func(c *gin.Context) {
		req := &rpcReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if mergedConf.DisableAdminHTTPAPI {
			c.JSON(http.StatusOK, &rpcResp{Error: errAdminHTTPAPIDisabled.Error()})
			return
		}
		resp := handleRequest(req, persistConf, mergedConf, tun, rpcPermissionWeb)
		c.JSON(http.StatusOK, resp)
	})

	r.StaticFile("/", "web/dist/index.html")
	r.Static("/static", "web/dist/static")

	return r.Run(listenAddr)
}
