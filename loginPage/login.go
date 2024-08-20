package loginpage

import (
	"sample/db"
	"sample/middleware/jwttoken"
	"sample/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoginPages(c *fiber.Ctx) error {
	log := &models.LoginPage{}

	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := &models.Registration{}
	if err := db.Database.Debug().Table("regist").Where("username = ?", log.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username",
		})
	}

	//Decypt
	UnhashedPassword, err := db.HashPassword((log.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "could not hash password",
			"details": err.Error(),
		})
	}
	log.Password = string(UnhashedPassword)

	// Generate token with 3-second expiration
	tokenString, err := jwttoken.GenerateTokens(user.Username, 3*time.Second)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "could not generate token",
			"details": err.Error(),
		})
	}

	response := &models.ResponseLogin{
		Username: user.Username,
		Password: user.Password,
		Status:   true,
		Token:    tokenString,
	}
	if err := db.Database.Debug().Exec("INSERT INTO public.login(username, password, status, token) VALUES (?, ?, ?, ?)", log.Username, log.Password, response.Status, response.Token).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
		})

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
		"details": "200",
	})
}

// package loginpage

// import (
// 	"sample/db"
// 	"sample/models"

// 	"github.com/gofiber/fiber/v2"
// )

// func LoginPage(c *fiber.Ctx) error {
// 	log := &models.LoginPage{}

// 	// Parse the request body into the log struct
// 	if err := c.BodyParser(&log); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}
// 	UnhashedPassword, err := db.HashAndComparePassword()((log.Password))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   "could not hash password",
// 			"details": err.Error(),
// 		})
// 	}
// 	log.Password = string(UnhashedPassword)

// 	user := &models.Registration{}
// 	if err := db.Database.Debug().Table("regist").Where("username = ?", log.Username).First(&user).Error; err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid username or password",
// 		})
// 	}

// 	// Compare the provided password with the stored password directly
// 	if user.Password != log.Password {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid username or password",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Login successful",
// 	})
// }
