package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"sheduler/internal/server"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	Init()
	err := server.StartServer()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}
