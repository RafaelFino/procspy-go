package procspy

import (
	"encoding/json"
	"log"
)

type Target struct {
	Elapsed float64 `json:"elapsed"`
	Limit   float64 `json:"limit"`
}

func NewTarget(limit float64) *Target {
	return &Target{
		Elapsed: 0,
		Limit:   limit,
	}
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

func (t *Target) IsExpired() bool {
	return t.GetElapsed() > t.GetLimit()
}

func (t *Target) ToJson() string {
	ret, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		log.Printf("Error marshalling target: %s", err)
	}

	return string(ret)
}

func (t *Target) FromJson(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), &t)
	if err != nil {
		log.Printf("Error unmarshalling target: %s", err)
		return err
	}

	return nil
}
