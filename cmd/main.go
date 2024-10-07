package main

import (
	"github.com/rs/zerolog/log"

	"sheduler/internal/server"
)

func main() {
	err := server.StartServer()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}
