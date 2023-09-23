package main

import (
	"github.com/rs/zerolog/log"
	"hackathon/api"
	db "hackathon/db/sqlc"
	"hackathon/utils"
)

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
