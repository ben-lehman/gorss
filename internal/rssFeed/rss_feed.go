package rssFeed

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
  req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
  if err != nil {
    return &RSSFeed{}, err
  }

  req.Header.Set("User-Agent", "gator")

  client := &http.Client{
    Timeout: 10 * time.Second,
  }
  res, err := client.Do(req)
  if err != nil {
    return &RSSFeed{}, err
  }
  defer res.Body.Close()

  data, err := io.ReadAll(res.Body)
  if err != nil {
    return &RSSFeed{}, err
  }
  
  var feed RSSFeed
  err = xml.Unmarshal(data, &feed)
  if err != nil {
    return &feed, err
  }

  cleanRSSFeedText(&feed)

  return &feed, err
}

func cleanRSSFeedText(feed *RSSFeed) {
  feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
  feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

  for i, item := range feed.Channel.Item {
    item.Title = html.UnescapeString(item.Title)
    item.Description = html.UnescapeString(item.Description)
    feed.Channel.Item[i] = item
  }
}
