package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

type Target struct {
	User           string  `json:"user"`
	Name           string  `json:"name"`
	Pattern        string  `json:"pattern"`
	Limit          float64 `json:"limit"`
	Elapsed        float64 `json:"elapsed,omitempty"`
	WarningOn      float64 `json:"warning_on"`
	Kill           bool    `json:"kill"`
	Source         string  `json:"source,omitempty"`
	CheckCommand   string  `json:"check_command,omitempty"`
	WarningCommand string  `json:"warning_command,omitempty"`
	LimitCommand   string  `json:"limit_command,omitempty"`
	rgx            *regexp.Regexp
}

func NewTarget(user string, name string, pattern string, limit float64, warningOn float64, kill bool, source string, checkCommand string, warningCommand string, limitCommand string) *Target {
	ret := &Target{
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
		rgx:            regexp.MustCompile(pattern),
	}

	if ret.rgx == nil {
		log.Printf("[domain.Target] Error compiling regex: %s to %s:%s", pattern, user, name)
	}

	return ret
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

func NewTargetList() *TargetList {
	return &TargetList{
		Targets: []*Target{},
	}
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

func (t *TargetList) ToLog() string {
	ret, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		log.Printf("[domain.TargetList] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}

func (t *TargetList) Hash() string {
	ret := ""
	for _, v := range t.Targets {
		ret += fmt.Sprintf("%s %s %s %f %f %t %s %s %s %s", v.User, v.Name, v.Pattern, v.Limit, v.WarningOn, v.Kill, v.Source, v.CheckCommand, v.WarningCommand, v.LimitCommand)
	}
	return ret
}

func (t *Target) Match(value string) bool {
	if t.rgx == nil {
		t.rgx = regexp.MustCompile(t.Pattern)
	}

	ret := t.rgx.MatchString(value)

	return ret
}

func (t *Target) AddElapsed(elapsed float64) {
	t.Elapsed += elapsed
}

func (t *Target) ResetElapsed() {
	t.Elapsed = 0
}

func (t *Target) CheckLimit() bool {
	if t.Limit == 0 {
		return false
	}

	return t.Elapsed > t.Limit
}

func (t *Target) CheckWarning() bool {
	if t.WarningOn == 0 {
		return false
	}

	return t.Elapsed > t.WarningOn
}
