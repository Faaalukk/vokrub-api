package handlers

import (
	"github.com/Faaalukk/vokrub-api.git/database"
	"github.com/Faaalukk/vokrub-api.git/models"
	"github.com/gofiber/fiber/v2"
)

// GET /api/metrics
func GetMetrics(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Customer{}).Count(&total)

	var proMonthly int64
	database.DB.Model(&models.Customer{}).Where("plan = ?", "pro_monthly").Count(&proMonthly)

	var proAnnual int64
	database.DB.Model(&models.Customer{}).Where("plan = ?", "pro_annual").Count(&proAnnual)

	var activeCount int64
	database.DB.Model(&models.Customer{}).Where("status = ?", "active").Count(&activeCount)

	mrr := float64(proMonthly)*6 + float64(proAnnual)*(58.0/12.0)
	proCount := proMonthly + proAnnual

	return c.JSON(fiber.Map{
		"mrr":            mrr,
		"mrr_delta":      0,
		"customers":      total,
		"customers_delta": 0,
		"pro_count":      proCount,
		"pro_delta":      0,
		"active_today":   activeCount,
		"active_delta":   0,
	})
}

// GET /api/analytics/top-words
func GetTopWords(c *fiber.Ctx) error {
	type WordCount struct {
		Word  string `json:"word"`
		Count int64  `json:"count"`
	}

	var results []WordCount
	database.DB.Model(&models.Word{}).
		Select("word, count(*) as count").
		Group("word").
		Order("count desc").
		Limit(8).
		Scan(&results)

	return c.JSON(results)
}

// GET /api/analytics/transactions
func GetTransactions(c *fiber.Ctx) error {
	type TransactionRow struct {
		CustomerID uint   `json:"customer_id"`
		Name       string `json:"name"`
		Plan       string `json:"plan"`
		Amount     int    `json:"amount"`
		Status     string `json:"status"`
	}

	var rows []TransactionRow
	database.DB.Model(&models.Customer{}).
		Select("id as customer_id, name, plan, CASE WHEN plan = 'pro_annual' THEN 58 ELSE 6 END as amount, 'paid' as status").
		Where("plan != ?", "free").
		Order("updated_at desc").
		Limit(10).
		Scan(&rows)

	return c.JSON(rows)
}
