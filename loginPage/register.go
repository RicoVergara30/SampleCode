package loginpage

import (
	"sample/db"
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func Registration(c *fiber.Ctx) error {

	register := &models.Registration{}

	if err := c.BodyParser(&register); err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, err := db.HashPassword((register.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "could not hash password",
			"details": err.Error(),
		})
	}
	register.Password = string(hashedPassword)

	if err := db.Database.Debug().Exec("INSERT INTO public.regist( email, username, password)VALUES ( ?, ?, ?)", register.Email, register.Username, register.Password).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "registration failed",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registration successful",
		"details": "200",
	})
}
