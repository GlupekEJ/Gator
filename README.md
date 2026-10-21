# Gator
CLI Blog aggregator in Go
Gator is a command-line RSS feed aggregator written in Go. It allows users to register, follow RSS feeds, collect posts, and browse recent posts from the feeds they follow.

Prerequisites

Before running Gator, make sure you have the following installed:

Go — required to install and run the CLI. Install it from the official Go website.
PostgreSQL — required to store users, feeds, feed follows, and posts. Install it from the official PostgreSQL website.

You will also need a running PostgreSQL server and a database connection URL.

Installation

Install Gator using go install:

go install <your-module-path>@latest

Replace <your-module-path> with the Go module path of this project, for example, github.com/yourusername/gator.

Make sure your Go binary directory is included in your system's PATH. You can check your installation with:

gator

If you are working on a local copy of the project and want to run it without installing it, you can use go run . from the project directory instead.

Configuration

Gator reads its configuration from a .gatorconfig.json file in your home directory.

Create the file and configure your PostgreSQL connection URL:

{
  "db_url": "postgres://postgres:your_password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

Replace your_password with your PostgreSQL password and gator with your database name if you chose a different one.

Create the PostgreSQL database before running Gator. For example, using psql:

CREATE DATABASE gator;

The current_user_name field is managed by Gator when you register or log in.

Database migrations

Before using Gator, apply all database migrations in the project's migrations directory to your PostgreSQL database. The migrations create the tables required for users, feeds, feed follows, and posts.

Use the migration tool and command configured by this project. If you are developing locally, consult the project's migration setup for the exact command and database URL format.

After changing SQL queries, regenerate the database code when required:

sqlc generate
Running Gator

Once PostgreSQL is running, the configuration is set up, and migrations have been applied, you can run Gator from your terminal.

If installed:

gator <command> [arguments]

From the project directory:

go run . <command> [arguments]
Available commands
Command	Description	Example
register	Register a new user	gator register alice
login	Log in as an existing user	gator login alice
users	List registered users	gator users
addfeed	Add a feed and follow it	gator addfeed "Example Feed" https://example.com/feed.xml
feeds	List registered feeds	gator feeds
follow	Follow a feed by URL	gator follow https://example.com/feed.xml
following	List feeds followed by the current user	gator following
unfollow	Unfollow a feed by URL	gator unfollow https://example.com/feed.xml
agg	Continuously fetch feeds and save their posts	gator agg 1m
browse	Browse recent posts from followed feeds	gator browse
browse 10	Browse up to 10 posts	gator browse 10

Some commands require a logged-in user. Use login or register before running them.

Collecting RSS posts

The agg command runs continuously, fetching one feed at a time at the specified interval:

gator agg 1m

The duration can be specified in Go duration format, such as 30s, 1m, or 1h. Leave the command running in one terminal while you use Gator in another terminal.

Press Ctrl+C to stop the aggregator.

Browsing posts

Once the aggregator has collected posts, browse the latest posts from feeds you follow:

gator browse

By default, browse displays up to two posts. Specify a positive integer to change the limit:

gator browse 10
Typical workflow
Create a PostgreSQL database and configure .gatorconfig.json.
Apply the database migrations.

Register a user:

gator register alice

Add an RSS feed:

gator addfeed "Example Feed" https://example.com/feed.xml

Start the aggregator in one terminal:

gator agg 1m

In another terminal, browse collected posts:

gator browse 10

Note: Replace the example feed URL with a valid RSS or Atom feed URL. The aggregator can only collect posts from feeds that have been added to the database.
