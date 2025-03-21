package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ben-lehman/gorss/internal/database"
	"github.com/google/uuid"
)


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

func handlerReset(state *State, cmd command) error {
  err := state.db.DeleteUsers(context.Background())
  if err != nil {
    return fmt.Errorf("Issue reseting users table: %v", err)
  }

  fmt.Println("Users table successfully reset")
  return nil
}
