package procspy_storage

import (
	"fmt"
	"log"
	procspy_domains "procspy/internal/procspy/domain"
)

type TargetStorage struct {
	conn *DbConnection
}

func NewTargetStorage(dbConn *DbConnection) *TargetStorage {
	ret := &TargetStorage{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[TargetStorage] Error initializing storage: %s", err)
	}

	return ret
}

func (t *TargetStorage) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS targets (     
	user_id INT REFERENCES users(id),
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
	PRIMARY KEY (user_id, name)
);	
	`

	if t.conn == nil {
		log.Printf("[TargetStorage] Error creating tables: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(create)

	if err != nil {
		log.Printf("[TargetStorage] Error creating tables: %s", err)
	}

	return err
}

func (t *TargetStorage) Close() error {
	if t.conn == nil {
		log.Printf("[TargetStorage] Database is already closed")
		return nil
	}

	return t.conn.Close()
}

func (t *TargetStorage) InsertTarget(target *procspy_domains.Target) error {
	insert := `
INSERT INTO targets (user_id, name, pattern, elapsed_cmd, check_cmd, warn_cmd, kill, so_source, limit) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
`
	if t.conn == nil {
		log.Printf("[TargetStorage] Error creating target: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(insert, target.UserID, target.Name, target.UserID, target.Name, target.Pattern, target.ElapsedCmd, target.CheckCmd, target.WarnCmd, target.Kill, target.SoSource, target.Limit)

	if err != nil {
		log.Printf("[TargetStorage] Error creating target: %s", err)
	}

	return err
}

func (t *TargetStorage) DeleteTargets(userID int) error {
	delete := `
DELETE FROM targets WHERE user_id = ?;	
`
	if t.conn == nil {
		log.Printf("[TargetStorage] Error deleting targets: db is nil")
		return fmt.Errorf("db is nil")
	}

	err := t.conn.Exec(delete, userID)

	if err != nil {
		log.Printf("[TargetStorage] Error deleting targets: %s", err)
	}

	return err
}

func (t *TargetStorage) GetTargets(userID int) (map[string]*procspy_domains.Target, error) {
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
	user_id = ?;
ORDER BY
	name;
`
	if t.conn == nil {
		log.Printf("[TargetStorage] Error getting targets: db is nil")
		return nil, fmt.Errorf("db is nil")
	}

	conn, err := t.conn.GetConn()

	if err != nil {
		log.Printf("[TargetStorage] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, userID)

	if err != nil {
		log.Printf("[TargetStorage] Error getting targets: %s", err)
		return nil, err
	}

	ret := make(map[string]*procspy_domains.Target)

	for rows.Next() {
		var name, pattern, elapsedCmd, checkCmd, warnCmd, soSource string
		var kill bool
		var limit float64

		err = rows.Scan(&name, &pattern, &elapsedCmd, &checkCmd, &warnCmd, &kill, &soSource, &limit)

		if err != nil {
			log.Printf("[TargetStorage] Error scanning targets: %s", err)
			return nil, err
		}

		target := procspy_domains.NewTarget(name, limit, pattern, kill)
		target.SetElapsedCommand(elapsedCmd)
		target.SetCheckCommand(checkCmd)
		target.SetWarnCommand(warnCmd)
		target.SetSoSource(soSource)
		target.SetUserID(userID)

		ret[name] = target
	}

	return ret, nil
}
