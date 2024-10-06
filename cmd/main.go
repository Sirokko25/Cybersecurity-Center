package main

import (
	"github.com/gin-gonic/gin"

	"sheduler/internal/handlers"
	"sheduler/internal/storage"

)

func main() {
	//открываем файл для логирования

	//создаём connection к БД
	db, err := storage.ConnectionDB()
	if err != nil {
		//log
	}

	handlers := handlers.Handlers{db}

	router := gin.Default()
	router.POST("/api/task/add", handlers.AddTask)
	router.GET("/api/tasks/:id", handlers.GetTask)
	router.PUT("/api/tasks", handlers.PutTask)
	router.DELETE("/api/tasks", handlers.DeleteTask)

	err = router.Run("localhost:7070")
	if err != nil {
		//log
		return
	}

}
