package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.args[0], err)
	}

	if timeBetweenRequests <= 0 {
		return fmt.Errorf("duration must be greater than zero")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	// Get the feed that should be fetched next.
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("could not get next feed: %w", err)
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)

	// Mark the feed as fetched.
	err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	// Fetch and parse the RSS feed using the existing function.
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed %s: %w", feed.Name, err)
	}

	// Print the title of every item in the feed.
	for _, item := range rssFeed.Channel.Item {
		now := time.Now()

		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      feed.ID,
		}

		if err := s.db.CreatePost(ctx, post); err != nil {
			log.Printf("could not save post %q: %v", item.Title, err)
			continue
		}
	}

	return nil
}

func parsePublishedAt(value string) sql.NullTime {
	value = strings.TrimSpace(value)
	if value == "" {
		return sql.NullTime{}
	}

	layouts := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
	}

	for _, layout := range layouts {
		publishedAt, err := time.Parse(layout, value)
		if err == nil {
			return sql.NullTime{
				Time:  publishedAt,
				Valid: true,
			}
		}
	}

	// An unrecognized date is treated as missing, not as a fatal error.
	return sql.NullTime{}
}
