package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type Match struct {
	storage *storage.Match
}

func NewMatch(conn *storage.DbConnection) *Match {
	ret := &Match{
		storage: storage.NewMatch(conn),
	}

	log.Printf("[service.Match] Initializing storage")

	err := ret.storage.Init()

	if err != nil {
		log.Printf("[service.Match] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (m *Match) Close() error {
	log.Printf("[service.Match] Closing storage")
	return m.storage.Close()
}

func (m *Match) InsertMatch(match *domain.Match) error {
	log.Printf("[service.Match] Inserting match: %s", match.Pattern)
	return m.storage.InsertMatch(match)
}

func (m *Match) GetMatches(user string) (map[string]float64, error) {
	data, err := m.storage.GetMatches(user)

	if err != nil {
		log.Printf("[service.Match] Error getting matches: %s", err)
	}

	return data, err
}

func (m *Match) GetMatchesInfo(user string) (map[string]*domain.MatchInfo, error) {
	data, err := m.storage.GetMatchesInfo(user)

	if err != nil {
		log.Printf("[service.Match] Error getting matches info: %s", err)
	}

	return data, err
}
