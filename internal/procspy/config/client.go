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
	APIPort   int    `json:"api_port,omitempty"`
	APIHost   string `json:"api_host,omitempty"`
}

func NewConfig() *Client {
	return &Client{
		Interval: 30,
		LogPath:  "logs",
		Debug:    false,
		APIPort:  8888,
		APIHost:  "localhost",
	}
}

func (c *Client) SetDefaults() {
	if c.Interval < 30 {
		c.Interval = 30
	}

	if c.LogPath == "" {
		c.LogPath = "logs"
	}

	if c.APIPort == 0 {
		c.APIPort = 8888
	}

	if c.APIHost == "" {
		c.APIHost = "localhost"
	}
}

func (c *Client) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Printf("[config.Client.ToJson] Failed to marshal client configuration to JSON: %v", err)
	}

	return string(ret)
}

func ClientConfigFromJson(jsonString string) (*Client, error) {
	ret := &Client{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[config.ClientConfigFromJson] Failed to unmarshal client configuration: %v", err)
		return nil, err
	}

	ret.SetDefaults()

	log.Printf("[config.ClientConfigFromJson] Client configuration loaded successfully: %s", ret.ToJson())

	return ret, nil
}

func ClientConfigFromFile(path string) (*Client, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[config.ClientConfigFromFile] Failed to read configuration file '%s': %v", path, err)
		return nil, err
	}

	return ClientConfigFromJson(string(byteValue))
}
