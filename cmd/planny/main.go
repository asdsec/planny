package main

import (
	"com.github/asdsec/planny/configs"
	"com.github/asdsec/planny/internal/api"
	"com.github/asdsec/planny/internal/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func main() {
	conf, err := configs.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	if conf.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := gorm.Open(mysql.Open(conf.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = conn.AutoMigrate(
		&db.UserEntity{},
		&db.SessionEntity{},
		&db.PlanEntity{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot migrate db")
	}
	store := db.NewStore(conn)

	serv, err := api.NewServer(conf, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = serv.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
