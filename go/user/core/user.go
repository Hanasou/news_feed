package core

import (
	"errors"
	"log"

	"github.com/Hanasou/news_feed/go/common/auth"
	"github.com/Hanasou/news_feed/go/common/db"
	"github.com/Hanasou/news_feed/go/common/db/memdb"
	"github.com/Hanasou/news_feed/go/common/models"
	"github.com/Hanasou/news_feed/go/user/config"
)

type UserService struct {
	// Add fields for user service if needed
	userTable  db.DbDriver[*models.User]
	jwtService *auth.JWTService
}

func InitializeService(userServiceConfig *config.UserServiceConfig) (*UserService, error) {
	service := &UserService{}
	userDb, err := CreateDb(userServiceConfig.Database.Type, userServiceConfig.Database.Table,
		userServiceConfig.Database.RootPath, userServiceConfig.Database.SaveToDisk)
	if err != nil {
		log.Printf("Could not create database driver for table: %s, %v", "users", err)
		return nil, err
	}
	service.userTable = userDb

	return service, nil
}

func CreateDb(dbType string, table string, rootPath string, saveToDisk bool) (db.DbDriver[*models.User], error) {
	if dbType == "local" {
		memDbDriver, err := memdb.Initialize[*models.User](table, rootPath, saveToDisk)
		if err != nil {
			log.Printf("Could not initialize db. Error: %v", err)
			return nil, err
		}
		return memDbDriver, nil
	} else {
		return nil, errors.New("CreateDb in User service failed. Db type not supported: " + dbType)
	}
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

func (service *UserService) AuthenticateUser(userIdentifier, password string) (*auth.TokenPair, *models.User, error) {
	if userIdentifier == "" || password == "" {
		log.Println("Authenticate user failed: missing username or password")
		return nil, nil, errors.New("username and password must be provided")
	}

	user, err := service.userTable.GetByField("username", userIdentifier)
	if err != nil {
		log.Printf("Authenticate user failed: %v", err)
		return nil, nil, err
	}
	if user == nil {
		log.Println("Authenticate user failed: user not found")
		return nil, nil, nil
	}

	// Check password
	if err = auth.ValidatePassword(password, user.Password); err != nil {
		log.Println("Authenticate user failed: invalid password")
		return nil, nil, nil
	}

	tokenPair, err := service.jwtService.GenerateTokenPair(user)
	if err != nil {
		log.Printf("Could not generate token pair: %v", err)
		return nil, nil, err
	}

	log.Printf("User authenticated successfully: %v", user)
	return tokenPair, user, nil
}
