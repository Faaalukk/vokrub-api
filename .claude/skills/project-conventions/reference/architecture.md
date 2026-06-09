## Architecture

### Startup sequence

`main.go` → `godotenv.Load()` → `database.Connect()` → `fiber.New()` → global middleware → `routes.Setup(app)` → `app.Listen(":3000")`

`database.Connect()` opens the GORM connection and runs `DB.AutoMigrate(...)` for all registered models. Add new models to the AutoMigrate call here.

### Request lifecycle

```
HTTP request
  → Fiber logger middleware
  → CORS middleware
  → routes.Setup routing tree
    → middleware.Protected (if protected route)
      → handler function
        → database.DB query
        → c.JSON(response)
```

### Global middleware

Registered in `main.go` before `routes.Setup`:
- `logger.New()` — request logging
- `cors.New(...)` — CORS; `AllowOrigins` lists permitted client origins

### Environment

All config via `.env` + `godotenv`. Keys: `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `JWT_SECRET`. Access with `os.Getenv("KEY")`.

### Module path

`github.com/Faaalukk/vokrub-api.git` — use this as the import prefix for all internal packages.
