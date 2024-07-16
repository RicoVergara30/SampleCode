package loginpage

import (
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func LogoutUser(c *fiber.Ctx) error {
	logout := &models.LogOut{}

	if err := c.BodyParser(&logout); err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(logout)
}
