# gorss

This repo is for personal education purposes to get some practice creating a command line utility in go and does not provide an elegant user experience.

## Requirements
[go](https://go.dev/doc/install)1.24
[postgres](https://www.postgresql.org/download/) 16.8

## Installation

```
go install github.com/ben-lehman/gorss
```

## Setup

Create `.gatorconfig.json` in your HOME directory with the following contents:

```
{
  "db_url": "<postgres connection string>",
  "current_user_name": ""
}
```

Set up a local postgres server and set `db_rul` in .gatorconfig.json to its connection string.

## Usage

#### Users
GoRSS supports the ability for multiple users to follow their own RSS feeds.
Here are commands useful for managing users:

```
# Registers a new user in the db
> gorss register <username>

# Lists registered users
> gorss users

# Login as a specific user
> gorss login <username>
```
#### Feeds
Use the following commands to manage RSS feeds:
```
# Adds an RSS feed to the db
> gorss addfeed <feed name> <feed url>

# Lists added RSS feeds
> gorss feeds

# Follows an RSS feed for the current user
> gorss follow <feed url>

# Unfollow an RSS feed for the current user
> gorss unfollow <feed url>

# Lists posts from followed feeds. Add a limit param to set amount of posts recieved
> gorss browse <optional limit>
```

#### Aggregate Posts
To add posts from the registered RSS feeds run the following command:
```
# fetches posts from RSS feeds at the duration interval
# duration examples: 5s, 10m, 1h
> gorss agg <duration>
```
