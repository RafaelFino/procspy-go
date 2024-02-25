package procspy_storage

import (
	"errors"
	"log"
)

type CommandStorage struct {
	conn *DbConnection
}

func NewCommandStorage(dbConn *DbConnection) *CommandStorage {
	ret := &CommandStorage{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[CommandStorage] Error initializing storage: %s", err)
	}

	return ret
}

func (c *CommandStorage) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS command_log (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	name varchar(128) NOT NULL,
	command_type varchar(128) DEFAULT NULL,
	command TEXT NOT NULL,
	comand_return TEXT DEFAULT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP()
);	
	`

	if c.conn == nil {
		log.Printf("[CommandStorage] Error creating tables: db is nil")
		return errors.New("db is nil")
	}

	err := c.conn.Exec(create)

	if err != nil {
		log.Printf("[CommandStorage] Error creating tables: %s", err)
	}

	return err
}

func (c *CommandStorage) Close() error {
	if c.conn == nil {
		log.Printf("[CommandStorage] Database is already closed")
		return nil
	}

	return c.conn.Close()
}

func (c *CommandStorage) LogCommand(userID int, name string, commandType string, command string, commandReturn string) error {
	insert := `
INSERT INTO command_log 
(
	user_id, 
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
		log.Printf("[CommandStorage] Error logging command: db is nil")
		return errors.New("db is nil")
	}

	err := c.conn.Exec(insert, userID, name, commandType, command, commandReturn)

	if err != nil {
		log.Printf("[CommandStorage] Error logging command: %s")
	}

	return err
}
