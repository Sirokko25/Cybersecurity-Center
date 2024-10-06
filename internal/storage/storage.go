package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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
		return DB{}, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS testtasks (
			task_id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			createdate TEXT NOT NULL,
			status TEXT NOT NULL
			);`)
	if err != nil {
		return DB{}, fmt.Errorf("failed to create table in db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return DB{conn: db}, nil
}

func (db *DB) AppendTask(task models.Task) (int, error) {
	query := `INSERT INTO testtasks (title, description, createdate, status) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status)
	fmt.Println(err)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Ошибка добавления записи")
	}
	return http.StatusOK, errors.New("Ошибка добавления записи")
}

func (db *DB) ChangeTask(task models.Task) (int, error) {
	query := `UPDATE testtasks SET title = $1, description= $2, createdate = $3, status = $4 WHERE task_id = $5`
	res, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status, task.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if rowsAffected == 0 {
		return http.StatusBadRequest, errors.New("Задача не найдена")
	}
	return http.StatusOK, nil
}

func (db *DB) FindTask(id string) (models.Task, int, error) {
	var task models.Task
	query := `SELECT task_id, title, description, createdate, status FROM testtasks WHERE task_id = $1`
	err := db.conn.QueryRow(query, id).Scan(&task.Id, &task.Title, &task.Description, &task.CreateDate, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, http.StatusBadRequest, errors.New("Задача не найдена")
		} else {
			return task, http.StatusBadRequest, errors.New("Ошибка выполнения запроса")
		}

	}
	return task, http.StatusOK, nil
}

func (db *DB) RemoveTask(id string) (int, error) {
	deleteQuery := `DELETE FROM testtasks WHERE task_id = $1`
	res, err := db.conn.Exec(deleteQuery, id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Ошибка выполнения запроса")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if rowsAffected == 0 {
		return http.StatusBadRequest, errors.New("Задача не найдена")
	}
	return http.StatusOK, nil
}
