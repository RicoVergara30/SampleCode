package routes

import (
	loginpage "sample/loginPage"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	apiEndpoint := app.Group("/api")
	v1Endpoint := apiEndpoint.Group("/v1")

	v1Endpoint.Post("/log", loginpage.LoginPage)
}
