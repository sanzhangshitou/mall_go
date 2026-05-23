# Repository Guidelines

## Project Structure & Module Organization

This project follows a Laravel-inspired directory layout for a Go + Gin + GORM + Redis mall API.

```
mall/                          # Go module: mall
├── main.go                    # Entry point with graceful shutdown
├── bootstrap/app.go           # Boot sequence: .env → logger → DB → Redis
├── routes/api.go              # Route groups under /api/v1
├── app/
│   ├── Http/
│   │   ├── Controllers/       # H5 API controllers
│   │   ├── AdminControllers/  # Admin API controllers
│   │   ├── Middleware/        # CORS, recovery, request logging, admin auth
│   │   ├── Requests/          # H5 request DTOs
│   │   ├── AdminRequests/     # Admin request DTOs
│   │   ├── Resources/         # H5 response DTOs
│   │   └── AdminResources/    # Admin response DTOs
│   ├── Models/                # GORM models + enum constants
│   ├── Services/              # Business logic (shared by H5 + Admin)
│   └── Support/               # Infrastructure packages
│       ├── config/            #   .env loader
│       ├── database/          #   GORM MySQL connection
│       ├── cache/             #   Redis client + helpers
│       ├── logger/            #   Zap logger with Lumberjack rotation
│       ├── response/          #   Unified JSON response builder
│       └── paginator/         #   Page/page_size parsing
├── database/schema.sql        # Reference schema
└── storage/logs/              # Log output
```

## Build, Test, and Development Commands

| Command | Description |
|---|---|
| `make dev` / `go run main.go` | Start development server |
| `make build` | Compile binary to `./mall` |
| `make tidy` | Run `go mod tidy` to prune dependencies |
| `make clean` | Remove binary and log files |

Set `APP_DEBUG=true` in `.env` to enable GORM SQL logging.

## Coding Style & Naming Conventions

- **Go standard formatting** — run `gofmt` or `goimports` before committing.
- Package names are lowercase, single-word (e.g., `services`, not `my_services`).
- Controller methods mirror REST actions: `Index`, `Show`, `Store`, `Update`, `Destroy`.
- Request DTOs use pointer fields (`*int64`, `*bool`) in Update structs to distinguish "not passed" from "zero value".
- `sql.NullString` / `sql.NullFloat64` for nullable DB columns.

### Enum Pattern

Define typed constants in `app/Models/`:

```go
type ProductStatus int
const (
    ProductStatusDraft  ProductStatus = 0
    ProductStatusOnSale ProductStatus = 1
)
```

Always use enum constants in queries, never magic numbers.

## Architecture Guidelines

### Query Strategy (Repository-less)

- **List endpoints**: query the main table with pagination first, then batch-fetch related display data (e.g., category names) using `GetByIDs`. Never join tables just for display fields.
- **Detail endpoints**: fetch the main record first, return 404 if absent, then fetch SKUs, images, and related names in separate queries.
- Only join tables when filtering, sorting, or aggregation genuinely depends on the joined table.

### Layering

```
Controller → Request (validation) → Service (logic) → Resource (formatting) → Response
```

- Controllers parse inputs, call services, format output via Resources.
- Services contain reusable business logic and never reference `gin.Context`.
- Resources transform models into API response DTOs, including enum labels and associated data.

## Testing Guidelines

- Place test files alongside the code they test: `*_test.go`.
- Use the standard `testing` package.
- Integration tests that require DB/Redis should read connection info from `.env` and skip with `t.Skip()` when unavailable.

## Commit & Pull Request Guidelines

- Use conventional commit prefixes: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`.
- Keep commits atomic — one logical change per commit.
- PR descriptions should summarize what changed and why, with links to related issues.
- For new API endpoints, include a sample request/response in the PR body.
