package routes

import (
	"github.com/Faaalukk/vokrub-api.git/handlers"
	"github.com/Faaalukk/vokrub-api.git/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// Auth - public
	auth := api.Group("/auth")
	auth.Post("/login", middleware.Protected, handlers.Login)
	auth.Post("/register", middleware.Protected, handlers.Register)
	auth.Get("/me", middleware.Protected, handlers.Me)

	customer := api.Group("/customer")
	customer.Get("/", middleware.Protected, handlers.GetCustomers)
	customer.Get("/:id", middleware.Protected, handlers.GetCustomer)
	customer.Post("/", middleware.Protected, handlers.CreateCustomer)
	customer.Delete("/:id", middleware.Protected, handlers.DeleteCustomer)

	user := api.Group("/user")
	user.Get("/", middleware.Protected, handlers.GetUsers)
	user.Get("/:id", middleware.Protected, handlers.GetUser)
	user.Post("/", middleware.Protected, handlers.CreateUser)
	user.Delete("/", middleware.Protected, handlers.DeleteUser)
}
