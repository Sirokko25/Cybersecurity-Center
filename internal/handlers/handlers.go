package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sheduler/internal/handlers/helpers"
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
		c.JSON(http.StatusBadRequest, err)
		return
	}
	code, err := userTask.CheckingFields()
	if err != nil {
		c.JSON(code, err)
		return
	}
	code, err = h.TaskStorage.AppendTask(userTask)
	if err != nil {
		//log
		c.JSON(code, err)
		return
	}
	//log
	c.IndentedJSON(http.StatusOK, "Запись успешно добавлена.")
}

func (h *Handlers) GetTask(c *gin.Context) {
	var returnedTask models.Task
	id := c.Param("id")
	returnedTask, code, err := h.TaskStorage.FindTask(id)
	if err != nil {
		helpers.ErrorHandle(code, err)
	}
	//log
	c.IndentedJSON(http.StatusOK, returnedTask)
}

func (h *Handlers) PutTask(c *gin.Context) {
	var userTask models.Task
	err := c.BindJSON(&userTask)
	if err != nil {
		helpers.ErrorHandle(http.StatusBadRequest, err)
	}
	code, err := h.TaskStorage.ChangeTask(userTask)
	if err != nil {
		helpers.ErrorHandle(code, err)
	}
	helpers.ErrorHandle(code, err)
	c.JSON(http.StatusOK, "Запись успешно добавлена.")
}

func (h *Handlers) DeleteTask(c *gin.Context) {
	var userTask models.Task
	err := c.BindJSON(&userTask)
	if err != nil {
		helpers.ErrorHandle(http.StatusBadRequest, err)
	}
	code, err := h.TaskStorage.ChangeTask(userTask)
	if err != nil {
		helpers.ErrorHandle(code, err)
	}
	helpers.ErrorHandle(code, err)
	c.JSON(http.StatusOK, "Запись успешно добавлена.")
}
