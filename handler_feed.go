package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ben-lehman/gorss/internal/database"
	"github.com/ben-lehman/gorss/internal/rssFeed"
	"github.com/google/uuid"
)

func handlerAddFeed(state *State, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("invalid add feed format. Usage: %s <name> <url>", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]
	user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUsername)
	if err != nil {
		return err
	}

	feed, err := state.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		})
	if err != nil {
		return err
	}

	_, err = state.db.CreateFeedFollow(
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

	printFeed(feed)

	return nil
}

func handlerFeeds(state *State, _ command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Unable to get feeds: %v", err)
	}

	for _, feed := range feeds {
		feedUser, err := state.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Unable to get user name tied to feed: %v", err)
		}
		fmt.Printf("name: %s\n", feed.Name)
		fmt.Printf("url: %s\n", feed.Url)
		fmt.Printf("user id: %s\n", feedUser.Name)
	}

	return nil
}

func handlerAggregate(state *State, cmd command) error {
	feed, err := rssFeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("id: %v\ncreated at: %v\n updated at:%v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt)
	fmt.Printf("name: %v\n url:%v\n user id:%v\n", feed.Name, feed.Url, feed.UserID)

}
