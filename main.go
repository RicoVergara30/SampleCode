package main

import (
	"log"
	"sample/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	app.Static("/", "./assets/...")
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Specify allowed headers
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white}| ${red}${latency} ${white}\n",
		TimeFormat: "01/02/2006 3:04 PM",
	}))

	log.Fatal(app.Listen(":1000"))
}
