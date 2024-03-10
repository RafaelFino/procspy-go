package service

import (
	"log"
	"procspy/internal/procspy/storage"
)

type Command struct {
	storage *storage.Command
	dbConn  *storage.DbConnection
}

func NewCommand(dbConn *storage.DbConnection) *Command {
	ret := &Command{
		dbConn:  dbConn,
		storage: storage.NewCommand(dbConn),
	}

	return ret
}

func (c *Command) InsertCommand(user string, name string, commandType string, command string, commandReturn string) error {
	log.Printf("[service.Command] Inserting command: %s", command)

	return c.storage.InsertCommand(user, name, commandType, command, commandReturn)
}
