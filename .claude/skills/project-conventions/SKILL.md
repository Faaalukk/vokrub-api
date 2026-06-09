---
name: project-conventions
description: Canonical conventions for vokrub-api — Go + Fiber + GORM. Load the relevant reference file before writing any handler, model, middleware, or route code.
metadata:
  type: project
---

# vokrub-api — Project Conventions Index

Stack: **Go 1.26 · Fiber v2 · GORM v2 · PostgreSQL · JWT (golang-jwt/jwt v5)**

## Reference files

Load the file(s) relevant to your task from `.claude/skills/project-conventions/reference/`:

| Topic | File | When to load |
|---|---|---|
| Package & file layout | `directory-structure.md` | Adding any new file or package |
| GORM model structs | `model-conventions.md` | Adding or editing a model |
| HTTP handler functions | `handler-patterns.md` | Writing any handler |
| Route registration | `routing.md` | Adding routes |
| JWT middleware | `middleware-patterns.md` | Auth, guards |
| Response shape & errors | `error-handling.md` | Any response code |
| Architecture overview | `architecture.md` | Orientation / onboarding |

## Quick rules (always active)

- One file per domain in `handlers/`, `models/` — no giant files.
- All routes register through `routes.Setup(app)` — never in `main.go`.
- `database.DB` is the single global GORM handle — import from `database` package.
- All non-public routes must pass `middleware.Protected` as first handler arg.
- Input structs for body parsing: `type <Action>Input struct { ... }` — never parse directly into a model.
- Return `fiber.Map{"error": "..."}` for errors, `fiber.Map{"message": "..."}` for success messages.
- Never return raw GORM errors to the client — map to HTTP status + message.
