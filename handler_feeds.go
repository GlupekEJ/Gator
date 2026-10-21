package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: gator feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf(
			"* %s\n  URL: %s\n  Added by: %s\n\n",
			feed.FeedName,
			feed.Url,
			feed.UserName,
		)
	}

	return nil
}
