package main

import (
	"fli-gateway-api/lib"
	"fli-gateway-api/router"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	_logger := new(lib.Logger)
	PORT := os.Getenv("APP_PORT")
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Register Router
	router.New(app)

	err := app.Listen(fmt.Sprintf(":%s", PORT))

	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot start server: %v", err), true)
	}
}
