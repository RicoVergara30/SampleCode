package loginpage

import (
	"sample/db"
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func LoginPage(c *fiber.Ctx) error {
	log := &models.LoginPage{}

	// Parse the request body into the log struct
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := &models.Registration{}
	if err := db.Database.Debug().Table("regist").Where("username = ?", log.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Compare the provided password with the stored password directly
	if user.Password != log.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}
