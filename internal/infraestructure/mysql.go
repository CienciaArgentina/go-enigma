package infraestructure

import (
	"errors"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	errNoConfig = errors.New("No hay archivo de configuración para la db")
)

func New(c *config.Configuration) *sqlx.DB {
	if c == nil {
		panic(errNoConfig)
	}
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Database.Username, c.Database.Password, c.Database.Hostname, c.Database.Port, c.Database.Database))
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
