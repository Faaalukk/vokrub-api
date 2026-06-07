package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

// GET /api/customer
func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer
	database.DB.Find(&customers)
	return c.JSON(customers)
}

// GET /api/customer/:id
func GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Customer
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}
	return c.JSON(user)
}

// POST /api/customer
func CreateCustomer(c *fiber.Ctx) error {
	user := new(models.Customer)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	database.DB.Create(&user)
	return c.Status(201).JSON(user)
}

// DELETE /api/customer/:id
func DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Customer{}, id)
	return c.JSON(fiber.Map{"message": "Deleted"})
}
