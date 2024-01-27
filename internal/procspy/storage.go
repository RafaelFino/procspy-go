package procspy

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DatabasePath string `json:"database_path"`
	db           *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	ret := &Storage{
		DatabasePath: path,
	}

	if err := os.Mkdir(path, 0755); !os.IsExist(err) {
		fmt.Printf("Error creating directory %s: %s", path, err)
		return nil, err
	}

	dbFile := fmt.Sprintf("%s/procspy.db", path)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Printf("Error opening database: %s on %s", err, dbFile)
		return nil, err
	}

	ret.db = db

	err = ret.preapreDatabase()

	if err != nil {
		log.Printf("Error creating tables: %s", err)
	}

	return ret, err
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		log.Printf("Error closing database: %s", err)
		return err
	}

	return nil
}

func (s *Storage) preapreDatabase() error {
	const command string = `
CREATE TABLE IF NOT EXISTS processes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,                
	name TEXT NOT NULL,
	pattern TEXT NOT NULL,
	command TEXT NOT NULL,
	kill BOOLEAN DEFAULT FALSE,
	elapsed REAL NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS processes_old (
	id INTEGER PRIMARY KEY AUTOINCREMENT,                
	name TEXT NOT NULL,
	kill INTEGER DEFAULT 0,
	elapsed REAL NOT NULL,
	created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS match (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	pattern TEXT NOT NULL,
	command TEXT NOT NULL,	
	kill BOOLEAN DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS match_old (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	pattern TEXT NOT NULL,
	command TEXT NOT NULL,
	kill BOOLEAN DEFAULT FALSE,
	created_at DATETIME NOT NULL
);

INSERT INTO processes_old 
SELECT
	min(id) id,
	name, 
	sum(kill) kill,
	sum(elapsed) elapsed,	
	date(created_at) created_at
FROM 
	processes 
WHERE created_at < datetime('now', '-1 day')
GROUP BY 
	name, date(created_at)
ORDER BY 
	created_at DESC;

DELETE FROM processes WHERE created_at < datetime('now', '-1 day');
DELETE FROM processes_old WHERE created_at < datetime('now', '-60 day');

INSERT INTO match_old
SELECT
	id,
	name,
	pattern,
	command,
	kill,
	created_at
FROM
	match
WHERE created_at < datetime('now', '-1 day')
ON CONFLICT DO NOTHING;

DELETE FROM match WHERE created_at < datetime('now', '-1 day');
DELETE FROM match_old WHERE created_at < datetime('now', '-60 day');
`
	_, err := s.db.Exec(command)
	if err != nil {
		log.Printf("Error creating table: %s", err)
	}

	return err
}

func (s *Storage) InsertProcess(name string, elapsed float64, pattern string, command string, kill bool) error {
	const sql string = `
INSERT INTO processes (name, elapsed, pattern, command, kill) VALUES (?, ?, ?, ?, ?);`
	_, err := s.db.Exec(sql, name, elapsed, pattern, command, kill)

	if err != nil {
		log.Printf("Error inserting process: %s", err)
	}

	return err
}

func (s *Storage) InsertMatch(name string, pattern string, command string, kill bool) error {
	const sql string = `
INSERT INTO match (name, pattern, command, kill) VALUES (?, ?, ?, ?);`
	_, err := s.db.Exec(sql, name, pattern, command, kill)

	if err != nil {
		log.Printf("Error inserting match: %s", err)
	}

	return err
}

func (s *Storage) GetElapsed() (map[string]float64, error) {
	const command string = `
SELECT name, SUM(elapsed) FROM processes GROUP BY name;`
	rows, err := s.db.Query(command)
	if err != nil {
		log.Printf("Error getting elapsed: %s", err)
		return nil, err
	}
	defer rows.Close()

	ret := make(map[string]float64)
	for rows.Next() {
		var name string
		var elapsed float64
		err = rows.Scan(&name, &elapsed)
		if err != nil {
			log.Printf("Error scanning elapsed: %s", err)
			return nil, err
		}
		ret[name] = elapsed
	}

	return ret, nil
}
