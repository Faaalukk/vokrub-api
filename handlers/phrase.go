package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type CreateCategoryInput struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Hue  int    `json:"hue"`
}

type CreateSentenceInput struct {
	Text    string `json:"text"`
	Meaning string `json:"meaning"`
	Note    string `json:"note"`
}

// GET /api/phrase/category
func GetCategories(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	var categories []models.Category
	database.DB.Preload("Sentences").Where("customer_id = ?", customerID).Find(&categories)
	return c.JSON(categories)
}

// POST /api/phrase/category
func CreateCategory(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	input := new(CreateCategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	cat := models.Category{
		CustomerID: customerID.(uint),
		Name:       input.Name,
		Icon:       input.Icon,
		Hue:        input.Hue,
	}
	database.DB.Create(&cat)
	cat.Sentences = []models.Sentence{}
	return c.Status(201).JSON(cat)
}

// DELETE /api/phrase/category/:id
func DeleteCategory(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	database.DB.Where("customer_id = ?", customerID).Delete(&models.Category{}, id)
	return c.JSON(fiber.Map{"message": "Deleted"})
}

// POST /api/phrase/category/:id/sentence
func AddSentence(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	catID := c.Params("id")

	var cat models.Category
	if result := database.DB.Where("id = ? AND customer_id = ?", catID, customerID).First(&cat); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	input := new(CreateSentenceInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	sentence := models.Sentence{
		CategoryID: cat.ID,
		Text:       input.Text,
		Meaning:    input.Meaning,
		Note:       input.Note,
	}
	database.DB.Create(&sentence)
	return c.Status(201).JSON(sentence)
}

// PUT /api/phrase/sentence/:id
func UpdateSentence(c *fiber.Ctx) error {
	id := c.Params("id")
	var sentence models.Sentence
	if result := database.DB.First(&sentence, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Sentence not found"})
	}

	input := new(CreateSentenceInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	database.DB.Model(&sentence).Updates(map[string]interface{}{
		"text":    input.Text,
		"meaning": input.Meaning,
		"note":    input.Note,
	})
	return c.JSON(sentence)
}

// DELETE /api/phrase/sentence/:id
func DeleteSentence(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Sentence{}, id)
	return c.JSON(fiber.Map{"message": "Deleted"})
}
