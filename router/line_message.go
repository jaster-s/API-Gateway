package router

import (
	"fli-gateway-api/controller"
	"fli-gateway-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Line(router fiber.Router) {
	lineMessage := new(controller.LineMessage)

	router.Post("/push/*", middleware.VerifyingSignaturesLINE, lineMessage.BuildMessage)
}
