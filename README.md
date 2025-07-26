# Gator - Your Personal CLI RSS Feed Aggregator

Gator is a command-line interface (CLI) application built with Go and PostgreSQL that allows you to aggregate, manage, and browse RSS (Really Simple Syndication) feeds directly from your terminal. It's designed for developers who prefer a minimalist, text-based approach to staying updated with their favorite blogs, news sites, and other content sources.

## Features

* **User Management**: Register and switch between multiple users.
* **Feed Management**: Add, view, follow, and unfollow RSS feeds.
* **Automated Aggregation**: A long-running process to periodically fetch and store posts from your followed feeds.
* **Post Browse**: View recent posts from all feeds you follow, directly in your terminal.

## Prerequisites

Before you can run `gator`, you need to have the following installed on your system:

* **Go**: Version 1.20 or newer. You can download it from [golang.org](https://golang.org/dl/).
* **PostgreSQL**: A running PostgreSQL database instance. You can download it from [postgresql.org](https://www.postgresql.org/download/).
* **Goose**: A database migration tool for Go. Install it with:
    ```bash
    go install [github.com/pressly/goose/v3/cmd/goose@latest](https://github.com/pressly/goose/v3/cmd/goose@latest)
    ```
* **SQLC**: A SQL compiler for Go. Install it with:
    ```bash
    go install [github.com/sqlc-dev/sqlc/cmd/sqlc@latest](https://github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
    ```

## Installation

You can install the `gator` CLI tool directly from your project's root directory:

1.  **Clone the repository**:
    ```bash
    git clone <YOUR_GITHUB_REPO_LINK_HERE>
    cd gator
    ```
    *(Remember to replace `<YOUR_GITHUB_REPO_LINK_HERE>` with the actual link after you push it to GitHub.)*

2.  **Install the `gator` binary**:
    From the root of the `gator` directory, run:
    ```bash
    go install .
    ```
    This command compiles your Go program into a standalone executable named `gator` (or `gator.exe` on Windows) and places it in your `GOPATH/bin` directory, which should be in your system's `PATH`. You can then run `gator` from any directory.

    **Note**: `go run .` is primarily for development and quick testing. For production use and to run `gator` as a standalone command, always use the compiled binary installed via `go install .`.

## Configuration

`gator` requires a `config.json` file in your user's home directory to store the PostgreSQL database connection string and the currently logged-in user.

1.  **Create the config directory**:
    ```bash
    mkdir -p ~/.gator
    ```

2.  **Create the `config.json` file**:
    Create a file named `config.json` inside the `~/.gator` directory with the following content:

    **`~/.gator/config.json`**
    ```json
    {
        "database_url": "postgres://postgres:postgres@localhost:5432/gator",
        "current_user_name": ""
    }
    ```
    * **`database_url`**: Replace `"postgres://postgres:postgres@localhost:5432/gator"` with your PostgreSQL connection string. Ensure the `gator` database exists.
    * **`current_user_name`**: This field will be automatically updated by `gator` when you use the `login` or `register` commands.

## Database Setup

Before running `gator` for the first time, you need to set up your PostgreSQL database and run the migrations.

1.  **Ensure your `config.json` is set up correctly** with your PostgreSQL `database_url`.
2.  **Run migrations**: From the `gator` project's root directory, navigate to the `sql/schema` folder and run `goose`:
    ```bash
    cd sql/schema
    goose postgres "$(cat ~/.gator/config.json | jq -r .database_url)" up
    # If jq is not installed, you can manually paste the URL:
    # goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
    cd ../..
    ```
    This will create all necessary tables in your database.

## Usage

Once `gator` is installed and configured, you can run commands from your terminal.

### Basic Commands:

* **`gator register <username>`**: Creates a new user in the database and automatically logs them in.
    ```bash
    gator register alice
    ```

* **`gator login <username>`**: Logs in an existing user, setting them as the current active user.
    ```bash
    gator login alice
    ```

* **`gator users`**: Lists all registered users and indicates the current logged-in user.
    ```bash
    gator users
    ```

* **`gator reset`**: **CAUTION!** This command deletes all users and their associated data (feeds, follows, posts) from the database. Use with care.
    ```bash
    gator reset
    ```

### Feed Management:

* **`gator addfeed "<feed name>" "<feed url>"`**: Adds a new RSS feed to the system and automatically follows it for the current user.
    ```bash
    gator addfeed "Hacker News" "[https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)"
    gator addfeed "TechCrunch" "[https://techcrunch.com/feed/](https://techcrunch.com/feed/)"
    ```

* **`gator feeds`**: Lists all feeds added to the system (regardless of who follows them).
    ```bash
    gator feeds
    ```

* **`gator follow "<feed url>"`**: Allows the current user to follow an existing feed by its URL.
    ```bash
    gator follow "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"
    ```

* **`gator unfollow "<feed url>"`**: Allows the current user to unfollow a feed by its URL.
    ```bash
    gator unfollow "[https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)"
    ```

* **`gator following`**: Lists all feeds currently followed by the logged-in user.
    ```bash
    gator following
    ```

### Aggregation and Browse:

* **`gator agg <time_between_requests>`**: Starts a long-running process to fetch posts from followed feeds. `time_between_requests` is a duration string (e.g., `1s`, `30s`, `1m`, `5m`).
    ```bash
    gator agg 30s
    # This command will run indefinitely until you press Ctrl+C
    ```
    It's recommended to run this command in a separate terminal or in the background.

* **`gator browse [limit]`**: Displays recent posts from feeds followed by the current user. Optionally specify a `limit` (defaults to 2 if not provided).
    ```bash
    gator browse
    gator browse 10
    ```

---

**Example Workflow:**

```bash
# Register a new user
gator register devuser

# Add some feeds
gator addfeed "Boot.dev Blog" "[https://blog.boot.dev/index.xml](https://blog.boot.dev/index.xml)"
gator addfeed "Hacker News" "[https://news.ycombinator.com/rss](https://news.ycombinator.com/rss)"

# Verify what you're following
gator following

# Start the aggregator (in a separate terminal)
# This will fetch and save posts every 10 seconds
gator agg 10s

# Back in your main terminal, browse the posts
gator browse
gator browse 5 # See more posts
