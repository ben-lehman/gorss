package main

import (
	"log"
	"os"

	"github.com/ben-lehman/gorss/internal/config"
)

type State struct {
	cfg *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
    log.Fatalf("Error reading config file: %v", err)
	}
	state := State{
		cfg: &conf,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*State, command) error),
	}

  commands.register("login", handlerLogin)

  args := os.Args
  if len(args) < 2 {
    log.Fatal("Usage: cli <command> [args...]")
    return
  }

  command := command{
    name: args[1],
    args: args[2:],
  }
  
  err = commands.run(&state, command)
  if err != nil {
    log.Fatal(err)
  }
	return
}
