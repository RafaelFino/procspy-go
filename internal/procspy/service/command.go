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

	log.Printf("[service.Command.NewCommand] Initializing command storage layer")
	err := ret.storage.Init()

	if err != nil {
		log.Printf("[service.Command.NewCommand] Failed to initialize command storage: %v", err)
		panic(err)
	}

	return ret
}

func (c *Command) Close() error {
	log.Printf("[service.Command.Close] Closing command storage connection")
	return c.storage.Close()
}

func (c *Command) InsertCommand(cmd *domain.Command) error {
	log.Printf("[service.Command.InsertCommand] Inserting command '%s' for user '%s'", cmd.CommandLine, cmd.User)
	return c.storage.InsertCommand(cmd)
}

func (c *Command) GetCommands(user string) ([]*domain.Command, error) {
	log.Printf("[service.Command.GetCommands] Retrieving commands for user '%s'", user)
	return c.storage.GetCommands(user)
}
