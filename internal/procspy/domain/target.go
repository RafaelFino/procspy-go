package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"
)

const DEFAULT_WEEKDAY_LIMIT = 0.5
const DEFAULT_WEEKEND_LIMIT = 1.0
const DEFAULT_BASE_LIMIT = 60 * 60
const DEFAULT_WARNING_ON = 0.95

type Target struct {
	User           string          `json:"user"`
	Name           string          `json:"name"`
	Pattern        string          `json:"pattern"`
	Source         string          `json:"source,omitempty"`
	Limit          float64         `json:"limit"`
	Elapsed        float64         `json:"elapsed,omitempty"`
	Remaining      float64         `json:"remaining"`
	Ocurrences     int             `json:"ocurrences,omitempty"`
	FirstMatch     string          `json:"first_match,omitempty"`
	LastMatch      string          `json:"last_match,omitempty"`
	Kill           bool            `json:"kill"`
	LimitCommand   string          `json:"limit_command,omitempty"`
	CheckCommand   string          `json:"check_command,omitempty"`
	WarningCommand string          `json:"warning_command,omitempty"`
	WarningOn      float64         `json:"warning_on,omitempty"`
	Weekdays       map[int]float64 `json:"weekdays,omitempty"`
	rgx            *regexp.Regexp
}

func (t *Target) setWeekdays() {
	if t.Weekdays == nil {
		t.Weekdays = map[int]float64{}
	}

	for i := 0; i < 7; i++ {
		if _, found := t.Weekdays[i]; !found {
			if i == 0 || i == 6 {
				t.Weekdays[i] = DEFAULT_WEEKEND_LIMIT
			} else {
				t.Weekdays[i] = DEFAULT_WEEKDAY_LIMIT
			}
		}
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
		v.setWeekdays()
		v.getLimit()
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
	t.Remaining = t.getLimit() - t.Elapsed
}

func (t *Target) AddElapsed(elapsed float64) {
	t.Elapsed += elapsed
	t.Remaining = t.getLimit() - t.Elapsed
}

func (t *Target) SetElapsed(elapsed float64) {
	t.Elapsed = elapsed
	t.Remaining = t.getLimit() - elapsed
}

func (t *Target) ResetElapsed() {
	t.Elapsed = 0
	t.Remaining = t.getLimit()
}

func (t *Target) CheckLimit() bool {
	limit := t.getLimit()
	if limit == 0 {
		return false
	}

	return t.Elapsed >= limit
}

func (t *Target) CheckWarning() bool {
	warn := t.getWarningOn()
	if warn == 0 {
		return false
	}

	return t.Elapsed > warn
}

func (t *Target) getLimit() float64 {
	today := int(time.Now().Weekday())

	factor, found := t.Weekdays[today]

	if !found {
		factor = DEFAULT_WEEKDAY_LIMIT
	}

	t.Limit = DEFAULT_BASE_LIMIT * factor
	return t.Limit
}

func (t *Target) getWarningOn() float64 {
	t.WarningOn = t.getLimit() * DEFAULT_WARNING_ON
	return t.WarningOn
}
