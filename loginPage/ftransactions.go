package loginpage

import (
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func FtransactionHandler(c *fiber.Ctx) error {
	ft := &models.Ftransaction{}

	// Parse the request body into the ft struct
	if err := c.BodyParser(ft); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ft)
}
