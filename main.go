package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ben-lehman/gorss/internal/config"
	"github.com/ben-lehman/gorss/internal/database"
	_ "github.com/lib/pq"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	dbConnectionURL := conf.DbURL
	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state := State{
		db:  dbQueries,
		cfg: &conf,
	}
	commands := commands{
		registeredCommands: make(map[string]func(*State, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
  commands.register("users", handlerUsers)
  commands.register("reset", handlerReset)

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
