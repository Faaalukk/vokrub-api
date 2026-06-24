package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ErrWordRequired is returned when an inline word has no text.
var ErrWordRequired = errors.New("word required")

// findOrCreateWord resolves a word for the customer: returns the existing word
// with the same text if present, otherwise creates a new one. Returns the ID.
// Used when callers (e.g. family creation) want to add a word that may not be
// stored yet without forcing the user to create it separately first.
func findOrCreateWord(customerID uint, input CreateWordInput) (uint, error) {
	term := strings.ToLower(strings.TrimSpace(input.Word))
	if term == "" {
		return 0, ErrWordRequired
	}

	var existing models.Word
	if result := database.DB.Where("customer_id = ? AND word = ?", customerID, term).First(&existing); result.Error == nil {
		return existing.ID, nil
	}

	synonyms := input.Synonyms
	if synonyms == nil {
		synonyms = models.StringSlice{}
	}
	word := models.Word{
		CustomerID: customerID,
		Word:       term,
		Pos:        input.Pos,
		Meaning:    input.Meaning,
		Note:       input.Note,
		CategoryID: input.CategoryID,
		Synonyms:   synonyms,
		Box:        1,
		Seen:       0,
		Due:        true,
		Added:      time.Now().Format("2006-01-02"),
	}
	if result := database.DB.Create(&word); result.Error != nil {
		return 0, result.Error
	}
	return word.ID, nil
}

// updateStreak increments or resets the customer streak based on activity date.
// Returns the new streak value.
func updateStreak(db *gorm.DB, customerID uint) int {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	var customer models.Customer
	if err := db.First(&customer, customerID).Error; err != nil {
		return 0
	}
	if customer.LastActiveDate == today {
		return customer.Streak // already counted today
	}

	newStreak := 1
	if customer.LastActiveDate == yesterday {
		newStreak = customer.Streak + 1
	}

	db.Model(&customer).Updates(map[string]any{
		"streak":           newStreak,
		"last_active_date": today,
	})
	return newStreak
}

type CreateWordInput struct {
	Word       string             `json:"word"`
	Pos        models.StringSlice `json:"pos"`
	Translate  string             `json:"translate"`
	Meaning    string             `json:"meaning"`
	Note       string             `json:"note"`
	CategoryID *uint              `json:"category_id"`
	Synonyms   models.StringSlice `json:"synonyms"`
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

// GET /api/word/check?word=xxx
func CheckWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	term := strings.ToLower(strings.TrimSpace(c.Query("word")))
	if term == "" {
		return c.Status(400).JSON(fiber.Map{"error": "word query param required"})
	}
	var existing models.Word
	result := database.DB.Where("customer_id = ? AND word = ?", customerID, term).First(&existing)
	if result.Error == nil {
		return c.JSON(fiber.Map{"duplicate": true, "word": existing})
	}
	return c.JSON(fiber.Map{"duplicate": false})
}

// POST /api/word
func CreateWord(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	input := new(CreateWordInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	input.Word = strings.ToLower(strings.TrimSpace(input.Word))
	if input.Word == "" {
		return c.Status(400).JSON(fiber.Map{"error": "word required"})
	}

	var existing models.Word
	if result := database.DB.Where("customer_id = ? AND word = ?", customerID, input.Word).First(&existing); result.Error == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Word already exists", "word": existing})
	}

	synonyms := input.Synonyms
	if synonyms == nil {
		synonyms = models.StringSlice{}
	}
	word := models.Word{
		CustomerID: customerID.(uint),
		Word:       input.Word,
		Pos:        input.Pos,
		Translate:  input.Translate,
		Meaning:    input.Meaning,
		Note:       input.Note,
		CategoryID: input.CategoryID,
		Synonyms:   synonyms,
		Box:        1,
		Seen:       0,
		Due:        true,
		Added:      time.Now().Format("2006-01-02"),
	}

	if result := database.DB.Create(&word); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create word"})
	}
	streak := updateStreak(database.DB, customerID.(uint))
	return c.Status(201).JSON(fiber.Map{"word": word, "streak": streak})
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

	input.Word = strings.ToLower(strings.TrimSpace(input.Word))

	synonyms := input.Synonyms
	if synonyms == nil {
		synonyms = models.StringSlice{}
	}
	database.DB.Model(&word).Updates(map[string]interface{}{
		"word":        input.Word,
		"pos":         input.Pos,
		"translate":   input.Translate,
		"meaning":     input.Meaning,
		"note":        input.Note,
		"category_id": input.CategoryID,
		"synonyms":    synonyms,
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
