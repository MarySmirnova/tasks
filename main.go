package main

import (
	"github.com/MarySmirnova/tasks/internal/config"
	"github.com/MarySmirnova/tasks/pkg/storage/postgres"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var cfg config.Application

func init() {
	_ = godotenv.Load(".env")
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := postgres.NewStorage(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.GetPGPool().Close()

	//Здесь должно запускаться наше приложение.

}
