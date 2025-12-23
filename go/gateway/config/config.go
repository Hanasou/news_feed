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
	UserClientConfig UserClientConfig `json:"user_client_config"`
}

type UserClientConfig struct {
	Protocol    string `json:"protocol"`
	ServiceHost string `json:"service_host"`
	ServicePort int    `json:"service_port"`
}

func InitConfig() (*GatewayConfig, error) {
	result, err := parsers.ParseJSONFile("/app/go/gateway/config/gateway_config.json", &GatewayConfig{})
	if err != nil {
		log.Println("Error parsing gateway_config.json:", err)
		return nil, err
	}
	return result.(*GatewayConfig), nil
}
