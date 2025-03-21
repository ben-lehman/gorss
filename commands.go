package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ben-lehman/gorss/internal/database"
	"github.com/google/uuid"
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
 
  _, err := state.db.GetUser(context.Background(), username)
  if err != nil {
    return fmt.Errorf("No user with that name: %v", err)
  }

	err = state.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("db url: %v\ncurrent user: %v\n", state.cfg.DbURL, state.cfg.CurrentUsername)
	return nil
}

func handlerRegister(state *State, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("invalid login format. usage: %s <name>", cmd.name)
  }

  username := cmd.args[0]
  dbUser := database.CreateUserParams{
    ID: uuid.New(),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
    Name: username,
  }

  user, err := state.db.CreateUser(context.Background(), dbUser)
  if err != nil {
    return err
  }

  state.cfg.SetUser(username)
  fmt.Printf("User was created: %v", user)
  return nil 
}
