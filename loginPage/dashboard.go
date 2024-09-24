package loginpage

import (
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func DashboardMenu(c *fiber.Ctx) error {
	dash := &models.MenuBar{}
	if parsErr := c.BodyParser(dash); parsErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": parsErr.Error(),
		})
	}

	return c.JSON(dash)
}
