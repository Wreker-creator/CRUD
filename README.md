# Task Manager REST API

A RESTful CRUD API built in Go using only the standard library — no frameworks. Built following a test-driven development approach with both unit and integration tests.

---

## Features

- Full CRUD operations via proper HTTP methods (GET, POST, PUT, DELETE)
- File system persistence using JSON
- In-memory store available for testing
- Interface-driven storage layer — swap between file system and in-memory store seamlessly
- Safe file writes using a custom `Tape` wrapper that truncates before rewriting

---

## Tech Stack

- **Language:** Go
- **Libraries:** Standard library only (`net/http`, `encoding/json`, `os`, `io`)
- **Testing:** `testing`, `net/http/httptest`

---

## Getting Started

### Run the server

```bash
go run cmd/webserver/main.go
```

Server starts on `http://localhost:5001`. A `task.db.json` file will be created in the project root to persist data.

### Run tests

```bash
go test ./...
```

---

## API Endpoints

| Method   | Endpoint      | Description          | Status Codes          |
|----------|---------------|----------------------|-----------------------|
| GET      | /tasks        | Get all tasks        | 200                   |
| GET      | /tasks/{id}   | Get task by ID       | 200, 404              |
| POST     | /tasks        | Create a new task    | 202, 400              |
| PUT      | /tasks/{id}   | Update a task        | 200, 400, 404         |
| DELETE   | /tasks/{id}   | Delete a task        | 200, 400, 404         |

---

## Example Requests

```bash
# Get all tasks
curl http://localhost:5001/tasks

# Get a specific task
curl http://localhost:5001/tasks/1

# Create a task
curl -X POST http://localhost:5001/tasks \
  -H "Content-Type: application/json" \
  -d '{"id":1,"title":"Buy groceries","description":"Milk, eggs, bread"}'

# Update a task
curl -X PUT http://localhost:5001/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Updated list"}'

# Delete a task
curl -X DELETE http://localhost:5001/tasks/1
```

---

## Testing Strategy

### Unit Tests
- **Server tests** (`fetch_test.go`) — test each HTTP handler in isolation using a `StubTaskStore`, verifying correct status codes and response bodies for all CRUD operations
- **Store tests** (`file_system_store_test.go`) — test `FileSystemTaskStore` directly, verifying reads, writes, updates, and deletes against a real temp file

### Integration Tests
- Wire a real `FileSystemTaskStore` with a real `TaskServer` and exercise the full request-response cycle end to end

---

## Design Decisions

**Interface-driven storage** — `TaskServer` depends on a `TaskStore` interface, not a concrete implementation. This made the `StubTaskStore` in tests trivial to write and will make swapping to a PostgreSQL store in future straightforward.

**TDD approach** — every feature was written test-first. Tests drove the API design, particularly around error handling and status codes.

**Tape wrapper** — a custom `Tape` type wraps `*os.File` and overrides `Write` to always seek to the start and truncate before writing. This prevents stale bytes from corrupting the JSON file after delete operations shorten the content.

**In-memory caching** — `FileSystemTaskStore` reads the file once at startup and caches the task list in memory. All subsequent reads are served from the cache and writes update both the cache and the file, avoiding repeated disk reads.
