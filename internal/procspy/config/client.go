package config

import (
	"encoding/json"
	"log"
	"os"

	"procspy/internal/procspy/domain"
)

type Client struct {
	Interval  int             `json:"interval"`
	LogPath   string          `json:"log_path"`
	ServerURL string          `json:"server_url"`
	User      string          `json:"user"`
	Targets   []domain.Target `json:"targets"`
}

func NewConfig() *Client {
	return &Client{
		Interval: 30,
		LogPath:  "log/procspy.log",
		Targets:  []domain.Target{},
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

func (c *Client) ClearTargets() {
	c.Targets = []domain.Target{}
}

func (c *Client) AddTargets(targets []domain.Target) {
	c.Targets = append(c.Targets, targets...)
}

func (c *Client) AddTarget(t domain.Target) {
	c.Targets = append(c.Targets, t)
}

func ConfigClientFromFile(path string) (*Client, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return nil, err
	}

	return ConfigClientFromJson(string(byteValue))
}
