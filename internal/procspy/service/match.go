package service

import "procspy/internal/procspy/domain"

type Match struct {
}

func NewMatch() *Match {
	ret := &Match{}

	return ret
}

func (m *Match) InsertMatch(match *domain.Match) error {
	return nil
}

func (m *Match) GetElapsed(user string) (map[string]float64, error) {
	return nil, nil
}
