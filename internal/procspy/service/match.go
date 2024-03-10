package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type Match struct {
	storage *storage.Match
	dbConn  *storage.DbConnection
}

func NewMatch(dbConn *storage.DbConnection) *Match {
	ret := &Match{
		dbConn:  dbConn,
		storage: storage.NewMatch(dbConn),
	}

	return ret
}

func (m *Match) InsertMatch(match *domain.Match) error {
	log.Printf("[service.Match] Inserting match: %s", match.Name)
	return m.storage.InsertMatch(match)
}

func (m *Match) GetMatches(user string) (map[string]float64, error) {
	log.Printf("[service.Match] Getting matches: %s", user)
	return m.storage.GetMatches(user)
}
