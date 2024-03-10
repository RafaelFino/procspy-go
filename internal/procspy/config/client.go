package config

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Client struct {
	Interval  int    `json:"interval"`
	LogPath   string `json:"log_path"`
	ServerURL string `json:"server_url"`
	Key       string `json:"key"`
	User      string `json:"user"`
}

func NewConfig() *Client {
	ret := &Client{
		Interval: 60,
		LogPath:  "logs",
	}

	return ret
}

func LoadClientConfig(filename string) (*Client, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return nil, err
	}
	defer file.Close()

	jsonString := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jsonString += scanner.Text()
	}

	ret, err := ClientConfigFromJson(jsonString)
	if err != nil {
		log.Printf("Error parsing json: %s", err)
		return nil, err
	}

	return ret, nil
}

func (c *Client) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Printf("Error marshalling config: %s", err)
	}

	return string(ret)
}

func ClientConfigFromJson(jsonString string) (*Client, error) {
	ret := NewConfig()
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("Error unmarshalling config: %s", err)
		return nil, err
	}

	return ret, nil
}
