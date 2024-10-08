package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"sheduler/internal/auth"
	"sheduler/internal/handlers"
	"sheduler/internal/storage"
)

func StartServer() error {
	db, err := storage.ConnectionDB()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка подключения к базе данных.")
		return err
	}

	handlers := handlers.Handlers{db}

	router := gin.Default()
	router.Use(auth.WithAuth)
	router.POST("/api/task/add", handlers.AddTask)
	router.GET("/api/tasks/:id", handlers.GetTask)
	router.PUT("/api/tasks", handlers.PutTask)
	router.DELETE("/api/tasks/:id", handlers.DeleteTask)

	err = router.Run(fmt.Sprintf("localhost:%s", os.Getenv("TODO_PORT")))
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при запуске сервера.")
		return err
	}
	return nil
}
