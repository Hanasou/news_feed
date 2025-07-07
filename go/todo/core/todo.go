package core

import (
	"errors"
	"fmt"
	"log"

	"github.com/Hanasou/news_feed/go/common/db/memdb"
	"github.com/Hanasou/news_feed/go/todo/models"
)

type TodoService struct {
	todoTable *memdb.MemDb
	userTable *memdb.MemDb
}

func CreateDb(dbType string, table string, rootPath string) (*memdb.MemDb, error) {
	if dbType == "mem" {
		memDbDriver, err := memdb.Initialize(table)
		if err != nil {
			log.Fatalf("Could not initialize db. Shutting down: %v", err)
			return nil, err
		}
		memDbDriver.Data, err = memdb.GetDataFromFile(fmt.Sprintf("%s/%s.json", rootPath, table))
		if err != nil {
			log.Printf("Could not get file on disk for data initialization: %s", table)
		}
		return memDbDriver, nil
	} else {
		return nil, errors.New("CreateDb in Todo service failed. Db type not supported: " + dbType)
	}
}

func InitializeService(dbType string, rootPath string) (*TodoService, error) {
	// TODO: Get tables from config
	service := &TodoService{}
	todoDb, err := CreateDb(dbType, "todos", rootPath)
	if err != nil {
		log.Fatalf("Could not create database driver for table: %s, %v", "todos", err)
	}
	service.todoTable = todoDb
	userDb, err := CreateDb(dbType, "users", rootPath)
	if err != nil {
		log.Fatalf("Could not create database driver for table: %s, %v", "users", err)
	}
	service.userTable = userDb
	if err != nil {
		log.Fatalf("Failed to create db: %v", err)
		return nil, err
	}

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

func (service *TodoService) GetTodos() ([]*models.Todo, error) {
	data, err := service.todoTable.GetData()
	todos := []*models.Todo{}
	if err != nil {
		log.Printf("Failed to get todos: %v", err)
		return nil, err
	}
	// Verify that all data is instance of todo
	for _, v := range data {
		if concreteData, ok := v.(*models.Todo); ok {
			todos = append(todos, concreteData)
		} else {
			log.Printf("Could not convert this data into todo: %v", v)
		}
	}
	return todos, nil
}
