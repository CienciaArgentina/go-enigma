package di

import (
	"github.com/CienciaArgentina/go-enigma/conf"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

var container = dig.New()

func BuildContainer() *dig.Container {
	container.Provide(conf.New)
	container.Provide(func(c *conf.Configuration) *gin.Engine {
		return rest.InitRouter(c)
	})
	return container
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}
