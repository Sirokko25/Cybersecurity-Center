package models

import (
	"errors"
	"net/http"
)

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreateDate  string `json:"createdate"`
	Status      string `json:"status"`
}

func (t *Task) CheckingFields() (int, error) {
	if t.Title == "" || t.Description == "" || t.CreateDate == "" || t.Status == "" {
		return http.StatusBadRequest, errors.New("Некорректно указаны данные")
	}
	return 0, nil
}
