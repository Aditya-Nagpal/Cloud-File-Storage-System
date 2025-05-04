package routes

import (
	// "net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/api-gateway/config"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authServiceUrl, err := url.Parse(config.AppConfig.AuthServiceUrl)
	if err != nil {
		panic(err)
	}

	// Reverse proxy handler
	authServiceProxy := httputil.NewSingleHostReverseProxy(authServiceUrl)

	r.Any("/auth/*proxyPath", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		authServiceProxy.ServeHTTP(c.Writer, c.Request)
	})
}
