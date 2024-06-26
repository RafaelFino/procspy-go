package config

import (
	"encoding/json"
	"log"
	"os"
)

type Client struct {
	Interval  int    `json:"interval"`
	LogPath   string `json:"log_path"`
	ServerURL string `json:"server_url"`
	User      string `json:"user"`
	Debug     bool   `json:"debug,omitempty"`
}

func NewConfig() *Client {
	return &Client{
		Interval: 30,
		LogPath:  "log/procspy.log",
	}
}

func (c *Client) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Printf("[Client] Error marshalling config: %s", err)
	}

	return string(ret)
}

func ConfigClientFromJson(jsonString string) (*Client, error) {
	ret := &Client{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[Client] Error unmarshalling config: %s", err)
		return nil, err
	}

	log.Printf("Client config: %s", ret.ToJson())

	return ret, nil
}

func ConfigClientFromFile(path string) (*Client, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return nil, err
	}

	return ConfigClientFromJson(string(byteValue))
}
