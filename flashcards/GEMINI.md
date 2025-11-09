# GEMINI.md

This file provides guidance to Gemini CLI (https://codeassist.google/) when working with code in this repository.

## Common Development Commands

**Build and Run:**
- `go build ./cmd/main.go` - Build the application (after making changes, always run this to verify compilation)
- `make build` - Build with proper binary name (creates `todo-api`)
- `make run` - Run the application directly
- `make clean` - Clean build artifacts

**Database Operations:**
- `make db-start` - Start Supabase local PostgreSQL environment
- `make db-up` - Run database migrations
- `make db-stop` - Stop Supabase environment
- `make db-reset` - Complete database reset (stop, start, migrate)

## Architecture Overview

This is a **dual-purpose REST API** serving both **todos and flashcards** with identical architectural patterns.

**Clean Architecture Pattern:**
- `cmd/main.go` - Application entry point with dependency injection
- `models/` - Data structures and DTOs (todo.go, flashcard.go)
- `handlers/` - HTTP request/response handling (todoHandler.go, flashcardHandler.go)
- `services/` - Business logic and validation (todoService.go, flashcardService.go)
- `db/` - Repository pattern with PostgreSQL implementation (todoDb.go, flashcardDb.go)
- `config/` - Environment-based configuration management

**Key Patterns:**
- Both todos and flashcards follow **identical architectural patterns**
- Repository interfaces for database abstraction
- Service layer handles validation and business logic
- Handlers manage HTTP concerns (JSON, status codes, error responses)
- Gorilla Mux for routing with pattern-based routes (`/todos/{id:[0-9]+}`, `/flashcards/{id:[0-9]+}`)

**Database:**
- PostgreSQL with schema `gocourse.todos` and `gocourse.flashcards`
- Migrations in `supabase/migrations/` with timestamp prefixes
- Local development via Supabase CLI (accessible at localhost:54322)

**Data Models:**
- **Todos**: ID, Title, Description, Completed, CreatedAt, UpdatedAt
- **Flashcards**: ID, Content, CreatedAt, UpdatedAt (minimal content-only design)

**When adding new entities**, follow the exact same layered pattern: model → migration → repository → service → handler → route registration in main.go.

## Environment Setup

Requires `.env` file with:
- `DB_URL` - PostgreSQL connection string
- `PORT` - Application port (defaults to 8080)

Both Docker and Supabase CLI must be installed and running for database operations.

## Development Workflow Guidance



- After each task is complete, go build the project to verify your work and then fix any issues which come up.



## Structured Logging



This project uses the `log/slog` package for structured, leveled logging.



**Key Principles:**

- **JSON Output:** Logs are formatted as JSON to `os.Stdout` for machine-readability.

- **Dependency Injection:** The `slog.Logger` is created in `main.go` and injected into handlers, services, and repositories.

- **Log Levels:**

    - `INFO`: Used for happy-path events (e.g., start and success of an operation).

    - `ERROR`: Used for errors.

- **Error Logging Strategy:** To avoid redundant logging, errors are logged at two main points:

    1.  **Repository Layer:** When a database operation fails, the error is logged with `slog.Error`.

    2.  **Handler Layer:** When a service returns an error, the handler logs it with `slog.Error` before sending an HTTP error response.

    - Services do **not** log errors from the repository; they simply return them up the call stack.



**When adding new features, follow this pattern:**

1.  Ensure your new handlers, services, and repositories accept a `*slog.Logger` in their constructor.

2.  Add `INFO` level logs at the beginning and successful completion of each public method.

3.  Log errors from external calls (like databases or other APIs) at the point of failure.

4.  In handlers, log any error received from the service layer before generating the HTTP response.
