package rest

import (
	"github.com/CienciaArgentina/go-enigma/conf"
	"github.com/gin-gonic/gin"
)

func InitRouter(c *conf.Configuration) *gin.Engine {
	router := gin.Default()
	gin.ForceConsoleColor()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	port := c.Server.Port
	if port == 0 {
		port = 8080
	}
	return router
}
