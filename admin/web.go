package admin

import (
	"errors"
	"net/http"
	"path"

	"github.com/gin-contrib/gzip"
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

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.POST("/rpc/admin", func(c *gin.Context) {
		req := &RpcReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if mergedConf.DisableAdminHTTPAPI {
			c.JSON(http.StatusOK, &RpcResp{Error: errAdminHTTPAPIDisabled.Error()})
			return
		}
		resp := handleRequest(req, persistConf, mergedConf, tun, rpcPermissionWeb)
		c.JSON(http.StatusOK, resp)
	})

	r.StaticFile("/", path.Join(mergedConf.WebRootPath, "index.html"))
	r.StaticFile("/favicon.ico", path.Join(mergedConf.WebRootPath, "favicon.ico"))
	r.StaticFile("/sw.js", path.Join(mergedConf.WebRootPath, "sw.js"))
	r.Static("/static", path.Join(mergedConf.WebRootPath, "static"))
	r.Static("/_nuxt", path.Join(mergedConf.WebRootPath, "_nuxt"))
	r.Static("/img", path.Join(mergedConf.WebRootPath, "img"))
	r.Static("/zh", path.Join(mergedConf.WebRootPath, "zh"))
	r.Static("/zh-TW", path.Join(mergedConf.WebRootPath, "zh-TW"))

	return r.Run(listenAddr)
}
