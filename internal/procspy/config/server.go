package config

import (
	"encoding/json"
	"log"
)

type ServerConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c *ServerConfig) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("[ServerConfig] Error parsing json: %s", err)
		return ""
	}

	return string(ret)
}

func ServerConfigFromJson(jsonString string) (*ServerConfig, error) {
	ret := &ServerConfig{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[ServerConfig] Error parsing json: %s", err)
		return nil, err
	}

	return ret, nil
}
