package domain

import (
	"encoding/json"
	"log"
)

type Target struct {
	User           string `json:"user"`
	Name           string `json:"name"`
	Pattern        string `json:"pattern"`
	Limit          int    `json:"limit"`
	WarningOn      int    `json:"warning_on"`
	Kill           bool   `json:"kill"`
	Source         string `json:"source,omitempty"`
	CheckCommand   string `json:"check_command,omitempty"`
	WarningCommand string `json:"warning_command,omitempty"`
	LimitCommand   string `json:"limit_command,omitempty"`
}

func NewTarget(user string, name string, pattern string, limit int, warningOn int, kill bool, source string, checkCommand string, warningCommand string, limitCommand string) *Target {
	return &Target{
		User:           user,
		Name:           name,
		Pattern:        pattern,
		Limit:          limit,
		WarningOn:      warningOn,
		Kill:           kill,
		Source:         source,
		CheckCommand:   checkCommand,
		WarningCommand: warningCommand,
		LimitCommand:   limitCommand,
	}
}

func (t *Target) ToLog() string {
	ret, err := json.Marshal(t)
	if err != nil {
		log.Printf("[domain.Target] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}
func (t *Target) ToJson() string {
	ret, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		log.Printf("[domain.Target] Error parsing json: %s", err)
	}

	return string(ret)
}

type TargetList struct {
	Targets []*Target `json:"targets"`
}

func TargetListFromJson(jsonString string) (*TargetList, error) {
	ret := &TargetList{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[domain.Target] Error parsing json: %s", err)
		return nil, err
	}

	return ret, nil
}
