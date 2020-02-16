package rest

import (
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func InitRouter(h *healthController, ur *registerController) *gin.Engine {
	r := gin.Default()
	MapRoutes(r, h, ur)
	return r
}

func MapRoutes(r *gin.Engine, h *healthController, ur *registerController) {
	// Health
	health := r.Group("/health")
	{
		health.GET("/ping", h.Ping)
	}

	user := r.Group("/users")
	{
		user.POST("/", ur.SignUp)
	}
}
