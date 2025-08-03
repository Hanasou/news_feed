package core

import (
	"errors"
	"log"

	"github.com/Hanasou/news_feed/go/common/db"
	"github.com/Hanasou/news_feed/go/common/db/memdb"
	"github.com/Hanasou/news_feed/go/todo/models"
)

type TodoService struct {
	todoTable db.DbDriver
	userTable db.DbDriver
}

func CreateDb(dbType string, table string, rootPath string, saveToDisk bool) (db.DbDriver, error) {
	if dbType == "mem" {
		memDbDriver, err := memdb.Initialize(table, rootPath, saveToDisk)
		if err != nil {
			log.Printf("Could not initialize db. Error: %v", err)
			return nil, err
		}
		return memDbDriver, nil
	} else {
		return nil, errors.New("CreateDb in Todo service failed. Db type not supported: " + dbType)
	}
}

func InitializeService(dbType string, rootPath string, saveToDisk bool) (*TodoService, error) {
	// TODO: Get tables from config
	service := &TodoService{}
	todoDb, err := CreateDb(dbType, "todos", rootPath, saveToDisk)
	if err != nil {
		log.Printf("Could not create database driver for table: %s, %v", "todos", err)
		return nil, err
	}
	service.todoTable = todoDb
	userDb, err := CreateDb(dbType, "users", rootPath, saveToDisk)
	if err != nil {
		log.Printf("Could not create database driver for table: %s, %v", "users", err)
		return nil, err
	}
	service.userTable = userDb

	return service, nil
}

func (service *TodoService) CreateTodo(todo *models.Todo) error {
	err := service.todoTable.Upsert(todo)
	if err != nil {
		log.Printf("Create todo failed: %v", err)
		return err
	}
	log.Printf("Insert succeeded: %v", todo)
	return nil
}

func (service *TodoService) GetTodos(userId string) ([]*models.Todo, error) {
	todoData, err := service.todoTable.GetData()
	todos := []*models.Todo{}
	if err != nil {
		log.Printf("Failed to get todos: %v", err)
		return nil, err
	}
	// Verify that all data is instance of todo
	for _, v := range todoData {
		if concreteData, ok := v.(*models.Todo); ok {
			todos = append(todos, concreteData)
		} else {
			log.Printf("Could not convert this data into todo: %v", v)
		}
	}

	// Filter by userId
	var filteredTodos []*models.Todo
	for _, todo := range todos {
		if todo.UserId == userId {
			filteredTodos = append(filteredTodos, todo)
		}
	}
	return filteredTodos, nil
}
