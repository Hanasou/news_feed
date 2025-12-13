package config

import (
	"log"

	"github.com/Hanasou/news_feed/go/common/parsers"
)

type GatewayConfig struct {
	Debug   bool          `json:"debug"`
	Clients ClientsConfig `json:"clients"`
}

type ClientsConfig struct {
	UserServiceUrl string `json:"user_service_url"`
}

func InitConfig() (*GatewayConfig, error) {
	result, err := parsers.ParseJSONFile("./config/gateway_config.json", &GatewayConfig{})
	if err != nil {
		log.Println("Error parsing gateway_config.json:", err)
		return nil, err
	}
	return result.(*GatewayConfig), nil
}
