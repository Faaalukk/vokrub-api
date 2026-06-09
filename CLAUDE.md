## Read this first — Project skill

Before doing any non-trivial work in this repo, read the project skill index at `.claude/skills/project-conventions/SKILL.md` and then load the relevant topic file under `.claude/skills/project-conventions/reference/`. The skill is the canonical convention set for this codebase (one H2 topic per reference file): architecture, directory structure, handler patterns, model conventions, middleware, routing, error handling.

Pull from it rather than inventing a parallel pattern.

## Commands

```bash
go run main.go          # start dev server on :3000
go build -o server .    # production binary
go mod tidy             # sync dependencies
```

Environment loaded from `.env` via `godotenv` — copy `.env.example` if it exists, or set:

```
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=5432
JWT_SECRET=
```

Health check: `GET /health` → `{"status":"ok"}`

## Architecture

### Entry point

`main.go` bootstraps in order: load `.env` → `database.Connect()` (AutoMigrate runs here) → create Fiber app → attach global middleware (logger, CORS) → `routes.Setup(app)` → listen `:3000`.

CORS currently allows `http://localhost:8080`. Update `AllowOrigins` in `main.go` for new clients.

### Package layout

```
main.go            entry — wires everything
database/db.go     GORM connection + AutoMigrate
models/            GORM model structs (one file per domain)
handlers/          HTTP handlers (one file per domain)
middleware/        Fiber middleware (auth JWT guard)
routes/routes.go   all route registration via routes.Setup(app)
```

### Route registration

All routes live in `routes/routes.go` → `routes.Setup(app)`. Group by domain under `/api/<domain>`. Protected routes pass `middleware.Protected` as the first handler argument.

### Auth

`middleware.Protected` validates `Authorization: Bearer <jwt>` and sets `c.Locals("user_id")` and `c.Locals("role")` for downstream handlers. JWT signed with `JWT_SECRET` env var, 24-hour expiry, claims: `user_id`, `email`, `role`.

### Database

`database.DB` is a package-level `*gorm.DB`. Import `github.com/Faaalukk/vokrub-api.git/database` and use `database.DB` directly in handlers. AutoMigrate runs at startup for all registered models.

### Models

Embed `gorm.Model`-style fields manually (ID, CreatedAt, UpdatedAt, DeletedAt) for full control over JSON tags. Use `gorm:"uniqueIndex"` for unique fields. Soft-delete via `gorm.DeletedAt`.

## Naming

- Package names: lowercase, single word matching directory name
- Files: `snake_case.go`
- Types/structs: `PascalCase`
- Handler functions: verb + noun (`GetCustomers`, `CreateCustomer`, `DeleteCustomer`)
- Input structs: `<Action>Input` (e.g. `LoginInput`)

## When in doubt

Consult `.claude/skills/project-conventions/reference/`. Quick map:

- Adding a new domain → `handler-patterns.md` + `routing.md` + `model-conventions.md`
- Auth / JWT → `middleware-patterns.md`
- Response shape → `error-handling.md`
- Package / file layout → `directory-structure.md`
