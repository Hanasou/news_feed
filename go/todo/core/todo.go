package core

import (
	"errors"
	"log"

	"github.com/Hanasou/news_feed/go/common/db"
	"github.com/Hanasou/news_feed/go/common/db/memdb"
	"github.com/Hanasou/news_feed/go/common/models"
)

type TodoService struct {
	todoTable db.DbDriver[*models.Todo]
}

func CreateDb(dbType string, table string, rootPath string, saveToDisk bool) (db.DbDriver[*models.Todo], error) {
	if dbType == "mem" {
		memDbDriver, err := memdb.Initialize[*models.Todo](table, rootPath, saveToDisk)
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
	todos, err := service.todoTable.GetData()
	if err != nil {
		log.Printf("Failed to get todos: %v", err)
		return nil, err
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
