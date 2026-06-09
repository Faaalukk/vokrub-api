## Routing

All routes register in `routes/routes.go` → `Setup(app *fiber.App)`. Never add routes in `main.go`.

### Pattern

```go
func Setup(app *fiber.App) {
    api := app.Group("/api")

    word := api.Group("/word")
    word.Get("/",    middleware.Protected, handlers.GetWords)
    word.Get("/:id", middleware.Protected, handlers.GetWord)
    word.Post("/",   middleware.Protected, handlers.CreateWord)
    word.Put("/:id", middleware.Protected, handlers.UpdateWord)
    word.Delete("/:id", middleware.Protected, handlers.DeleteWord)
}
```

### Conventions

- Top-level prefix: `/api`
- Domain group: `/api/<domain>` (singular noun: `customer`, `word`, `phrase`)
- `middleware.Protected` is always the first handler argument on protected routes
- Public routes (auth) still use `middleware.Protected` where token is expected — if a route truly needs no auth, omit the middleware
- Route parameter: `/:id` for single-resource access

### URL shape

| Action | Method | Path |
|---|---|---|
| List all | GET | `/api/word` |
| Get one | GET | `/api/word/:id` |
| Create | POST | `/api/word` |
| Update | PUT | `/api/word/:id` |
| Delete | DELETE | `/api/word/:id` |
