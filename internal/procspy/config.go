package procspy

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
	"strings"
)

type Config struct {
	Interval    int      `json:"interval"`
	DBPath      string   `json:"db_path"`
	LogPath     string   `json:"log_path"`
	ConfigUrl   string   `json:"config_url"`
	Targets     []Target `json:"targets"`
	LoadFromUrl bool     `json:"load_from_url"`
	localFile   string
}

func NewConfig() *Config {
	return &Config{
		Interval:    60,
		DBPath:      "data",
		LogPath:     "logs",
		ConfigUrl:   "http://rgt-tools.duckdns/config.json",
		Targets:     make([]Target, 0),
		LoadFromUrl: false,
	}
}

func (c *Config) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return err
	}
	defer file.Close()

	jsonString := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jsonString += scanner.Text()
	}

	err = c.FromJson(jsonString)
	if err != nil {
		log.Printf("Error parsing json: %s", err)
		return err
	}

	c.localFile = filename

	return nil
}

func (c *Config) SaveToFile(filename string) error {
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

func (c *Config) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Printf("Error marshalling config: %s", err)
	}

	return string(ret)
}

func (c *Config) FromJson(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &c)
	if err != nil {
		log.Printf("Error unmarshalling config: %s", err)
		return err
	}

	for _, target := range c.Targets {
		target.FromJson(target.ToJson())
		c.Targets = append(c.Targets, target)
	}

	return nil
}

func (c *Config) calcCheckSum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (c *Config) FromUrl() error {
	if !c.LoadFromUrl {
		return nil
	}

	log.Printf("Loading config from url: %s", c.ConfigUrl)

	jsonString, err := c.downloadConfigFromURL()
	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return err
	}

	current := c.calcCheckSum(c.ToJson())
	fromUrl := c.calcCheckSum(jsonString)

	if fromUrl != current {
		log.Printf("Config changed, updating")
		err = c.FromJson(jsonString)
		if err != nil {
			log.Printf("Error parsing json: %s", err)
			return err
		}

		c.SaveToFile(c.localFile)
	}

	return nil
}

func (c *Config) downloadConfigFromURL() (string, error) {
	log.Printf("Downloading string from url: %s", c.ConfigUrl)

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
		log.Printf("Error downloading config: %d bytes read", err, n)
		return "", errors.New("No data read")
	}

	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return "", err
	}

	return buf.String(), err
}

func (c *Config) GetTargets() []Target {
	return c.Targets
}

func (c *Config) GetInterval() int {
	return c.Interval
}
