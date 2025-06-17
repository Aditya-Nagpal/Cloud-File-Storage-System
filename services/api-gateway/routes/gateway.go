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

	// USER-SERVICE
	userServiceUrl, err := url.Parse((config.AppConfig.UserServiceUrl))
	if err != nil {
		panic(err)
	}

	// Reverse proxy handlers
	authServiceProxy := httputil.NewSingleHostReverseProxy(authServiceUrl)
	fileServiceProxy := httputil.NewSingleHostReverseProxy(fileServiceUrl)
	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceUrl)

	// PUBLIC AUTH ROUTES
	r.Any("/auth/*proxyPath", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		authServiceProxy.ServeHTTP(c.Writer, c.Request)
	})

	// PROTECTED FILE ROUTES
	fileGroup := r.Group("/file")
	fileGroup.Use(middleware.JWTMiddleware(config.AppConfig.JwtSecret))
	{
		fileGroup.Any("/*proxyPath", func(c *gin.Context) {
			c.Request.URL.Path = c.Param("proxyPath")
			fileServiceProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// PROTECTED USER ROUTES
	userGroup := r.Group("/user")
	userGroup.Use(middleware.JWTMiddleware(config.AppConfig.JwtSecret))
	{
		userGroup.Any("/*proxyPath", func(c *gin.Context) {
			c.Request.URL.Path = c.Param("proxyPath")
			userServiceProxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}
