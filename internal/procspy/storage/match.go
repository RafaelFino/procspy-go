package storage

import (
	"errors"
	"log"
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
	}

	return ret
}

func (m *Match) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS matches (
	id SERIAL PRIMARY KEY,
	when DATE DEFAULT CURRENT_DATE(),
	user varchar(128) REFERENCES users(id),
	name varchar(128) NOT NULL,
	pattern TEXT NOT NULL,
	match TEXT NOT NULL,
	elapsed REAL NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP()
);	
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

func (m *Match) InsertMatch(user string, name string, pattern string, match string, elapsed float64) error {
	insert := `
INSERT INTO matches
(
	user_id,
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

	err := m.conn.Exec(insert)

	if err != nil {
		log.Printf("[storage.Match] Error logging match: %s")
	}

	return err
}

func (m *Match) GetElapsed(user string) (map[string]float64, error) {
	query := `
SELECT
	name,
	sum(elapsed) elapsed
FROM
	matches
WHERE
	user = ?
	and when = current_date
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
