package storage

import (
	"fmt"
	"log"

	domain "procspy/internal/procspy/domain"
)

type Target struct {
	conn *DbConnection
}

func NewTarget(dbConn *DbConnection) *Target {
	ret := &Target{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Target] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (t *Target) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS targets (     
	user varchar(128) REFERENCES users(name),
	name varchar(128) NOT NULL,
	pattern TEXT NOT NULL,
	elapsed_cmd TEXT NOT NULL,
	check_cmd TEXT NOT NULL,
	warn_cmd TEXT NOT NULL,
	elapsed_cmd TEXT NOT NULL,
	kill BOOLEAN DEFAULT FALSE,
	so_source TEXT DEFAULT NULL,
	limit REAL NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP(),
	PRIMARY KEY (user, name)
);	
	`

	if t.conn == nil {
		log.Printf("[storage.Target] Error creating tables: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Target] Error creating tables: %s", err)
	}

	return err
}

func (t *Target) Close() error {
	if t.conn == nil {
		log.Printf("[storage.Target] Database is already closed")
		return nil
	}

	return t.conn.Close()
}

func (t *Target) InsertTarget(target *domain.Target) error {
	insert := `
INSERT INTO targets (user, name, pattern, elapsed_cmd, check_cmd, warn_cmd, kill, so_source, limit)
`
	if t.conn == nil {
		log.Printf("[storage.Target] Error creating target: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(insert, target.GetUser(), target.GetName(), target.GetPattern(), target.GetElapsedCommand(), target.GetCheckCommand(), target.GetWarnCommand(), target.GetKill(), target.GetSoSource(), target.GetLimit())

	if err != nil {
		log.Printf("[storage.Target] Error creating target: %s", err)
	}

	return err
}

func (t *Target) DeleteTargets(user string) error {
	delete := `
DELETE FROM targets WHERE user = ?;	
`
	if t.conn == nil {
		log.Printf("[storage.Target] Error deleting targets: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(delete, user)

	if err != nil {
		log.Printf("[storage.Target] Error deleting targets: %s", err)
	}

	return err
}

func (t *Target) GetTargets(user string) (map[string]*domain.Target, error) {
	query := `
SELECT
	name,
	pattern,
	elapsed_cmd,
	check_cmd,
	warn_cmd,
	kill,
	so_source,
	limit
FROM
	targets
WHERE
	user = ?;
ORDER BY
	name;
`
	if t.conn == nil {
		log.Printf("[storage.Target] Error getting targets: db is nil")
		return nil, fmt.Errorf("db is nil")
	}

	conn, err := t.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Target] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, user)

	if err != nil {
		log.Printf("[storage.Target] Error getting targets: %s", err)
		return nil, err
	}

	ret := make(map[string]*domain.Target)

	for rows.Next() {
		var name, pattern, elapsedCmd, checkCmd, warnCmd, soSource string
		var kill bool
		var limit float64

		err = rows.Scan(&name, &pattern, &elapsedCmd, &checkCmd, &warnCmd, &kill, &soSource, &limit)

		if err != nil {
			log.Printf("[storage.Target] Error scanning targets: %s", err)
			return nil, err
		}

		target := domain.NewTarget(user, name, limit, pattern, kill)
		target.SetElapsedCommand(elapsedCmd)
		target.SetCheckCommand(checkCmd)
		target.SetWarnCommand(warnCmd)
		target.SetSoSource(soSource)

		ret[name] = target
	}

	return ret, nil
}
