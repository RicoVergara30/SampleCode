package web

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// |-----|
// | WEB |
// |-----|
// Show Login
func ShowLogin(c *fiber.Ctx) error {
	// logResult := &[]models.LoginPage{}
	// db.DB.Raw("SELECT * FROM rbi_instapay.view_transactions").Scan(&logResult)
	fmt.Println("error checker")
	return c.Render("footer", fiber.Map{
		"title": "Log-In",
		// "logsResponse": logResult,
	})
}
