package routes

import (
	"github.com/Faaalukk/vokrub-api.git/handlers"
	"github.com/Faaalukk/vokrub-api.git/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// ── Admin auth ─────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", middleware.Protected, handlers.Register)
	auth.Get("/me", middleware.Protected, handlers.Me)

	// ── Customer auth (mobile) ──────────────────────────────────
	customerAuth := api.Group("/customer/auth")
	customerAuth.Post("/login", handlers.CustomerLogin)
	customerAuth.Post("/register", handlers.CustomerRegister)
	customerAuth.Get("/me", middleware.CustomerProtected, handlers.CustomerMe)

	// ── Admin: customers ────────────────────────────────────────
	customer := api.Group("/customer")
	customer.Get("/", middleware.Protected, handlers.GetCustomers)
	customer.Get("/:id", middleware.Protected, handlers.GetCustomer)
	customer.Post("/", middleware.Protected, handlers.CreateCustomer)
	customer.Delete("/:id", middleware.Protected, handlers.DeleteCustomer)

	// ── Admin: users ────────────────────────────────────────────
	user := api.Group("/user")
	user.Get("/", middleware.Protected, handlers.GetUsers)
	user.Get("/:id", middleware.Protected, handlers.GetUser)
	user.Post("/", middleware.Protected, handlers.CreateUser)
	user.Delete("/", middleware.Protected, handlers.DeleteUser)

	// ── Admin: analytics / metrics ──────────────────────────────
	api.Get("/metrics", middleware.Protected, handlers.GetMetrics)
	api.Get("/analytics/top-words", middleware.Protected, handlers.GetTopWords)
	api.Get("/analytics/transactions", middleware.Protected, handlers.GetTransactions)

	// ── Customer: words (mobile) ────────────────────────────────
	word := api.Group("/word", middleware.CustomerProtected)
	word.Get("/", handlers.GetWords)
	word.Post("/", handlers.CreateWord)
	word.Get("/:id", handlers.GetWord)
	word.Put("/:id", handlers.UpdateWord)
	word.Delete("/:id", handlers.DeleteWord)
	word.Post("/:id/review", handlers.ReviewWord)

	// ── Customer: phrases (mobile) ──────────────────────────────
	phrase := api.Group("/phrase", middleware.CustomerProtected)
	phrase.Get("/category", handlers.GetCategories)
	phrase.Post("/category", handlers.CreateCategory)
	phrase.Delete("/category/:id", handlers.DeleteCategory)
	phrase.Post("/category/:id/sentence", handlers.AddSentence)
	phrase.Put("/sentence/:id", handlers.UpdateSentence)
	phrase.Delete("/sentence/:id", handlers.DeleteSentence)
}
