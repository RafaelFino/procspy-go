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
	remoteCS    string
	onUpdate    chan bool
}

func NewConfig() *Config {
	ret := &Config{
		Interval:    60,
		DBPath:      "data",
		LogPath:     "logs",
		ConfigUrl:   "http://rgt-tools.duckdns/config.json",
		Targets:     make([]Target, 0),
		LoadFromUrl: false,
	}

	return ret
}

func InitConfig(filename string, onUpdate chan bool) (*Config, error) {
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

	return ret, nil
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

func ConfigFromJson(jsonString string) (*Config, error) {
	ret := NewConfig()
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("Error unmarshalling config: %s", err)
		return nil, err
	}

	return ret, nil
}

func (c *Config) GetRemoteCS() string {
	return c.remoteCS
}

func (c *Config) SetRemoteCS(cs string) {
	c.remoteCS = cs
}

func (c *Config) GetLocalFile() string {
	return c.localFile
}

func (c *Config) SetLocalFile(filename string) {
	c.localFile = filename
}

func (c *Config) calcCheckSum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (c *Config) UpdateFromUrl() error {
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
		data := make(map[string][]Target, 0)
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			log.Printf("Error unmarshalling config: %s from %s", err, jsonString)
			return err
		}

		c.Targets = make([]Target, 0)

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

func (c *Config) downloadConfigFromURL() (string, error) {
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

func (c *Config) GetTargets() []Target {
	return c.Targets
}

func (c *Config) GetInterval() int {
	return c.Interval
}
