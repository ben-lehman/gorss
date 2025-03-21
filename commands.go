package main

import (
	"fmt"
)


type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*State, command) error
}

func (c *commands) register(name string, f func(*State, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(state *State, cmd command) error {
	f, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("command handler has not been registered")
	}

	return f(state, cmd)
}



