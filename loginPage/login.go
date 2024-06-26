package loginpage

import (
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func LoginPage(c *fiber.Ctx) error {
	log := &models.LoginPage{}
	if err := c.BodyParser(&log); err != nil {
		return c.JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	return c.JSON(log)
}
