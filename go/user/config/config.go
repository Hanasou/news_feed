package config

import (
	"log"
	"os"

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

const configName = "user_service_config.json"

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/"
	}
	return configPath + configName
}

func InitConfig() (*UserServiceConfig, error) {
	result, err := parsers.ParseJSONFile(getConfigPath(), &UserServiceConfig{})
	if err != nil {
		log.Printf("Error parsing %s: %s", configName, err)
		return nil, err
	}
	return result.(*UserServiceConfig), nil
}
