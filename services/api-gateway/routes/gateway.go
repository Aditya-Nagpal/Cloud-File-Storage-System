package routes

import (
	// "net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/api-gateway/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// AUTH-SERVICE
	authServiceUrl, err := url.Parse(config.AppConfig.AuthServiceUrl)
	if err != nil {
		panic(err)
	}

	// FILE-SERVICE
	fileServiceUrl, err := url.Parse(config.AppConfig.FileServiceUrl)
	if err != nil {
		panic(err)
	}

	// Reverse proxy handlers
	authServiceProxy := httputil.NewSingleHostReverseProxy(authServiceUrl)
	fileServiceProxy := httputil.NewSingleHostReverseProxy(fileServiceUrl)

	r.Any("/auth/*proxyPath", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		authServiceProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Protected file routes
	// Assumption is auth-service doesn't require middleware check
	protected := r.Group("/file")
	protected.Use(middleware.JWTMiddleware(config.AppConfig.JwtSecret))
	{
		protected.Any("/*proxyPath", func(c *gin.Context) {
			c.Request.URL.Path = c.Param("proxyPath")
			fileServiceProxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}
