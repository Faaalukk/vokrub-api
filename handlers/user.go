package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

// GET /api/customer
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)
	return c.JSON(users)
}

// GET /api/user/:id
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// POST /api/user
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	database.DB.Create(&user)
	return c.Status(201).JSON(user)
}

// DELETE /api/user/:id
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.User{}, id)
	return c.JSON(fiber.Map{"message": "Deleted"})
}
