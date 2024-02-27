package config

import (
	"encoding/json"
	"log"
)

type Server struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c *Server) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "  ")
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
