package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"sheduler/models"
)

type DB struct {
	conn *bun.DB
}

type storageStruct struct {
	bun.BaseModel `bun:"table:testtasks,alias:u"`

	Task_id     int64     `bun:",pk,autoincrement"`
	Title       string    `bun:",notnull"`
	Description string    `bun:",notnull"`
	Createdate  time.Time `bun:",nullzero,default:now()"`
	Status      string    `bun:",notnull"`
}

type StorageInterface interface {
	AppendTask(models.Task) (int64, error)
	ChangeTask(models.Task) (int64, error)
	FindTask(string) (models.Task, error)
	RemoveTask(string) (int64, error)
}

func ConnectionDB() (DB, error) {

	//connStr := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	"localhost", 5432, "postgres", "Sirokko25", "Tasks_db")
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	return DB{}, err
	//}
	//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS testtasks (
	//		task_id SERIAL PRIMARY KEY,
	//		title TEXT NOT NULL,
	//		description TEXT NOT NULL,
	//		createdate TIMESTAMP WITH TIME ZONE DEFAULT now(),
	//		status TEXT NOT NULL
	//		);`)
	//if err != nil {
	//	return DB{}, fmt.Errorf("failed to create table in db: %w", err)
	//}
	//err = db.Ping()
	//if err != nil {
	//	panic(err)
	//}

	dsn := "postgres://postgres:Sirokko25@localhost:5432/Tasks_db?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunConn := bun.NewDB(sqldb, pgdialect.New())
	err := bunConn.ResetModel(context.Background(), &storageStruct{})
	if err != nil {
		return DB{}, errors.New("Ошибка при создании таблицы в базе данных.")
	}
	fmt.Println("Соединение создано.")
	return DB{conn: bunConn}, nil
}

func (db DB) AppendTask(task models.Task) (int64, error) {
	storageObject := &storageStruct{Task_id: task.Id, Title: task.Title, Description: task.Description, Status: task.Status}
	var id int64
	//query := `INSERT INTO testtasks (title, description, createdate, status) VALUES ($1, $2, $3, $4)`
	//_, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status)\
	_, err := db.conn.NewInsert().Model(storageObject).Returning("task_id").Exec(context.Background(), &id)
	if err != nil {
		return 0, errors.New("Ошибка добавления записи")
	}
	return id, nil
}

func (db DB) ChangeTask(task models.Task) (int64, error) {
	storageObject := &storageStruct{}
	var idCounter int64
	//query := `UPDATE testtasks SET title = $1, description= $2, createdate = $3, status = $4 WHERE task_id = $5`
	//res, err := db.conn.Exec(query, task.Title, task.Description, task.CreateDate, task.Status, task.Id)
	res, err := db.conn.NewUpdate().Model(storageObject).Set("title = ?, description= ?, status = ?", task.Title, task.Description, task.Status).Where("task_id = ?", task.Id).Returning("task_id").Exec(context.Background(), &idCounter)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, errors.New("Задача не найдена.")
	}

	return idCounter, nil
}

func (db DB) FindTask(id string) (models.Task, error) {
	var task models.Task
	storageObject := &storageStruct{}
	//query := `SELECT task_id, title, description, createdate, status FROM testtasks WHERE task_id = $1`
	//err := db.conn.QueryRow(query, id).Scan(&task.Id, &task.Title, &task.Description, &task.CreateDate, &task.Status)
	err := db.conn.NewSelect().Model(storageObject).Where("task_id = ?", id).Scan(context.Background(), storageObject)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, errors.New("Задача не найдена")
		} else {
			return task, errors.New("Ошибка выполнения запроса")
		}
	}
	parseDate := storageObject.Createdate.Format("2006-01-02T15:04:05")
	task = models.Task{Id: storageObject.Task_id, Title: storageObject.Title, Description: storageObject.Description, CreateDate: parseDate, Status: storageObject.Status}
	return task, nil
}

func (db DB) RemoveTask(id string) (int64, error) {
	var idCounter []int64
	storageObject := &storageStruct{}
	//deleteQuery := `DELETE FROM testtasks WHERE task_id = $1`
	//res, err := db.conn.Exec(deleteQuery, id)
	_, err := db.conn.NewDelete().Model(storageObject).Where("task_id = ?", id).Returning("task_id").Exec(context.Background(), &idCounter)
	if err != nil {
		return 0, errors.New("Ошибка выполнения запроса")
	}
	if len(idCounter) == 0 {
		return 0, nil
	}
	return idCounter[0], nil
}
