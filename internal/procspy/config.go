package procspy

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Interval int               `json:"interval"`
	Database string            `json:"database"`
	Logfile  string            `json:"logfile"`
	Targets  map[string]Target `json:"targets"`
}

func NewConfig() *Config {
	return &Config{
		Interval: 60,
		Database: "procspy.db",
		Logfile:  "procspy.log",
		Targets:  make(map[string]Target),
	}
}

func (c *Config) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	return nil
}

func (c *Config) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(c.ToJson())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Config) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Println(err)
	}

	return string(ret)
}

func (c *Config) FromJson(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &c)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Config) AddTarget(name string, limit float64) {
	c.Targets = append(c.Targets, Target{name: name, limit: limit})
}

func (c *Config) GetTargets() map[string]Target {
	return c.Targets
}

func (c *Config) GetInterval() int {
	return c.Interval
}

func (c *Config) GetDatabase() string {
	return c.Database
}

func (c *Config) GetLogfile() string {
	return c.Logfile
}

func (c *Config) GetTarget(name string) Target {
	return c.Targets[name]
}

func (c *Config) GetTargetNames() []string {
	var names []string
	for _, target := range c.Targets {
		names = append(names, target.name)
	}
	return names
}

func (c *Config) GetTargetLimit(name string) float64 {
	return c.Targets[name].limit
}
