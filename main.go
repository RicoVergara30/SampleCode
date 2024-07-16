package main

import (
	"fmt"
	"log"
	"sample/db"
	"sample/models"
	"sample/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
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

	// .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Logger middleware with custom format and time format
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white}| ${red}${latency} ${white}\n",
		TimeFormat: "01/02/2006 3:04 PM",
	}))

	// Database connection
	err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}

	// Auto-migrate models
	db.Database.AutoMigrate(&models.Registration{})

	// Setup routes - Ensure this function is properly defined and routes are correctly set up
	routes.SetupRoutes(app)

	// Start server
	err = app.Listen(":2000")
	if err != nil {
		log.Fatal("Error starting the server")
	}

}
