package main

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/CienciaArgentina/go-enigma/internal/storage/database"
	"github.com/CienciaArgentina/go-enigma/internal/storage/database/repositories"
)

func main() {

	cfg := config.New()
	db := database.New(cfg)

	h := rest.NewHealthController()

	regRepo := repositories.NewRegisterRepository(db)
	regSvc := register.NewService(regRepo, nil, cfg)
	regCtrl := rest.NewRegisterController(regSvc)

	logRepo := repositories.NewLoginRepository(db)
	logsvc := login.NewService(logRepo, nil, cfg)
	logCtrl := rest.NewLoginController(logsvc)

	if err := rest.InitRouter(h, regCtrl, logCtrl).Run(cfg.Server.Port); err != nil {
		panic(err)
	}
}
