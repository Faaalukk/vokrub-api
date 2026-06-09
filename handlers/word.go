package handlers

import (
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type CreateWordInput struct {
	Word    string `json:"word"`
	Pos     string `json:"pos"`
	Meaning string `json:"meaning"`
	Note    string `json:"note"`
}

type ReviewWordInput struct {
	Correct bool `json:"correct"`
}

// GET /api/word
func GetWords(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	var words []models.Word
	database.DB.Where("customer_id = ?", customerID).Order("created_at desc").Find(&words)
	return c.JSON(words)
}

// GET /api/word/:id
func GetWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var word models.Word
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&word); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Word not found"})
	}
	return c.JSON(word)
}

// POST /api/word
func CreateWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	input := new(CreateWordInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	word := models.Word{
		CustomerID: customerID.(uint),
		Word:       input.Word,
		Pos:        input.Pos,
		Meaning:    input.Meaning,
		Note:       input.Note,
		Box:        1,
		Seen:       0,
		Due:        true,
		Added:      time.Now().Format("2006-01-02"),
	}

	if result := database.DB.Create(&word); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create word"})
	}
	return c.Status(201).JSON(word)
}

// PUT /api/word/:id
func UpdateWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var word models.Word
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&word); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Word not found"})
	}

	input := new(CreateWordInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	database.DB.Model(&word).Updates(map[string]interface{}{
		"word":    input.Word,
		"pos":     input.Pos,
		"meaning": input.Meaning,
		"note":    input.Note,
	})
	return c.JSON(word)
}

// DELETE /api/word/:id
func DeleteWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	database.DB.Where("customer_id = ?", customerID).Delete(&models.Word{}, id)
	return c.JSON(fiber.Map{"message": "Deleted"})
}

// POST /api/word/:id/review
func ReviewWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var word models.Word
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&word); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Word not found"})
	}

	input := new(ReviewWordInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	newBox := word.Box + 1
	if !input.Correct {
		newBox = word.Box - 1
	}
	if newBox > 5 {
		newBox = 5
	}
	if newBox < 1 {
		newBox = 1
	}

	database.DB.Model(&word).Updates(map[string]interface{}{
		"box":  newBox,
		"seen": word.Seen + 1,
		"due":  newBox <= 2,
	})
	return c.JSON(word)
}
