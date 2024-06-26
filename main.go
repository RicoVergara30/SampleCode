package main

import (
	"fmt"
	"log"
	"sample/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {

	yearNow := time.Now().Format("2006")
	appTitle := fmt.Sprintf("FDSAP RICO SAMPLE CODE  - %s", yearNow)

	app := fiber.New(fiber.Config{
		ServerHeader: appTitle,
		AppName:      appTitle,
		Views:        html.New("./template", ".html"),
	})

	if app == nil {
		log.Fatal("Failed to initialize Fiber app")
	}

	// Serve static files from the specified directory
	app.Static("/", "./assets/css")

	// Middleware to recover from panics
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Specify allowed headers
	}))

	// Logger middleware with custom format and time format
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white}| ${red}${latency} ${white}\n",
		TimeFormat: "01/02/2006 3:04 PM",
	}))

	// Setup routes - Ensure this function is properly defined and routes are correctly set up
	routes.SetupRoutes(app)

	// Start server on port 1000
	err := app.Listen(":1000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
