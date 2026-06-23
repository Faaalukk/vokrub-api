package handlers

import (
	"os"
	"time"

	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/golang-jwt/jwt/v5"
)

// makeCustomerJWT creates a signed 7-day JWT for a customer.
func makeCustomerJWT(customer *models.Customer) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customer.ID,
		"type":        "customer",
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

// findOrCreateOAuthCustomer resolves a customer from an OAuth identity.
// Order: existing identity → existing email (link) → create new.
func findOrCreateOAuthCustomer(provider, providerID, name, email, image string) (*models.Customer, error) {
	var identity models.CustomerIdentity
	if err := database.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&identity).Error; err == nil {
		var customer models.Customer
		if err := database.DB.First(&customer, identity.CustomerID).Error; err != nil {
			return nil, err
		}
		return &customer, nil
	}

	// No identity yet — try linking by email
	if email != "" {
		var customer models.Customer
		if err := database.DB.Where("email = ?", email).First(&customer).Error; err == nil {
			database.DB.Create(&models.CustomerIdentity{CustomerID: customer.ID, Provider: provider, ProviderID: providerID})
			return &customer, nil
		}
	}

	// Create new customer
	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}
	customer := models.Customer{
		Name:   name,
		Email:  emailPtr,
		Image:  image,
		Plan:   "free",
		Status: "active",
	}
	if err := database.DB.Create(&customer).Error; err != nil {
		return nil, err
	}
	database.DB.Create(&models.CustomerIdentity{CustomerID: customer.ID, Provider: provider, ProviderID: providerID})
	notifyDiscordNewCustomer(&customer, provider)
	return &customer, nil
}
