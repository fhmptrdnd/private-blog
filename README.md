# Telegraph Clone (refactored)

This repository is a small example web application (Telegraph-like) that was refactored into an idiomatic Go project layout and updated to use the Repository pattern.

The refactor separates concerns into packages and makes it easier to swap data storage implementations, add tests, and extend the app.

---

## Project layout

```
cmd/
  web/                # application entrypoint (main)
internal/
  handler/            # HTTP handlers + templates
    handler.go
  models/             # domain models
    article.go
  repository/         # repository interface + implementations
    repository.go
    memory.go          # in-memory implementation
  service/            # business logic
    article_service.go
go.mod
README.md
```

Why this layout?
- `cmd/web` contains the program entrypoint. Additional commands can be added later.
- `internal/*` packages are private to this module and keep implementation details inside the repository.
- `repository` implements the Repository pattern: code depends on an interface, and concrete storage implementations (in-memory, DB-backed) implement it.

---

## What changed during refactor

- The original single-file `main.go` app has been split into packages: models, repository (interface + memory), service and handler.
- `service.ArticleService` contains business logic (ID generation, sanitization, create/get/update/delete operations).
- `repository.Repository` is an interface; `repository.MemoryRepo` is a simple thread-safe in-memory implementation.
- `handler` is responsible for HTTP endpoints and templates.
- Entrypoint `cmd/web/main.go` wires repository -> service -> handler and starts the HTTP server on :8080.

---

## How to run

Make sure you have Go installed (Go 1.20+ recommended). From the project root run:

Windows (cmd.exe) - run directly:

```cmd
go run ./cmd/web
```

or build and run executable:

```cmd
go build -o web ./cmd/web
.\web
```

The server will run at `http://localhost:8080`.

---

## Endpoints

- GET `/` — home editor page
- POST `/create` — create article
- GET `/view/{id}` — view article
- GET `/edit/{id}` — edit page (ownership required)
- POST `/update/{id}` — update article (ownership required)
- POST `/delete/{id}` — delete article (ownership required)

Ownership is tracked by a cookie `user_id` created on first visit.

---

## Notes, limitations and recommendations

- The repository currently uses an in-memory store. All data is lost on process exit.
- Sanitization is simplistic (`\n` -> `<br>`). For a production app use a proper HTML sanitizer and ensure templates are escaped when rendering untrusted HTML.
- Templates are embedded as constants in `internal/handler/handler.go` for simplicity. For maintainability, consider moving them to `templates/` files and loading via `template.ParseGlob`.
- Add structured logging and graceful shutdown for production readiness.

Recommended next steps you can ask me to implement:
- Add `templates/` directory and load templates from files.
- Add unit tests for `service` and `repository` (memory implementation).
- Add a SQL-backed repository (SQLite/Postgres) and migrations.
- Improve sanitization using a library (e.g., bluemonday) or proper escaping.

---

If you want, I can now:
- Move the templates to disk and update the handler to parse them from `templates/`.
- Add a basic `README` (done) and a couple of unit tests.
- Implement a persistent repository (SQLite) and demonstrate migrations.

Tell me which follow-up you prefer and I'll implement it next.
# weather-api-test-web-based
Nyoba Web Cek Cuaca berdasarkan Kota
