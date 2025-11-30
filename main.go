package main

import (
	"log"
	"os"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/routes"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	} else {
		port = ":" + port
	}

	return port
}

func main() {

	database.Connect()

	app := fiber.New(fiber.Config{
		// Trust proxy headers
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
	})

	// Initialize default config
	app.Use(logger.New())

	// Middleware
	// Allow origins can be configured via the ALLOWED_ORIGINS env var (comma-separated).
	allowedOrigins := utils.Env("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://192.168.0.70:3000,https://certikiosk.up.railway.app"
	}
	log.Printf("[info] CORS allowed origins: %s", allowedOrigins)

	// CORS Configuration - Must be before routes
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Type",
		MaxAge:           86400,
	}))

	// Health check endpoint (doesn't require DB)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "certikiosk-api",
		})
	})

	// routes.Setup(app)
	routes.Setup(app)

	log.Fatal(app.Listen(getPort()))

}
