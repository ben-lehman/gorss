package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ben-lehman/gorss/internal/rssFeed"
)

func handlerAggregate(state *State, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Invalid number of args. Usage agg <time-duration>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Invalid time duration format: %v", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		fmt.Println("Scraping...")
		err := scrapeFeeds(state)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(state *State) error {
	feedToFetch, err := state.db.GetNextFeedToFetch(context.Background())
	if err != nil {
	return fmt.Errorf("Issue fetching next feed: %v", err)
	}

	err = state.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return fmt.Errorf("Issue marking feed as fetched: %v", err)
	}

	feed, err := rssFeed.FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Issue fetching RSS feed: %v", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
