package procspy

import (
	"encoding/json"
	"log"
	"regexp"
)

type Target struct {
	Name    string `json:"name"`
	elapsed float64
	Limit   float64 `json:"limit"`
	Pattern string  `json:"pattern"`
	Kill    bool    `json:"kill"`
	Command string  `json:"command"`
	regex   *regexp.Regexp
}

func NewTarget(name string, limit float64, pattern string, kill bool, command string) *Target {
	return &Target{
		Name:    name,
		elapsed: 0,
		Limit:   limit,
		Pattern: pattern,
		Kill:    kill,
		Command: command,
		regex:   regexp.MustCompile(pattern),
	}
}

func (t *Target) GetName() string {
	return t.Name
}

func (t *Target) GetPattern() string {
	return t.Pattern
}

func (t *Target) GetCommand() string {
	return t.Command
}

func (t *Target) AddElapsed(elapsed float64) {
	t.elapsed += elapsed
}

func (t *Target) GetElapsed() float64 {
	return t.elapsed
}

func (t *Target) GetLimit() float64 {
	return t.Limit
}

func (t *Target) GetKill() bool {
	return t.Kill
}

func (t *Target) Match(command string) bool {
	if t.regex == nil {
		log.Printf("Trying to compile regex for target %s -> regex: [%s]", t.Name, t.Pattern)
		t.regex = regexp.MustCompile(t.Pattern)
	}

	if t.regex == nil {
		log.Printf("Error matching target %s: regex is nil", t.Name)
		return false
	}

	return t.regex.MatchString(command)
}

func (t *Target) IsExpired() bool {
	if t.GetLimit() <= 0 {
		return false
	}

	return t.GetElapsed() > t.GetLimit()
}

func (t *Target) ToJson() string {
	ret, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		log.Printf("Error marshalling target: %s", err)
	}

	return string(ret)
}

func TargetFromJson(jsonString string) (*Target, error) {
	ret := NewTarget("", 0, "", false, "")
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("Error unmarshalling target: %s", err)
		return nil, err
	}

	return ret, nil
}

func (t *Target) ResetElapsed() {
	t.elapsed = 0
}
