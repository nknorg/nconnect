package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nknorg/nconnect/config"
	tunnel "github.com/nknorg/nkn-tunnel"
)

func StartWeb(listenAddr string, tun *tunnel.Tunnel, conf *config.Config) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/rpc/admin", func(c *gin.Context) {
		req := &rpcReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := handleRequest(req, conf, tun)
		c.JSON(http.StatusOK, resp)
	})

	r.StaticFile("/", "web/dist/index.html")
	r.Static("/static", "web/dist/static")

	return r.Run(listenAddr)
}
