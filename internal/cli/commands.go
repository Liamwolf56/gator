package cli

import "fmt"

type Commands struct {
	handlers map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*State, Command) error),
	}
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}
