package domain

import (
	"encoding/json"
	"log"
	"time"
)

type Match struct {
	User      string    `json:"user"`
	Name      string    `json:"name"`
	Pattern   string    `json:"pattern"`
	Match     string    `json:"match"`
	Elapsed   float64   `json:"elapsed"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type MatchList struct {
	Matches map[string]float64 `json:"matches"`
}

func NewMatch(user string, name string, pattern string, match string, elapsed float64) *Match {
	ret := &Match{
		User:      user,
		Name:      name,
		Pattern:   pattern,
		Match:     match,
		Elapsed:   elapsed,
		CreatedAt: time.Now(),
	}

	return ret
}

func (m *Match) ToLog() string {
	ret, err := json.Marshal(m)
	if err != nil {
		log.Printf("[domain.Match] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}

func (m *Match) ToJson() string {
	ret, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("[domain.Match] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}

func MatchFromJson(jsonString string) (*Match, error) {
	ret := &Match{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[domain.Match] Error parsing json: %s", err)
		return nil, err
	}
	return ret, nil
}

func MatchListFromJson(jsonString string) (*MatchList, error) {
	ret := &MatchList{}
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("[domain.Match] Error parsing json: %s", err)
		return nil, err
	}
	return ret, nil
}
