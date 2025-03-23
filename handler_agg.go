package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ben-lehman/gorss/internal/database"
	"github.com/ben-lehman/gorss/internal/rssFeed"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		fmt.Println("saving post: ", item.PubDate)
		publishedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return fmt.Errorf("Issue parsing published date: %v", err)
		}

		post, err := state.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description},
			PublishedAt: sql.NullTime{Time: publishedDate},
			FeedID:      feedToFetch.ID,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" { // 23505 is the error code for unique violation
					fmt.Println("Duplicate Url, skipping insertion")
					continue
				}
			}
			return fmt.Errorf("Issue creating post: %v", err)
		}

		fmt.Println("post created: ", post.Title)
	}

	return nil
}

func handlerBrowse(state *State, cmd command, user database.User) error {
	var postsLimit int32 = 2
	if len(cmd.args) == 1 {
		num, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid param. Usage: browse <int>")
		}
		postsLimit = int32(num)
	}

	posts, err := state.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			ID:    user.ID,
			Limit: postsLimit,
		})
	if err != nil {
		return fmt.Errorf("Issue fetching posts: %v", err)
	}

	for _, post := range posts {
		fmt.Println("title: ", post.Title)
	}

	return nil
}

func handlerResetPosts(state *State, _ command) error {
	err := state.db.DeletePosts(context.Background())
	if err != nil {
		return fmt.Errorf("Issue deleteing posts: %v", err)
	}

	return nil
}
