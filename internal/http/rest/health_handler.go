package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct{}

func NewHealthController() *healthController {
	return &healthController{}
}

func (h *healthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
