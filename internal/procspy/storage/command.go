package storage

import (
	"errors"
	"log"
	"procspy/internal/procspy/domain"
)

type Command struct {
	conn *DbConnection
}

func NewCommand(dbConn *DbConnection) *Command {
	ret := &Command{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Command] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (c *Command) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS command_log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT NOT NULL,
	name TEXT NOT NULL,
	command_line TEXT NOT NULL,
	command_return TEXT DEFAULT NULL,
	source TEXT NOT NULL,
	command_log TEXT DEFAULT NULL,
	created_at TIMESTAMP DEFAULT (datetime('now', 'localtime'))
);	

CREATE TABLE IF NOT EXISTS command_log_old (
	id INTEGER,
	user TEXT NOT NULL,
	name TEXT NOT NULL,
	command_line TEXT NOT NULL,
	command_return TEXT DEFAULT NULL,
	source TEXT NOT NULL,
	command_log TEXT DEFAULT NULL,
	created_at TIMESTAMP DEFAULT (datetime('now', 'localtime'))
);

INSERT INTO command_log_old
SELECT
	id,
	user,
	name,
	command_line,
	command_return,
	source,
	command_log,
	created_at
FROM	
	command_log
WHERE
	date(created_at) < date(date('now', 'localtime'), '-1 day')
ORDER BY
	created_at DESC;

DELETE FROM command_log
WHERE
	date(created_at) < date(date('now', 'localtime'), '-1 day');
	`
	if c.conn == nil {
		log.Printf("[storage.Command] Error creating tables: db is nil")
		return errors.New("db is nil")
	}

	err := c.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Command] Error creating tables: %s", err)
	}

	return err
}

func (c *Command) Close() error {
	if c.conn == nil {
		log.Printf("[storage.Command] Database is already closed")
		return nil
	}

	return c.conn.Close()
}

func (c *Command) InsertCommand(cmd *domain.Command) error {
	insert := `
INSERT INTO command_log (
	user, 
	name, 
	command_line, 
	command_return, 
	source, 
	command_log)
VALUES 
	(?, ?, ?, ?, ?, ?)
`
	err := c.conn.Exec(insert, cmd.User, cmd.Name, cmd.CommandLine, cmd.Return, cmd.Source, cmd.CommandLog)

	if err != nil {
		log.Printf("[storage.Command] Error executing query: %s -> error: %s", insert, err)
	}

	return err
}
