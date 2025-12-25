package config

import (
	"log"
	"os"

	"github.com/Hanasou/news_feed/go/common/parsers"
)

type GatewayConfig struct {
	Debug   bool          `json:"debug"`
	Clients ClientsConfig `json:"clients"`
}

type ClientsConfig struct {
	UserClientConfig UserClientConfig `json:"user_client_config"`
}

type UserClientConfig struct {
	Protocol    string `json:"protocol"`
	ServiceHost string `json:"service_host"`
	ServicePort int    `json:"service_port"`
}

const configName = "gateway_config.json"

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/"
	}
	return configPath + configName
}

func InitConfig() (*GatewayConfig, error) {
	result, err := parsers.ParseJSONFile(getConfigPath(), &GatewayConfig{})
	if err != nil {
		log.Printf("Error parsing %s: %s", configName, err)
		return nil, err
	}
	return result.(*GatewayConfig), nil
}
