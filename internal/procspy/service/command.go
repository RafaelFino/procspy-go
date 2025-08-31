package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type Command struct {
	storage *storage.Command
}

func NewCommand(conn *storage.DbConnection) *Command {
	ret := &Command{
		storage: storage.NewCommand(conn),
	}

	log.Printf("[service.Command] Initializing storage")
	err := ret.storage.Init()

	if err != nil {
		log.Printf("[service.Command] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (c *Command) Close() error {
	log.Printf("[service.Command] Closing storage")
	return c.storage.Close()
}

func (c *Command) InsertCommand(cmd *domain.Command) error {
	log.Printf("[service.Command] Inserting command: %s", cmd.CommandLine)
	return c.storage.InsertCommand(cmd)
}

func (c *Command) GetCommands(user string) ([]*domain.Command, error) {
	log.Printf("[service.Command] Get commands from user: %s", user)
	return c.storage.GetCommands(user)
}
