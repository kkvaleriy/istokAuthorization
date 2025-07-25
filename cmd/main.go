package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkvaleriy/istokAuthorization/internal/app"
	"github.com/kkvaleriy/istokAuthorization/internal/config"
	"github.com/kkvaleriy/istokAuthorization/pkg/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	log := logger.New(cfg.Logger.Level)
	server := echo.New()
	connString := cfg.DataSource.PostgresConnString()
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error("can't parse connection string", "error", err.Error())
		panic(err)
	}

	dbConfig.MaxConns = int32(cfg.DataSource.MaxConnection)
	dbConfig.MinConns = int32(cfg.DataSource.MinConnection)
	dbConfig.MaxConnLifetime = cfg.DataSource.LifeTime()

	log.Info("attempting to connect to the database", "connectionString", connString)
	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("the attempt to connect to the database failed", "error", err.Error())
	}
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("database ping error", "error", err.Error())
	}
	log.Info("successful connection to the database")

	app := app.New(db, server, cfg, log)
	app.Run()
}
