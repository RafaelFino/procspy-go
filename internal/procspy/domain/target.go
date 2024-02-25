package procspy_domains

import (
	"encoding/json"
	"log"
	"regexp"
)

type Target struct {
	UserID     int     `json:"user_id"`
	Name       string  `json:"name"`
	Pattern    string  `json:"pattern"`
	Elapsed    float64 `json:"elapsed,omitempty"`
	Limit      float64 `json:"limit"`
	CheckCmd   string  `json:"check_cmd,omitempty"`
	WarnCmd    string  `json:"warn_cmd,omitempty"`
	ElapsedCmd string  `json:"elapsed_cmd,omitempty"`
	Kill       bool    `json:"kill"`
	SoSource   string  `json:"so_source,omitempty"`
	regex      *regexp.Regexp
}

func NewTarget(name string, limit float64, pattern string, kill bool) *Target {
	return &Target{
		Name:    name,
		Elapsed: 0,
		Limit:   limit,
		Pattern: pattern,
		Kill:    kill,
		regex:   regexp.MustCompile(pattern),
	}
}

func (t *Target) GetName() string {
	return t.Name
}

func (t *Target) GetPattern() string {
	return t.Pattern
}

func (t *Target) SetPattern(pattern string) {
	t.Pattern = pattern
}

func (t *Target) SetName(name string) {
	t.Name = name
}

func (t *Target) SetLimit(limit float64) {
	t.Limit = limit
}

func (t *Target) SetKill(kill bool) {
	t.Kill = kill
}

func (t *Target) SetElapsedCommand(command string) {
	t.ElapsedCmd = command
}

func (t *Target) SetCheckCommand(command string) {
	t.CheckCmd = command
}

func (t *Target) SetWarnCommand(command string) {
	t.WarnCmd = command
}

func (t *Target) GetCheckCommand() string {
	return t.CheckCmd
}

func (t *Target) GetWarnCommand() string {
	return t.WarnCmd
}

func (t *Target) SetElapsed(elapsed float64) {
	t.Elapsed = elapsed
}

func (t *Target) AddElapsed(elapsed float64) {
	t.Elapsed += elapsed
}

func (t *Target) GetElapsed() float64 {
	return t.Elapsed
}

func (t *Target) GetLimit() float64 {
	return t.Limit
}

func (t *Target) GetKill() bool {
	return t.Kill
}

func (t *Target) SetUserID(id int) {
	t.UserID = id
}

func (t *Target) GetElapsedCommand() string {
	return t.ElapsedCmd
}

func (t *Target) SetSoSource(soSource string) {
	t.SoSource = soSource
}

func (t *Target) GetSoSource() string {
	return t.SoSource
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
		log.Printf("[Target] Error marshalling target: %s", err)
	}

	return string(ret)
}

func TargetFromJson(jsonString string) (*Target, error) {
	ret := NewTarget("", 0, "", false, "")
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("[Target] Error unmarshalling target: %s", err)
		return nil, err
	}

	return ret, nil
}

func (t *Target) ResetElapsed() {
	t.Elapsed = 0
}
