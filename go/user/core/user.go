package core

import (
	"errors"
	"log"

	"github.com/Hanasou/news_feed/go/common/db"
	"github.com/Hanasou/news_feed/go/common/db/memdb"
	"github.com/Hanasou/news_feed/go/user/auth"
	"github.com/Hanasou/news_feed/go/user/models"
)

type UserService struct {
	// Add fields for user service if needed
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

func InitializeService(dbType string, rootPath string, saveToDisk bool) (*UserService, error) {
	// TODO: Get tables from config
	service := &UserService{}
	userDb, err := CreateDb(dbType, "users", rootPath, saveToDisk)
	if err != nil {
		log.Printf("Could not create database driver for table: %s, %v", "users", err)
		return nil, err
	}
	service.userTable = userDb

	return service, nil
}

func (service *UserService) CreateUser(user *models.User) error {
	if user == nil {
		log.Println("Create user failed: user is nil")
		return errors.New("user cannot be nil")
	}
	if user.Username == "" || user.Password == "" || user.Email == "" {
		log.Println("Create user failed: missing required fields")
		return errors.New("user must have username, password, and email")
	}
	if !user.Role.IsValid() {
		user.Role = models.Default
	}

	// Hash password before storing
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Printf("Could not hash password: %v", err)
		return err
	}
	user.Password = hashedPassword

	// Insert user into the database
	log.Printf("Inserting user: %v", user)
	err = service.userTable.Upsert(user)
	if err != nil {
		log.Printf("Create user failed: %v", err)
		return err
	}
	log.Printf("Insert succeeded: %v", user)
	return nil
}
