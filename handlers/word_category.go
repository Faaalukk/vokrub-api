package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type CreateWordCategoryInput struct {
	Name  string `json:"name"`
	Color int    `json:"color"`
}

// GET /api/word/category
func GetWordCategories(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	var cats []models.WordCategory
	database.DB.Where("customer_id = ?", customerID).Order("name asc").Find(&cats)
	return c.JSON(cats)
}

// POST /api/word/category
func CreateWordCategory(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	input := new(CreateWordCategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name required"})
	}

	var existing models.WordCategory
	if result := database.DB.Where("customer_id = ? AND name = ?", customerID, input.Name).First(&existing); result.Error == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Category already exists", "category": existing})
	}

	cat := models.WordCategory{
		CustomerID: customerID.(uint),
		Name:       input.Name,
		Color:      input.Color,
	}
	if result := database.DB.Create(&cat); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create category"})
	}
	return c.Status(201).JSON(cat)
}

// PUT /api/word/category/:id
func UpdateWordCategory(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var cat models.WordCategory
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&cat); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	input := new(CreateWordCategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	database.DB.Model(&cat).Updates(map[string]interface{}{
		"name":  input.Name,
		"color": input.Color,
	})
	return c.JSON(cat)
}

// DELETE /api/word/category/:id
func DeleteWordCategory(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")

	database.DB.Where("customer_id = ?", customerID).Delete(&models.WordCategory{}, id)
	// Unlink words from this category
	database.DB.Model(&models.Word{}).Where("customer_id = ? AND category_id = ?", customerID, id).
		Update("category_id", nil)
	return c.JSON(fiber.Map{"message": "Deleted"})
}
