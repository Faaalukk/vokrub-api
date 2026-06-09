## Handler Patterns

### Signature

All handlers have signature `func Name(c *fiber.Ctx) error`.

### Standard CRUD template

```go
// GET /api/word
func GetWords(c *fiber.Ctx) error {
    var words []models.Word
    database.DB.Find(&words)
    return c.JSON(words)
}

// GET /api/word/:id
func GetWord(c *fiber.Ctx) error {
    id := c.Params("id")
    var word models.Word
    if result := database.DB.First(&word, id); result.Error != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Word not found"})
    }
    return c.JSON(word)
}

// POST /api/word
func CreateWord(c *fiber.Ctx) error {
    input := new(CreateWordInput)
    if err := c.BodyParser(input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }
    word := models.Word{Word: input.Word, Meaning: input.Meaning}
    if result := database.DB.Create(&word); result.Error != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Could not create word"})
    }
    return c.Status(201).JSON(word)
}

// DELETE /api/word/:id
func DeleteWord(c *fiber.Ctx) error {
    id := c.Params("id")
    database.DB.Delete(&models.Word{}, id)
    return c.JSON(fiber.Map{"message": "Deleted"})
}
```

### Getting auth context

```go
userID := c.Locals("user_id")   // set by middleware.Protected
role   := c.Locals("role")
```

### Query params

```go
page := c.QueryInt("page", 1)
limit := c.QueryInt("limit", 20)
```

### Never

- Return raw `result.Error.Error()` to the client.
- Put business logic in `routes.go`.
- Parse body into a model struct directly.
