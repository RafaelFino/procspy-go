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
		log.Printf("[storage.Match] Error initializing storage: %s", err)
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
		log.Printf("[storage.Match] Error creating tables: db is nil")
		return errors.New("db is nil")
	}

	err := m.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Match] Error creating tables: %s", err)
	}

	return err
}

func (m *Match) Close() error {
	if m.conn == nil {
		log.Printf("[storage.Match] Database is already closed")
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
		log.Printf("[storage.Match] Error logging match: db is nil")
		return errors.New("db is nil")
	}

	err := m.conn.Exec(insert, match.User, match.Name, match.Pattern, match.Match, match.Elapsed)

	if err != nil {
		log.Printf("[storage.Match] Error logging match: %s", err)
	}

	return err
}

func (m *Match) GetMatches(user string) (map[string]float64, error) {
	query := `
SELECT
	name,
	sum(elapsed) elapsed
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
		log.Printf("[storage.Match] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, user)

	if err != nil {
		log.Printf("[storage.Match] Error getting matches: %s", err)
		return nil, err
	}

	defer rows.Close()

	ret := make(map[string]float64)

	for rows.Next() {
		var name string
		var elapsed float64
		err = rows.Scan(&name, &elapsed)

		if err != nil {
			log.Printf("[storage.Match] Error scanning matches: %s", err)
			return nil, err
		}

		ret[name] = elapsed
	}

	return ret, nil
}
