package procspy

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Interval int               `json:"interval"`
	DBPath   string            `json:"db_path"`
	LogPath  string            `json:"log_path"`
	Targets  map[string]Target `json:"targets"`
}

func NewConfig() *Config {
	return &Config{
		Interval: 60,
		DBPath:   "data",
		LogPath:  "logs",
		Targets:  make(map[string]Target),
	}
}

func (c *Config) LoadFromFile(filename string) error {
	log.Println("Loading config from file: ", filename)

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

	for name, target := range c.Targets {
		target.FromJson(target.ToJson())
		c.Targets[name] = target
	}

	return nil
}

func (c *Config) AddTarget(name string, limit float64) {
	c.Targets[name] = *NewTarget(limit)
}

func (c *Config) GetTargets() map[string]Target {
	return c.Targets
}

func (c *Config) GetInterval() int {
	return c.Interval
}
