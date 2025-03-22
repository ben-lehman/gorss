package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ben-lehman/gorss/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(state *State, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("Invalid number of args. Usage: follow <url>")
  }

  user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUsername)
  if err != nil {
    return err
  }

  feed, err := state.db.GetFeedByUrl(context.Background(), cmd.args[0])
  if err != nil {
    return err
  }

  feed_follow, err := state.db.CreateFeedFollow(
    context.Background(),
    database.CreateFeedFollowParams{
      ID: uuid.New(),
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
      UserID: user.ID,
      FeedID: feed.ID,
  })
  if err != nil {
    return err
  }
  fmt.Println("Feed name: ", feed_follow.FeedName)
  fmt.Println("User name: ", feed_follow.UserName)

  return nil
}

func handlerFollowing(state *State, cmd command) error {
  username := state.cfg.CurrentUsername
  feed_follows, err := state.db.GetFeedFollowsForUser(context.Background(), username)
  if err != nil {
    return err
  }

  for _, feed_follow := range feed_follows {
    fmt.Println(feed_follow.FeedName)
  }

  return nil
}
