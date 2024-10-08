package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
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

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DSN"))))
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
	_, err := db.conn.NewInsert().Model(storageObject).Returning("task_id").Exec(context.Background(), &id)
	if err != nil {
		return 0, errors.New("Ошибка добавления записи")
	}
	return id, nil
}

func (db DB) ChangeTask(task models.Task) (int64, error) {
	storageObject := &storageStruct{}
	var idCounter int64
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
	err := db.conn.NewSelect().Model(storageObject).Where("task_id = ?", id).Scan(context.Background(), storageObject)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, errors.New("Задача не найдена")
		} else {
			return task, errors.New("Ошибка выполнения запроса")
		}
	}
	parseDate := storageObject.Createdate.Format(os.Getenv("TIME_FORMAT"))
	task = models.Task{Id: storageObject.Task_id, Title: storageObject.Title, Description: storageObject.Description, CreateDate: parseDate, Status: storageObject.Status}
	return task, nil
}

func (db DB) RemoveTask(id string) (int64, error) {
	var idCounter int64 = 0
	storageObject := &storageStruct{}
	_, err := db.conn.NewDelete().Model(storageObject).Where("task_id = ?", id).Returning("task_id").Exec(context.Background(), &idCounter)
	if err != nil {
		return 0, errors.New("Ошибка выполнения запроса")
	}
	if idCounter == 0 {
		return 0, nil
	}
	return idCounter, nil
}
