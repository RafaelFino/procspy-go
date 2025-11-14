package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type Match struct {
	storage *storage.Match
}

var MATCH_MAX_ELAPSED float64 = 120

func NewMatch(conn *storage.DbConnection) *Match {
	ret := &Match{
		storage: storage.NewMatch(conn),
	}

	log.Printf("[service.Match.NewMatch] Initializing match storage layer")

	err := ret.storage.Init()

	if err != nil {
		log.Printf("[service.Match.NewMatch] Failed to initialize match storage: %v", err)
		panic(err)
	}

	return ret
}

func (m *Match) Close() error {
	log.Printf("[service.Match.Close] Closing match storage connection")
	return m.storage.Close()
}

func (m *Match) InsertMatch(match *domain.Match) error {
	log.Printf("[service.Match.InsertMatch] Inserting match for pattern '%s' (user: '%s')", match.Pattern, match.User)

	if match.Elapsed > MATCH_MAX_ELAPSED {
		log.Printf("[service.Match.InsertMatch] Match elapsed time exceeds maximum allowed (%f > %f), capping to maximum", match.Elapsed, MATCH_MAX_ELAPSED)
		match.Elapsed = MATCH_MAX_ELAPSED
	}

	return m.storage.InsertMatch(match)
}

func (m *Match) GetMatches(user string) (map[string]float64, error) {
	data, err := m.storage.GetMatches(user)

	if err != nil {
		log.Printf("[service.Match.GetMatches] Failed to retrieve matches for user '%s': %v", user, err)
	}

	return data, err
}

func (m *Match) GetMatchesInfo(user string) (map[string]*domain.MatchInfo, error) {
	data, err := m.storage.GetMatchesInfo(user)

	if err != nil {
		log.Printf("[service.Match.GetMatchesInfo] Failed to retrieve match information for user '%s': %v", user, err)
	}

	return data, err
}
