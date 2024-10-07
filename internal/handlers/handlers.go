package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"sheduler/internal/storage"
	"sheduler/models"
)

type Handlers struct {
	TaskStorage storage.StorageInterface
}

func (h *Handlers) AddTask(c *gin.Context) {
	var userTask models.Task
	err := c.BindJSON(&userTask)
	if err != nil {
		log.Error().Err(err).Msg("Введены некорректные данные в теле запроса.")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ok := userTask.PostCheckingFields()
	if !ok {
		log.Error().Err(err).Msg("Неверно указанны данные в запросе.")
		c.JSON(http.StatusBadRequest, "Неверно указанны данные в запросе.")
		return
	}
	id, err := h.TaskStorage.AppendTask(userTask)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при добавлении задачи в базу данных")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Info().Msg("Запись успешно добавлена.")
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Запись с id = %d успешно добавлена.", id))
}

func (h *Handlers) GetTask(c *gin.Context) {
	var returnedTask models.Task
	id := c.Param("id")
	returnedTask, err := h.TaskStorage.FindTask(id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при поиске задачи в базе данных.")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Info().Msg("Запись успешно получена.")
	c.IndentedJSON(http.StatusOK, returnedTask)
}

func (h *Handlers) PutTask(c *gin.Context) {
	var userTask models.Task
	err := c.BindJSON(&userTask)
	if err != nil {
		log.Error().Err(err).Msg("Введены некорректные данные в теле запроса.")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ok := userTask.PutCheckingFields()
	if !ok {
		log.Error().Err(err).Msg("Неверно указанны данные в запросе.")
		c.JSON(http.StatusBadRequest, "Неверно указанны данные в запросе.")
		return
	}
	id, err := h.TaskStorage.ChangeTask(userTask)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при обновлении задачи.")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Info().Msg("Запись успешно изменена.")
	c.JSON(http.StatusOK, fmt.Sprintf("Запись c id = %d успешно изменена.", id))
}

func (h *Handlers) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	idTask, err := h.TaskStorage.RemoveTask(id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удалении задачи.")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if idTask == 0 {
		log.Error().Err(err).Msg("Задача отсутствует.")
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Задача отсутствует."))
		return
	}
	log.Info().Msg("Запись успешно удалена.")
	c.JSON(http.StatusOK, fmt.Sprintf("Запись успешно c id = %d удалена.", idTask))
}
