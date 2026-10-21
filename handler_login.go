package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login requires a username")
	}

	if len(cmd.args) > 1 {
		return errors.New("usage: gator login <username>")
	}

	username := cmd.args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("user %q does not exist", username)
	}
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Current user has been set to %s\n", username)

	return nil
}
