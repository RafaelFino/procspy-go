package domain

import (
	"encoding/json"
	"log"
	"time"
)

type Match struct {
	User      string    `json:"user"`
	When      time.Time `json:"when"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Pattern   string    `json:"pattern"`
	Match     string    `json:"match"`
	Elapsed   float64   `json:"elapsed"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMatch(user string, name string, pattern string, match string, elapsed float64) *Match {
	return &Match{
		User:      user,
		Name:      name,
		Pattern:   pattern,
		Match:     match,
		Elapsed:   elapsed,
		When:      time.Now(),
		CreatedAt: time.Now(),
	}
}

func (m *Match) SetUser(user string) {
	m.User = user
}

func (m *Match) GetUser() string {
	return m.User
}

func (m *Match) GetName() string {
	return m.Name
}

func (m *Match) GetPattern() string {
	return m.Pattern
}

func (m *Match) GetMatch() string {
	return m.Match
}

func (m *Match) GetElapsed() float64 {
	return m.Elapsed
}

func (m *Match) SetCreatedAt(created_at time.Time) {
	m.CreatedAt = created_at
}

func (m *Match) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Match) SetWhen(when time.Time) {
	m.When = when
}

func (m *Match) GetWhen() time.Time {
	return m.When
}

func (m *Match) ToJson() string {
	ret, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("[Match] Error parsing json: %s", err)
		return ""
	}
	return string(ret)
}

func MatchFromJson(jsonString string) (*Match, error) {
	ret := &Match{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[Match] Error parsing json: %s", err)
		return nil, err
	}
	return ret, nil
}
