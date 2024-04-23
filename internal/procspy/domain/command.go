package domain

import (
	"encoding/json"
	"log"
	"time"
)

type Command struct {
	User        string    `json:"user"`
	Name        string    `json:"name"`
	CommandLine string    `json:"command_line"`
	Return      string    `json:"command_return"`
	Source      string    `json:"source"`
	CommandLog  string    `json:"command_log"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

func NewCommand(user string, name string, commandLine string, commandReturn string) *Command {
	return &Command{
		User:        user,
		Name:        name,
		CommandLine: commandLine,
		Return:      commandReturn,
		Source:      "procspy",
		CommandLog:  "",
		CreatedAt:   time.Now(),
	}
}

func (c *Command) ToLog() string {
	ret, err := json.Marshal(c)
	if err != nil {
		log.Printf("[domain.Command] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}
func (c *Command) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("[domain.Command] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}

func CommandFromJson(jsonString string) (*Command, error) {
	ret := &Command{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[domain.Command] Error parsing json: %s", err)
		return nil, err
	}
	return ret, nil
}
