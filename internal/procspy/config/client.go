package config

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"procspy/internal/procspy/domain"
	"strings"
)

type ClientConfig struct {
	Interval    int    `json:"interval"`
	LogPath     string `json:"log_path"`
	ConfigUrl   string `json:"config_url"`
	LoadFromUrl bool   `json:"load_from_url"`
	Targets     []domain.Target
	localFile   string
	remoteCS    string
	onUpdate    chan bool
}

func NewConfig() *ClientConfig {
	ret := &ClientConfig{
		Interval:    60,
		LogPath:     "logs",
		ConfigUrl:   "http://rgt-tools.duckdns/config.json",
		Targets:     make([]domain.Target, 0),
		LoadFromUrl: false,
	}

	return ret
}

func InitConfig(filename string, onUpdate chan bool) (*ClientConfig, error) {
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

	ret, err := ConfigFromJson(jsonString)
	if err != nil {
		log.Printf("Error parsing json: %s", err)
		return nil, err
	}

	ret.localFile = filename

	err = ret.UpdateFromUrl()
	ret.onUpdate = onUpdate

	if err != nil {
		log.Printf("Error updating from url: %s", err)
	}

	return ret, nil
}

func (c *ClientConfig) SaveToFile(filename string) error {
	log.Printf("Saving config to file: %s", filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(c.ToJson())
	if err != nil {
		log.Printf("Error writing to file: %s", err)
		return err
	}

	return nil
}

func (c *ClientConfig) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Printf("Error marshalling config: %s", err)
	}

	return string(ret)
}

func ConfigFromJson(jsonString string) (*ClientConfig, error) {
	ret := NewConfig()
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("Error unmarshalling config: %s", err)
		return nil, err
	}

	return ret, nil
}

func (c *ClientConfig) GetRemoteCS() string {
	return c.remoteCS
}

func (c *ClientConfig) SetRemoteCS(cs string) {
	c.remoteCS = cs
}

func (c *ClientConfig) GetLocalFile() string {
	return c.localFile
}

func (c *ClientConfig) SetLocalFile(filename string) {
	c.localFile = filename
}

func (c *ClientConfig) calcCheckSum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (c *ClientConfig) UpdateFromUrl() error {
	if !c.LoadFromUrl {
		return nil
	}

	jsonString, err := c.downloadConfigFromURL()
	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return err
	}

	remoteCS := c.calcCheckSum(jsonString)
	if remoteCS != c.GetRemoteCS() {
		data := make(map[string][]domain.Target, 0)
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			log.Printf("Error unmarshalling config: %s from %s", err, jsonString)
			return err
		}

		c.Targets = make([]domain.Target, 0)

		for _, target := range data["targets"] {
			c.Targets = append(c.Targets, target)
			log.Printf("Updating target %s", target.GetName())
		}

		err = c.SaveToFile(c.GetLocalFile())
		if err != nil {
			log.Printf("Error saving config: %s", err)
			return err
		}

		c.SetRemoteCS(remoteCS)
		log.Printf("Config updated from url: [%s]\n%s", remoteCS, c.ToJson())
		go func() {
			c.onUpdate <- true
		}()
	}

	return nil
}

func (c *ClientConfig) downloadConfigFromURL() (string, error) {
	// Get the data
	resp, err := http.Get(c.ConfigUrl)
	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)

	if n == 0 {
		log.Printf("Error downloading config: %d bytes read: %s", n, err)
		return "", errors.New("No data read")
	}

	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return "", err
	}

	return buf.String(), err
}

func (c *ClientConfig) GetTargets() []domain.Target {
	return c.Targets
}

func (c *ClientConfig) GetInterval() int {
	return c.Interval
}
