package rest

import (
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func InitRouter(h *healthController, ur *registerController, l *loginController) *gin.Engine {
	r := gin.Default()
	MapRoutes(r, h, ur, l)
	return r
}

func MapRoutes(r *gin.Engine, h *healthController, ur *registerController, l *loginController) {
	// Health
	health := r.Group("/")
	{
		health.GET("/ping", h.Ping)
	}

	user := r.Group("/users")
	{
		user.POST("/", ur.SignUp)
		user.POST("/login", l.Login)
	}
}
