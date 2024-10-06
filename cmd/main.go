package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"sheduler/internal/handlers"
	"sheduler/internal/storage"

)

func main() {
	//открываем файл для логирования

	//создаём connection к БД
	db, err := storage.ConnectionDB()
	if err != nil {
		log.Error().Err(err).Msg("Error connect to database")
		return
	}

	handlers := handlers.Handlers{db}

	router := gin.Default()
	router.POST("/api/task/add", handlers.AddTask)
	router.GET("/api/tasks/:id", handlers.GetTask)
	router.PUT("/api/tasks", handlers.PutTask)
	router.DELETE("/api/tasks/:id", handlers.DeleteTask)

	err = router.Run("localhost:7070")
	if err != nil {
		log.Error().Err(err).Msg("Error starting server")
		return
	}

}
