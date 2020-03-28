package main

import (
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
	"github.com/CienciaArgentina/go-enigma/internal/listing"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/CienciaArgentina/go-enigma/internal/recovery"
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/CienciaArgentina/go-enigma/internal/storage/database"
	"github.com/CienciaArgentina/go-enigma/internal/storage/database/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {

	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Inicializando configuración")
	start := time.Now()
	cfg := config.New()
	elapsed := time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Configuración cargada")

	logrus.Info("Inicializando db")
	start = time.Now()
	db := database.New(cfg)
	elapsed = time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Db cargada")

	h := rest.NewHealthController()

	regRepo := repositories.NewRegisterRepository(db)
	regSvc := register.NewService(regRepo, nil, cfg)
	regCtrl := rest.NewRegisterController(regSvc)

	logRepo := repositories.NewLoginRepository(db)
	logsvc := login.NewService(logRepo, nil, cfg)
	logCtrl := rest.NewLoginController(logsvc)

	recRepo := repositories.NewRecoveryRepository(db)
	recSvc := recovery.NewService(recRepo, cfg)
	recCtlr := rest.NewRecoveryController(recSvc)

	lisRepo := repositories.NewListingRepository(db)
	lisSvc := listing.NewService(lisRepo)
	lisCtlr := rest.NewListingController(lisSvc)

	if err := rest.InitRouter(h, regCtrl, logCtrl, recCtlr, lisCtlr).Run(cfg.Server.Port); err != nil {
		panic(err)
	}
}
