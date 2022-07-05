package router

import (
	"fli-gateway-api/controller"
	"fli-gateway-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Config(router fiber.Router) {
	configController := new(controller.ConfigController)

	router.Post("/add", middleware.VerifySignatureAdmin, configController.AddConfig)
}
