package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	)

type Server struct {
	router *gin.Engine
	container *dig.Container
}

func New