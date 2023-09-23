package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"hackathon/api"
	db "hackathon/db/sqlc"
	"hackathon/utils"
)

func main() {
	ctx := context.Background()
	configuration, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("could not find configuration file")
	}

	conn, err := pgxpool.ParseConfig(configuration.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing database connection string failed")
	}

	pg, err := pgxpool.NewWithConfig(ctx, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to database failed")
	}
	defer pg.Close()

	store := db.NewStore(pg)

	runGinServer(configuration, store)
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("can not start server")
	}
}
