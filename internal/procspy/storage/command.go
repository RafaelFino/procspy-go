package storage

import (
	"errors"
	"log"
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
	}

	return ret
}

func (c *Command) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS command_log (
	id SERIAL PRIMARY KEY,
	user varchar(128) REFERENCES users(id),
	name varchar(128) NOT NULL,
	command_type varchar(128) DEFAULT NULL,
	command TEXT NOT NULL,
	comand_return TEXT DEFAULT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP()
);	
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

func (c *Command) InsertCommand(user string, name string, commandType string, command string, commandReturn string) error {
	insert := `
INSERT INTO command_log 
(
	user, 
	name, 
	command_type, 
	command, 
	command_return
) 
VALUES 
(	?, 
	?, 
	?, 
	?, 
	?
);`
	if c.conn == nil {
		log.Printf("[storage.Command] Error logging command: db is nil")
		return errors.New("db is nil")
	}

	err := c.conn.Exec(insert, user, name, commandType, command, commandReturn)

	if err != nil {
		log.Printf("[storage.Command] Error logging command: %s")
	}

	return err
}
