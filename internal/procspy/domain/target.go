package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"
)

const DEFAULT_WEEKDAY_FACTOR = 0.5

type Target struct {
	User               string  `json:"user"`
	Name               string  `json:"name"`
	Pattern            string  `json:"pattern"`
	Limit              float64 `json:"limit"`
	LimitHours         float64 `json:"limit_hours,omitempty"`
	LimitWeekdays      float64 `json:"limit_weekdays,omitempty"`
	LimitHoursWeekDays float64 `json:"limit_hours_weekdays,omitempty"`
	Ocurrences         int     `json:"ocurrences,omitempty"`
	Elapsed            float64 `json:"elapsed,omitempty"`
	ElapsedHours       float64 `json:"elapsed_hours,omitempty"`
	FirstMatch         string  `json:"first_match,omitempty"`
	LastMatch          string  `json:"last_match,omitempty"`
	WarningOn          float64 `json:"warning_on"`
	Kill               bool    `json:"kill"`
	Source             string  `json:"source,omitempty"`
	CheckCommand       string  `json:"check_command,omitempty"`
	WarningCommand     string  `json:"warning_command,omitempty"`
	LimitCommand       string  `json:"limit_command,omitempty"`
	WeekdayFactor      float64 `json:"weekday_factor,omitempty"`
	rgx                *regexp.Regexp
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
		WeekdayFactor:  DEFAULT_WEEKDAY_FACTOR,
		rgx:            regexp.MustCompile(pattern),
	}

	if ret.rgx == nil {
		log.Printf("[domain.Target] Error compiling regex: %s to %s:%s", pattern, user, name)
	}

	ret.SetReportValues()

	return ret
}

func (t *Target) SetReportValues() {
	t.ElapsedHours = math.Round(t.Elapsed*100/3600) / 100
	t.LimitHours = math.Round(t.Limit*100/3600) / 100

	t.LimitWeekdays = t.Limit * t.WeekdayFactor
	t.LimitHoursWeekDays = math.Round(t.Limit*t.WeekdayFactor*100/3600) / 100
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

	for _, v := range ret.Targets {
		v.SetReportValues()
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
		ret += fmt.Sprintf("%s %s %s %f %f %t %s %s %s %s", v.User, v.Name, v.Pattern, v.getLimit(), v.getWarningOn(), v.Kill, v.Source, v.CheckCommand, v.WarningCommand, v.LimitCommand)
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

func (t *Target) AddMatchInfo(info *MatchInfo) {
	t.Elapsed += info.Elapsed
	t.FirstMatch = info.FirstMatch
	t.LastMatch = info.LastMatch
	t.Ocurrences = info.Ocurrences
	t.SetReportValues()
}

func (t *Target) AddElapsed(elapsed float64) {
	t.Elapsed += elapsed
	t.SetReportValues()
}

func (t *Target) SetElapsed(elapsed float64) {
	t.Elapsed = elapsed
	t.SetReportValues()
}

func (t *Target) ResetElapsed() {
	t.Elapsed = 0
	t.SetReportValues()
}

func (t *Target) CheckLimit() bool {
	if t.Limit == 0 {
		return false
	}

	return t.Elapsed > t.getLimit()
}

func (t *Target) CheckWarning() bool {
	if t.WarningOn == 0 {
		return false
	}

	return t.Elapsed > t.getWarningOn()
}

func (t *Target) applyFactor(limit float64) float64 {
	today := int(time.Now().Weekday())

	if t.WeekdayFactor <= 0 {
		t.WeekdayFactor = DEFAULT_WEEKDAY_FACTOR
	}

	if today == 0 || today == 6 {
		return limit
	}

	return t.LimitWeekdays
}

func (t *Target) getLimit() float64 {
	return t.applyFactor(t.Limit)
}

func (t *Target) getWarningOn() float64 {
	return t.applyFactor(t.WarningOn)
}
