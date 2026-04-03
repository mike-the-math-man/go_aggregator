# Gator - RSS Feed Aggregator CLI

Gator is a command-line RSS feed aggregator built in Go. It allows users to register, follow RSS feeds, and browse the latest posts from their followed feeds. The application uses PostgreSQL as its database backend and provides a simple CLI interface for managing feeds and users.

## Prerequisites

Before running Gator, ensure you have the following installed:

- **Go**: Version 1.19 or later. Download from [golang.org](https://golang.org/dl/).
- **PostgreSQL**: A running PostgreSQL database server. You can install it via your package manager (e.g., `apt install postgresql` on Ubuntu, `brew install postgresql` on macOS) or download from [postgresql.org](https://www.postgresql.org/download/).

## Installation

Install the Gator CLI using Go's package manager:

```bash
go install github.com/mike-the-math-man/go_aggregator@latest
```

This will download, compile, and install the `gator` binary to your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH.

**Note**: Go programs are statically compiled binaries. After running `go install`, you can run the `gator` command directly without needing the Go toolchain. Use `go run .` only during development.

## Setup

### 1. Database Setup

Create a PostgreSQL database for Gator. You can do this using the `psql` command-line tool or your preferred database management tool:

```bash
createdb gator
```

Run the database migrations to set up the schema. The project uses Goose for migrations. From the project root:

```bash
goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

Replace `username`, `password`, and connection details with your PostgreSQL credentials.

### 2. Configuration

Gator stores its configuration in a JSON file at `~/.gatorconfig.json`. Create this file with your database connection URL:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url`: Your PostgreSQL connection string.
- `current_user_name`: This will be set automatically when you log in.

## Usage

### Running the Program

After installation and setup, you can run Gator commands directly:

```bash
gator <command> [arguments]
```

### Available Commands

Here are some key commands you can use:

- **Register a new user**:
  ```bash
  gator register <username>
  ```

- **Login as an existing user**:
  ```bash
  gator login <username>
  ```

- **List all users**:
  ```bash
  gator users
  ```

- **Add a new RSS feed** (requires login):
  ```bash
  gator addfeed <name> <url>
  ```

- **List all feeds**:
  ```bash
  gator feeds
  ```

- **Follow a feed** (requires login):
  ```bash
  gator follow <url>
  ```

- **List followed feeds** (requires login):
  ```bash
  gator following
  ```

- **Unfollow a feed** (requires login):
  ```bash
  gator unfollow <url>
  ```

- **Browse recent posts** (requires login):
  ```bash
  gator browse [limit]
  ```
  The optional `limit` parameter specifies how many posts to display (default: 2).

- **Start the aggregator** (fetches new posts from feeds):
  ```bash
  gator agg <time_between_requests>
  ```
  Example: `gator agg 30s` to fetch every 30 seconds.

- **Reset the database** (deletes all users):
  ```bash
  gator reset
  ```

### Development

For development, you can run the program without installing it:

```bash
go run .
```

To build a local binary:

```bash
go build -o gator .
```

Then run `./gator <command>`.

## Project Structure

- `main.go`: Entry point and command registration.
- `commands.go`: Command handlers and business logic.
- `rss.go`: RSS feed parsing structures.
- `internal/config/`: Configuration management.
- `internal/database/`: Generated database code (using sqlc).
- `sql/schema/`: Database schema migrations.
- `sql/queries/`: SQL queries for database operations.

## Dependencies

- `github.com/google/uuid`: For generating unique IDs.
- `github.com/lib/pq`: PostgreSQL driver for Go.
- `sqlc`: For generating type-safe Go code from SQL queries.
- `goose`: For database migrations.

