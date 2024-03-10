package config

import (
	"encoding/json"
	"log"
)

type Server struct {
	APIPort    string `json:"api_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	LogPath    string `json:"log_path"`
}

func (s *Server) ToJson() string {
	ret, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("[Server] Error parsing json: %s", err)
		return ""
	}

	return string(ret)
}

func ServerFromJson(jsonString string) (*Server, error) {
	ret := &Server{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[Server] Error parsing json: %s", err)
		return nil, err
	}

	return ret, nil
}
