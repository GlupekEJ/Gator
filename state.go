package main

import (
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
