## Error Handling & Response Shape

### Success responses

```go
// Single resource
return c.JSON(resource)

// Created
return c.Status(201).JSON(resource)

// Action success (no body)
return c.JSON(fiber.Map{"message": "Deleted"})

// Custom shape
return c.JSON(fiber.Map{"token": tokenString, "user": fiber.Map{...}})
```

### Error responses

Always use `fiber.Map{"error": "human readable message"}`. Never expose raw Go errors or GORM error strings.

```go
// 400 — bad request / validation
return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})

// 401 — unauthenticated
return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})

// 403 — forbidden
return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})

// 404 — not found
return c.Status(404).JSON(fiber.Map{"error": "Word not found"})

// 500 — internal
return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
```

### GORM error check pattern

```go
if result := database.DB.First(&word, id); result.Error != nil {
    return c.Status(404).JSON(fiber.Map{"error": "Word not found"})
}
```

### List endpoints

Return the array directly (not wrapped) unless pagination metadata is needed:

```go
// simple
return c.JSON(words)

// with pagination (future)
return c.JSON(fiber.Map{"data": words, "total": total, "page": page})
```
