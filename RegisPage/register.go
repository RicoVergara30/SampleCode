package regispage

import (
	"sample/db"
	"sample/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Registration(c *fiber.Ctx) error {
	// Create a new instance of Registration model
	register := &models.Registration{}

	// Parse the body into the register object
	if err := c.BodyParser(&register); err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "could not hash password",
			"details": err.Error(),
		})
	}
	register.Password = string(hashedPassword)

	// Insert the new registration into the database
	if err := db.DB.Debug().Exec("INSERT INTO public.registration ( email, username, password) VALUES (?, ?, ?, ?)", register.Email, register.Username, register.Password).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "registration failed",
			"details": err.Error(),
		})
	}

	// Return the registration details as response
	return c.JSON(register)
}
