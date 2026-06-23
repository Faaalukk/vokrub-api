package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

type SendOTPInput struct {
	Phone string `json:"phone"`
}

type VerifyOTPInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type SMSRequest struct {
	To      []string `json:"to"`
	Message string   `json:"message"`
	From    string   `json:"from"`
}

// POST /api/customer/auth/otp/send
func SendOTP(c *fiber.Ctx) error {
	input := new(SendOTPInput)
	if err := c.BodyParser(input); err != nil || input.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phone required"})
	}

	code := generateOTP()

	// Invalidate any existing unused OTPs for this phone
	database.DB.Model(&models.OTPCode{}).
		Where("phone = ? AND used = false", input.Phone).
		Updates(map[string]interface{}{"used": true})

	database.DB.Create(&models.OTPCode{
		Phone:     input.Phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if err := sendSMSKub(input.Phone, code); err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send OTP"})
	}

	return c.JSON(fiber.Map{"message": "OTP sent"})
}

// POST /api/customer/auth/otp/verify
func VerifyOTP(c *fiber.Ctx) error {
	input := new(VerifyOTPInput)
	if err := c.BodyParser(input); err != nil || input.Phone == "" || input.Code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phone and code required"})
	}

	var otp models.OTPCode
	if err := database.DB.Where(
		"phone = ? AND code = ? AND used = false AND expires_at > ?",
		input.Phone, input.Code, time.Now(),
	).First(&otp).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	database.DB.Model(&otp).Update("used", true)

	// Find existing phone identity or create new customer
	var identity models.CustomerIdentity
	var customer models.Customer

	if err := database.DB.Where("provider = ? AND provider_id = ?", "phone", input.Phone).First(&identity).Error; err == nil {
		database.DB.First(&customer, identity.CustomerID)
	} else {
		phone := input.Phone
		customer = models.Customer{
			Name:   input.Phone, // user can update name in profile
			Phone:  &phone,
			Plan:   "free",
			Status: "active",
		}
		database.DB.Create(&customer)
		database.DB.Create(&models.CustomerIdentity{
			CustomerID: customer.ID,
			Provider:   "phone",
			ProviderID: input.Phone,
		})
		notifyDiscordNewCustomer(&customer, "phone")
	}

	return c.JSON(fiber.Map{
		"token":    makeCustomerJWT(&customer),
		"customer": fiber.Map{"id": customer.ID, "name": customer.Name, "plan": customer.Plan, "streak": customer.Streak},
	})
}

func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// sendSMSKub sends an OTP via SMS_KUB.
// If SMSKUB_API_KEY is not set, logs to stdout (dev mode).
// TODO: verify exact request format at https://www.smskub.com/docs
func sendSMSKub(phone string, code string) error {
	apiKey := os.Getenv("SMSKUB_API_KEY")
	sender := os.Getenv("SMSKUB_SENDER")
	if sender == "" {
		sender = "Vokrub"
	}

	if apiKey == "" {
		fmt.Printf("[SMS_KUB DEV] To: %s | Code: %s\n", phone, code)
		return nil
	}

	payload, _ := json.Marshal(SMSRequest{
		To:      []string{phone},
		Message: fmt.Sprintf("รหัส OTP Vokrub: %s (ใช้ได้ 5 นาที)", code),
		From:    sender,
	})

	req, err := http.NewRequest("POST", "https://console.sms-kub.com/api/messages", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("smskub error: %d", resp.StatusCode)
	}
	return nil
}
