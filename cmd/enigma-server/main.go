package main

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.InfoLevel)

	cfg := config.New()

	if err := rest.InitRouter(cfg).Run(cfg.Server.Port); err != nil {
		panic(err)
	}
}
