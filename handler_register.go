package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register requires a username")
	}

	if len(cmd.args) > 1 {
		return errors.New("usage: gator register <username>")
	}

	username := cmd.args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		return fmt.Errorf("user %q already exists", username)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not check whether user exists: %w", err)
	}

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("user was created, but could not update config: %w", err)
	}

	fmt.Printf("User %s created successfully!\n", username)
	// Log the user data for debugging.
	log.Printf("Created user: %+v", user)
	return nil
}

func handlerDeleteAllUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("delete-all-users requires no arguments")
	}

	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}

	fmt.Println("All users have been deleted successfully.")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("list-users requires no arguments")
	}

	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	fmt.Println("List of users:")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
