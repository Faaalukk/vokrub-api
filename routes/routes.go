package routes

import (
	"github.com/Faaalukk/vokrub-api.git/handlers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	customer := api.Group("/customer")
	customer.Get("/", handlers.GetCustomers)
	customer.Get("/:id", handlers.GetCustomer)
	customer.Post("/", handlers.CreateCustomer)
	customer.Delete("/:id", handlers.DeleteCustomer)

	user := api.Group("/user")
	user.Get("/", handlers.GetUsers)
	user.Get("/:id", handlers.GetUser)
	user.Post("/", handlers.CreateUser)
	user.Delete("/", handlers.DeleteUser)
}
