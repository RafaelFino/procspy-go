package storage

import (
	"errors"
	"log"
	"procspy/internal/procspy/domain"
)

type Match struct {
	conn *DbConnection
}

func NewMatch(dbConn *DbConnection) *Match {
	ret := &Match{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Match.NewMatch] Failed to initialize match storage: %v", err)
		panic(err)
	}

	return ret
}

func (m *Match) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS matches (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT NOT NULL,
	name TEXT NOT NULL,
	pattern TEXT NOT NULL,
	match TEXT NOT NULL,
	elapsed int DEFAULT 60,
	created_at TIMESTAMP DEFAULT (datetime('now', 'localtime'))
);	

CREATE TABLE IF NOT EXISTS matches_old (
	id INTEGER,
	user TEXT NOT NULL,
	name TEXT NOT NULL,
	pattern TEXT NOT NULL,
	match TEXT NOT NULL,
	elapsed int DEFAULT 60,
	created_at TIMESTAMP DEFAULT (datetime('now', 'localtime'))
);

INSERT INTO matches_old
SELECT
	id,
	user,
	name,
	pattern,
	match,
	elapsed,
	created_at
FROM
	matches
WHERE
	date(created_at) <= date(date('now', 'localtime'), '-1 day')
ORDER BY 
	created_at DESC;

DELETE FROM matches
WHERE
	created_at < date(date('now', 'localtime'), '-1 day');
`
	if m.conn == nil {
		log.Printf("[storage.Match.Init] Cannot create tables: database connection is nil")
		return errors.New("db is nil")
	}

	err := m.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Match.Init] Failed to create match tables: %v", err)
	}

	return err
}

func (m *Match) Close() error {
	if m.conn == nil {
		log.Printf("[storage.Match.Close] Database connection is already closed")
		return nil
	}

	return m.conn.Close()
}

func (m *Match) InsertMatch(match *domain.Match) error {
	insert := `
INSERT INTO matches
(
	user,
	name,
	pattern,
	match,
	elapsed
)
VALUES
(	
	?,
	?,
	?,
	?,
	?
);`

	if m.conn == nil {
		log.Printf("[storage.Match.InsertMatch] Cannot insert match: database connection is nil")
		return errors.New("db is nil")
	}

	err := m.conn.Exec(insert, match.User, match.Name, match.Pattern, match.Match, match.Elapsed)

	if err != nil {
		log.Printf("[storage.Match.InsertMatch] Failed to insert match for user '%s', pattern '%s': %v", match.User, match.Pattern, err)
	}

	return err
}

func (m *Match) GetMatches(user string) (map[string]float64, error) {
	query := `
SELECT
	name,
	sum(elapsed) elapsed,
    min(created_at) first_elapsed,
    max(created_at) last_elapsed	
FROM
	matches
WHERE
	user = ?
	and date(created_at) >= date('now', 'localtime')
GROUP BY
	name
ORDER BY	
	name DESC;
`
	conn, err := m.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Match.GetMatches] Failed to get database connection: %v", err)
		return nil, err
	}

	rows, err := conn.Query(query, user)

	if err != nil {
		log.Printf("[storage.Match.GetMatches] Failed to query matches for user '%s': %v", user, err)
		return nil, err
	}

	defer rows.Close()

	ret := make(map[string]float64)

	for rows.Next() {
		var name string
		var elapsed float64
		var firstElapsed string
		var lastElapsed string
		err = rows.Scan(&name, &elapsed, &firstElapsed, &lastElapsed)

		if err != nil {
			log.Printf("[storage.Match.GetMatches] Failed to scan match row for user '%s': %v", user, err)
			return nil, err
		}

		ret[name] = elapsed
	}

	return ret, nil
}

func (m *Match) GetMatchesInfo(user string) (map[string]*domain.MatchInfo, error) {
	query := `
SELECT
	name,
	sum(elapsed) elapsed,
    min(created_at) first,
    max(created_at) last,
	count(*) ocurrences
FROM
	matches
WHERE
	user = ?
	and date(created_at) >= date('now', 'localtime')
GROUP BY
	name
ORDER BY	
	name DESC;
`
	conn, err := m.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Match.GetMatchesInfo] Failed to get database connection: %v", err)
		return nil, err
	}

	rows, err := conn.Query(query, user)

	if err != nil {
		log.Printf("[storage.Match.GetMatchesInfo] Failed to query match info for user '%s': %v", user, err)
		return nil, err
	}

	defer rows.Close()

	ret := make(map[string]*domain.MatchInfo)

	for rows.Next() {
		var name string

		var elapsed float64
		var first string
		var last string
		var ocurrences int

		err = rows.Scan(&name, &elapsed, &first, &last, &ocurrences)

		if err != nil {
			log.Printf("[storage.Match.GetMatchesInfo] Failed to scan match info row for user '%s': %v", user, err)
			return nil, err
		}

		ret[name] = &domain.MatchInfo{
			Elapsed:    elapsed,
			FirstMatch: first,
			LastMatch:  last,
			Ocurrences: ocurrences,
		}
	}

	return ret, nil
}
