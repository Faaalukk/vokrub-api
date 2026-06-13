package handlers

import (
	"os"
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomerLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustomerRegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/customer/auth/login
func CustomerLogin(c *fiber.Ctx) error {
	input := new(CustomerLoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	var customer models.Customer
	if result := database.DB.Where("email = ?", input.Email).First(&customer); result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customer.ID,
		"type":        "customer",
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token":    tokenString,
		"customer": fiber.Map{"id": customer.ID, "name": customer.Name, "email": customer.Email, "plan": customer.Plan, "streak": customer.Streak},
	})
}

// POST /api/customer/auth/register
func CustomerRegister(c *fiber.Ctx) error {
	input := new(CustomerRegisterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}

	email := input.Email
	customer := models.Customer{
		Name:     input.Name,
		Email:    &email,
		Password: string(hashed),
		Plan:     "free",
		Status:   "active",
	}

	if result := database.DB.Create(&customer); result.Error != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	return c.Status(201).JSON(fiber.Map{
		"id": customer.ID, "name": customer.Name, "email": customer.Email, "plan": customer.Plan,
	})
}

// GET /api/customer/auth/me
func CustomerMe(c *fiber.Ctx) error {
	customerID := c.Locals("customer_id")

	var customer models.Customer
	if result := database.DB.First(&customer, customerID); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}

	wordCount := int64(0)
	database.DB.Model(&models.Word{}).Where("customer_id = ?", customer.ID).Count(&wordCount)

	return c.JSON(fiber.Map{
		"id":     customer.ID,
		"name":   customer.Name,
		"email":  customer.Email, // *string — null for phone/oauth-only customers
		"phone":  customer.Phone, // *string — null for email customers
		"plan":   customer.Plan,
		"streak": customer.Streak,
		"words":  wordCount,
		"status": customer.Status,
		"image":  customer.Image,
	})
}
