package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/auth/login
func Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Find user by email

	var user models.User
	result := database.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	fmt.Printf("Hash: [%s] len=%d\n", user.Password, len(user.Password))
	fmt.Printf("Input: [%s] len=%d\n", input.Password, len(input.Password))
	// Check password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		fmt.Printf(err.Error())
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// GET /api/auth/me
func Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"role":   user.Role,
		"avatar": user.Avatar,
	})
}

// POST /api/auth/register (create first admin)
func Register(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	input.Password = string(hashed)
	fmt.Printf("Saving hash: [%s]\n", input.Password)

	result := database.DB.Create(&input)
	if result.Error != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":    input.ID,
		"name":  input.Name,
		"email": input.Email,
		"role":  input.Role,
	})
}
