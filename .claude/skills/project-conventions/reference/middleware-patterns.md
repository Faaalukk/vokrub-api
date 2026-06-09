## Middleware Patterns

### Protected (JWT guard)

`middleware.Protected` in `middleware/auth.go`:

1. Reads `Authorization: Bearer <token>` header
2. Parses and validates JWT with `JWT_SECRET`
3. Sets `c.Locals("user_id", claims["user_id"])` and `c.Locals("role", claims["role"])`
4. Returns `401` with `{"error": "..."}` if missing or invalid

### Adding new middleware

Create `middleware/<name>.go`, define `func Name(c *fiber.Ctx) error`. Register in `routes.Setup` or in `main.go` for global middleware.

### Role guard example

```go
func AdminOnly(c *fiber.Ctx) error {
    role := c.Locals("role")
    if role != "admin" {
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }
    return c.Next()
}
```

Chain after `Protected`: `route.Delete("/", middleware.Protected, middleware.AdminOnly, handlers.DeleteUser)`

### JWT token shape

Claims: `user_id` (uint), `email` (string), `role` (string), `exp` (Unix timestamp, 24h). Signed HS256 with `JWT_SECRET` env var.
