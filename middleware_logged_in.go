package main

import (
	"context"
	"fmt"

	"gator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Look up the currently logged-in user.
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not get current user: %w", err)
		}

		// Call the handler with the authenticated user's data.
		return handler(s, cmd, user)
	}
}
