package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"sheduler/internal/storage"
	"sheduler/models"

)

type Handlers struct {
	TaskStorage storage.DB
}

func (h *Handlers) AddTask(c *gin.Context) {
	var userTask models.Task
	err := c.BindJSON(&userTask)
	if err != nil {
		log.Error().Err(err).Msg("Введены некорректные данные в теле запроса.")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	code, err := userTask.PostCheckingFields()
	if err != nil {
		log.Error().Err(err).Msg("Недостаточно информации в теле запроса.")
		c.JSON(code, err)
		return
	}
	code, err = h.TaskStorage.AppendTask(userTask)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при добавлении задачи в базу данных")
		c.JSON(code, err)
		return
	}
	log.Info().Msg("Запись успешно добавлена.")
	c.IndentedJSON(http.StatusOK, "Запись успешно добавлена.")
}

func (h *Handlers) GetTask(c *gin.Context) {
	var returnedTask models.Task
	id := c.Param("id")
	returnedTask, code, err := h.TaskStorage.FindTask(id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при поиске задачи в базе данных")
		c.JSON(code, err)
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
	code, err := userTask.PutCheckingFields()
	if err != nil {
		log.Error().Err(err).Msg("Недостаточно информации в теле запроса.")
		c.JSON(code, err)
		return
	}
	code, err = h.TaskStorage.ChangeTask(userTask)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при обновлении задачи.")
		c.JSON(code, err)
		return
	}
	log.Info().Msg("Запись успешно изменена.")
	c.JSON(http.StatusOK, "Запись успешно изменена.")
}

func (h *Handlers) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	code, err := h.TaskStorage.RemoveTask(id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удалении задачи.")
		c.JSON(code, err)
		return
	}
	log.Info().Msg("Запись успешно удалена.")
	c.JSON(http.StatusOK, "Запись успешно удалена.")
}
