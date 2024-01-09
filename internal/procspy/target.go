package procspy

import (
	"encoding/json"
	"log"
)

type Target struct {
	Elapsed []float64 `json:"elapsed"`
	Limit   float64   `json:"limit"`
}

func NewTarget(limit float64) *Target {
	return &Target{
		Elapsed: make([]float64, 0),
		Limit:   limit,
	}
}

func (t *Target) AddElapsed(elapsed float64) {
	t.Elapsed = append(t.Elapsed, elapsed)
}

func (t *Target) GetElapsed() float64 {
	ret := 0.0

	for _, v := range t.Elapsed {
		ret += v
	}

	return ret
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

	if t.Elapsed == nil {
		t.Elapsed = make([]float64, 0)
	}

	return nil
}
