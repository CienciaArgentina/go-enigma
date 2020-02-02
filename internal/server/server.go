package server

import (
	"fmt"
	"github.com/CienciaArgentina/go-enigma/conf"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type Server struct {
	container *dig.Container
}

func New(c *dig.Container) *Server {
	return &Server{
		container: c,
	}
}

func (s Server) Start() error {
	var config *conf.Configuration
	if err := s.container.Invoke(func(c *conf.Configuration) { config = c }); err != nil {
		return err
	}
	return s.container.Invoke(func(r *gin.Engine) {
		r.Run(fmt.Sprintf(":%v", config.Server.Port))
	})
}
