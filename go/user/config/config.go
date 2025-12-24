package config

import (
	"log"

	"github.com/Hanasou/news_feed/go/common/parsers"
)

type UserServiceConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

type DatabaseConfig struct {
	Type       string `json:"type"`
	RootPath   string `json:"root_path"`
	SaveToDisk bool   `json:"save_to_disk"`
	Table      string `json:"table"`
}

type ServerConfig struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func InitConfig() (*UserServiceConfig, error) {
	result, err := parsers.ParseJSONFile("/app/go/user/config/user_config.json", &UserServiceConfig{})
	if err != nil {
		log.Println("Error parsing user_config.json:", err)
		return nil, err
	}
	return result.(*UserServiceConfig), nil
}
