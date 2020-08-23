package main

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
)

func main() {
	clog.SetLogLevel(clog.InfoLevel)

	if err := rest.InitRouter().Run(); err != nil {
		clog.Panic("Error starting app", "main", err, nil)
	}
}
