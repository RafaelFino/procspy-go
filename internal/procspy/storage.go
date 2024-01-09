package procspy

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	DatabasePath string `json:"database_path"`
	db *sql.DB	
)

func NewStorage(databasePath string) *Storage {
	return &Storage{
		DatabasePath: databasePath,
	}
}

func (s *Storage) Connect() error {
	db, err := sql.Open("sqlite3", s.DatabasePath)
	if err != nil {
		log.Println(err)
		return err
	}
	s.db = db

	return nil
}

func (s *Storage) Disconnect() error {
	err := s.db.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) CreateTables() error {
	const command string = `
CREATE TABLE IF NOT EXISTS processes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,                
	name TEXT NOT NULL,
	elapsed REAL NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	_, err := s.db.Exec(command)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) InsertProcess(name string, elapsed float64) error {
	const command string = `
INSERT INTO processes (name, elapsed) VALUES (?, ?);`
	_, err := s.db.Exec(command, name, elapsed)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) GetElapsed() (map[string]float64, error) {
	const command string = `
SELECT name, SUM(elapsed) FROM processes GROUP BY name;`	
	rows, err := s.db.Query(command)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	ret := make(map[string]float64)
	for rows.Next() {
		var name string
		var elapsed float64
		err = rows.Scan(&name, &elapsed)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		ret[name] = elapsed
	}
	
	return ret, nil
}
