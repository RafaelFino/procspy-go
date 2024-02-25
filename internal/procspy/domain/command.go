package procspy_domains

import (
	"encoding/json"
	"time"
)

type Command struct {
	UserID        int       `json:"user_id"`
	Name          string    `json:"name"`
	Command       string    `json:"command"`
	CommandReturn string    `json:"command_return"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewCommand(name string, command string) *Command {
	return &Command{
		Name:      name,
		Command:   command,
		CreatedAt: time.Now(),
	}
}

func (c *Command) SetUserID(id int) {
	c.UserID = id
}

func (c *Command) SetCommandReturn(commandReturn string) {
	c.CommandReturn = commandReturn
}

func (c *Command) GetUserID() int {
	return c.UserID
}

func (c *Command) GetName() string {
	return c.Name
}

func (c *Command) GetCommand() string {
	return c.Command
}

func (c *Command) GetCommandReturn() string {
	return c.CommandReturn
}

func (c *Command) GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c *Command) SetCreatedAt(created_at time.Time) {
	c.CreatedAt = created_at
}

func (c *Command) ToJson() string {
	ret, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return ""
	}
	return string(ret)
}

func CommandFromJson(jsonString string) (*Command, error) {
	ret := &Command{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
