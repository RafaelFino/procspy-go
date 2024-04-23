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
		log.Printf("Error marshalling target: %s", err)
	}

	return string(ret)
}

func ConfigServerFromJson(jsonString string) (*Server, error) {
	ret := &Server{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("Error unmarshalling target: %s", err)
		return nil, err
	}

	log.Printf("Server config: %s", ret.ToJson())

	return ret, nil
}

func ConfigServerFromFile(path string) (*Server, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return nil, err
	}

	return ConfigServerFromJson(string(byteValue))
}
