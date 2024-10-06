package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"sheduler/models"
)

type DB struct {
	conn *sql.DB
}

func ConnectionDB() (DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "Sirokko25", "Tasks_db")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{conn: nil}, err
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return DB{conn: db}, nil
}

func (db *DB) AppendTask(task models.Task) (int, error) {
	query := `INSERT INTO taskstorage (title, description, createdate, status) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Ошибка добавления записи")
	}
	return http.StatusOK, errors.New("Ошибка добавления записи")
}

func (db *DB) ChangeTask(task models.Task) (int, error) {
	query := `UPDATE tasks SET title = ?, description= ?, createdate = ?, status = ? WHERE id = ?`
	res, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status, task.Id)
	if err != nil {
		//log
		return http.StatusInternalServerError, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		//log
		return http.StatusInternalServerError, err
	}

	if rowsAffected == 0 {
		//log
		return http.StatusBadRequest, errors.New("Задача не найдена")
	}
	return http.StatusOK, nil
}

func (db *DB) FindTask(id string) (models.Task, int, error) {
	var task models.Task
	taskId, _ := strconv.Atoi("id")
	query := `SELECT id, title, description, createdate, status FROM tasks WHERE id = ?`
	err := db.conn.QueryRow(query, taskId).Scan(&task.Id, &task.Title, &task.Description, &task.CreateDate, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, http.StatusBadRequest, errors.New("Задача не найдена")
		} else {
			return task, http.StatusBadRequest, errors.New("Ошибка выполнения запроса")
		}

	}
	return task, http.StatusOK, nil
}

func (db *DB) DeleteTask(id string) (int, error) {
	deleteQuery := `DELETE FROM tasks WHERE id = ?`
	res, err := db.conn.Exec(deleteQuery, id)
	if err != nil {
		//log
		return http.StatusInternalServerError, errors.New("Ошибка выполнения запроса")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		//log
		return http.StatusInternalServerError, err
	}
	if rowsAffected == 0 {
		//log
		return http.StatusBadRequest, errors.New("Задача не найдена")
	}
	return http.StatusOK, nil
}
