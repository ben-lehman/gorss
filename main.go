package main

import (
	"context"
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
	commands.register("agg", handlerAggregate)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

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

func middlewareLoggedIn(handler func(state *State, cmd command, user database.User) error) func(*State, command) error {
	return func(state *State, cmd command) error {
		user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUsername)
		if err != nil {
			return err
		}

		return handler(state, cmd, user)
	}
}
