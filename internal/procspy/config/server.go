package config

import (
	"encoding/json"
	"log"
	"os"
)

type Server struct {
	DBPath     string            `json:"db_path"`
	LogPath    string            `json:"log_path"`
	APIPort    int               `json:"api_port"`
	APIHost    string            `json:"api_host"`
	UserTarges map[string]string `json:"user_targets"`
	Debug      bool              `json:"debug"`
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ToJson() string {
	ret, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Printf("[config.Server.ToJson] Failed to marshal server configuration to JSON: %v", err)
	}

	return string(ret)
}

func ServerConfigFromJson(jsonString string) (*Server, error) {
	ret := &Server{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[config.ServerConfigFromJson] Failed to unmarshal server configuration: %v", err)
		return nil, err
	}

	log.Printf("[config.ServerConfigFromJson] Server configuration loaded successfully: %s", ret.ToJson())

	return ret, nil
}

func ServerConfigFromFile(path string) (*Server, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[config.ServerConfigFromFile] Failed to read configuration file '%s': %v", path, err)
		return nil, err
	}

	return ServerConfigFromJson(string(byteValue))
}
