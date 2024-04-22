package config

import (
	"encoding/json"
	"log"
)

type Server struct {
	DBPath     string            `json:"db_path"`
	LogPath    string            `json:"log_path"`
	APIPort    string            `json:"api_port"`
	UserTarges map[string]string `json:"user_targets"`
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

	return ret, nil
}
