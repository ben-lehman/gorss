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

func handlerLogin(state *State, cmd command) error {
	if len(cmd.args) != 1  {
    return fmt.Errorf("invalid login format. usage: %s <name>", cmd.name)
	}

  username := cmd.args[0]
	err := state.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("db url: %v\ncurrent user: %v\n", state.cfg.DbURL, state.cfg.CurrentUsername)
	return nil
}
