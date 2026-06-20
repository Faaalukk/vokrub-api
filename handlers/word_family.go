package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type CreateFamilyInput struct {
	Name    string           `json:"name"`
	WordID  uint             `json:"word_id"`  // existing word as initial member (optional)
	NewWord *CreateWordInput `json:"new_word"` // create-and-add a not-yet-stored word (optional)
}

type AddFamilyMemberInput struct {
	WordID  uint             `json:"word_id"`  // existing word (optional)
	NewWord *CreateWordInput `json:"new_word"` // create-and-add a not-yet-stored word (optional)
}

// resolveMemberWordID picks the word to add: an existing word_id, or an inline
// new_word that gets created (or matched if already stored). Returns 0 when
// neither is supplied. wordErr maps to 400 on missing text, 500 otherwise.
func resolveMemberWordID(customerID uint, wordID uint, newWord *CreateWordInput) (uint, error) {
	if wordID > 0 {
		return wordID, nil
	}
	if newWord != nil {
		return findOrCreateWord(customerID, *newWord)
	}
	return 0, nil
}

func mapWordErr(c *fiber.Ctx, err error) error {
	if err == ErrWordRequired {
		return c.Status(400).JSON(fiber.Map{"error": "word required"})
	}
	return c.Status(500).JSON(fiber.Map{"error": "Could not create word"})
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

	wordID, err := resolveMemberWordID(customerID.(uint), input.WordID, input.NewWord)
	if err != nil {
		return mapWordErr(c, err)
	}
	if wordID > 0 {
		member := models.WordFamilyMember{FamilyID: family.ID, WordID: wordID}
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
	wordID, err := resolveMemberWordID(customerID.(uint), input.WordID, input.NewWord)
	if err != nil {
		return mapWordErr(c, err)
	}
	if wordID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "word_id or new_word required"})
	}
	member := models.WordFamilyMember{FamilyID: family.ID, WordID: wordID}
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
