package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type CreateFamilyInput struct {
	Name   string `json:"name"`
	WordID uint   `json:"word_id"` // initial member
}

type AddFamilyMemberInput struct {
	WordID uint `json:"word_id"`
}

// GET /api/word/family
func GetWordFamilies(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	var families []models.WordFamily
	database.DB.Where("customer_id = ?", customerID).Preload("Members").Find(&families)
	return c.JSON(families)
}

// POST /api/word/family
func CreateWordFamily(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	input := new(CreateFamilyInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	family := models.WordFamily{
		CustomerID: customerID.(uint),
		Name:       input.Name,
	}
	if result := database.DB.Create(&family); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create family"})
	}

	if input.WordID > 0 {
		member := models.WordFamilyMember{FamilyID: family.ID, WordID: input.WordID}
		database.DB.Create(&member)
	}

	database.DB.Preload("Members").First(&family, family.ID)
	return c.Status(201).JSON(family)
}

// PUT /api/word/family/:id
func UpdateWordFamily(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var family models.WordFamily
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&family); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Family not found"})
	}
	input := new(CreateFamilyInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	database.DB.Model(&family).Update("name", input.Name)
	return c.JSON(family)
}

// DELETE /api/word/family/:id
func DeleteWordFamily(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var family models.WordFamily
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&family); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Family not found"})
	}
	database.DB.Where("family_id = ?", family.ID).Delete(&models.WordFamilyMember{})
	database.DB.Delete(&family)
	return c.JSON(fiber.Map{"message": "Deleted"})
}

// POST /api/word/family/:id/member
func AddFamilyMember(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	var family models.WordFamily
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&family); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Family not found"})
	}
	input := new(AddFamilyMemberInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	member := models.WordFamilyMember{FamilyID: family.ID, WordID: input.WordID}
	if result := database.DB.Create(&member); result.Error != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Already a member"})
	}
	database.DB.Preload("Members").First(&family, family.ID)
	return c.Status(201).JSON(family)
}

// DELETE /api/word/family/:id/member/:wordId
func RemoveFamilyMember(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")
	id := c.Params("id")
	wordId := c.Params("wordId")
	var family models.WordFamily
	if result := database.DB.Where("id = ? AND customer_id = ?", id, customerID).First(&family); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Family not found"})
	}
	database.DB.Where("family_id = ? AND word_id = ?", family.ID, wordId).Delete(&models.WordFamilyMember{})
	database.DB.Preload("Members").First(&family, family.ID)
	return c.JSON(family)
}
