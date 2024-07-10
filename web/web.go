package web

import (
	"github.com/gofiber/fiber/v2"
)

// |-----|
// | WEB |
// |-----|

// Show Login
func ShowLogin(c *fiber.Ctx) error {

	return c.Render("login", fiber.Map{
		"title": "Log-In",
		// "logsResponse": logResult,
	})
}

// Show Registewr
func ShowRegister(c *fiber.Ctx) error {

	return c.Render("registration", fiber.Map{
		"title": "Registration",
		// "logsResponse": logResult,
	})
}
