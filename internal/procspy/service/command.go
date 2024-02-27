package service

type Command struct {
}

func NewCommand() *Command {
	ret := &Command{}

	return ret
}

func (c *Command) InsertCommand(user string, name string, commandType string, command string, commandReturn string) error {
	return nil
}
