package main

import (
	"context"
	"fmt"

	"github.com/ben-lehman/gorss/internal/rssFeed"
)

func handlerAggregate(state *State, cmd command) error {
  feed, err := rssFeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
  if err != nil {
    return err
  }

  fmt.Println(feed)
  return nil
}
