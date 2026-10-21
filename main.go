package main

import (
	"database/sql"
	"log"
	"os"

	"gator/internal/config"
	"gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Read the existing configuration.
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the shared application state.
	s := &state{
		cfg: &cfg,
		db:  database.New(db),
	}

	// Initialize the command registry.
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register the login command.
	cmds.register("login", handlerLogin)

	// Register the register command.
	cmds.register("register", handlerRegister)

	cmds.register("reset", handlerDeleteAllUsers)

	cmds.register("users", handlerListUsers)

	cmds.register("agg", handlerAgg)

	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))

	cmds.register("feeds", handlerFeeds)

	cmds.register("follow", middlewareLoggedIn(handlerFollow))

	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Read the command-line arguments.
	if len(os.Args) < 2 {
		log.Fatal("a command is required")
	}

	// Split the command name from its arguments.
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	// Run the selected command.
	if err := cmds.run(s, cmd); err != nil {
		log.Fatal(err)
	}
}
