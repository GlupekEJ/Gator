package main

import (
	"context"
	"fmt"
	"strconv"

	"gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}

	limit := 2

	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("limit must be a positive integer")
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Try following a feed and running agg.")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published: ")

		if post.PublishedAt.Valid {
			fmt.Println(post.PublishedAt.Time.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("Unknown")
		}

		fmt.Printf("Description: %s\n", post.Description)
		fmt.Println("---")
	}

	return nil
}
