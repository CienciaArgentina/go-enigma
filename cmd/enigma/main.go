package main

import (
	"fmt"
	"github.com/CienciaArgentina/go-enigma/conf"
	"github.com/CienciaArgentina/go-enigma/internal/di"
	"github.com/CienciaArgentina/go-enigma/internal/server"
)

func main() {
	conf.ConsolePrintMessageByCienciaArgentina("Starting Enigma, please stand by...\n")
	d := di.BuildContainer()

	srv := server.New(d)
	if err := srv.Start(); err != nil {
		fmt.Print(err)
	}
}
