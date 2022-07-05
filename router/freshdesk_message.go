package router

import (
	"fli-gateway-api/controller"
	"fli-gateway-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Freshdesk(router fiber.Router) {
	freshdeskMessage := new(controller.FreshdeskMessage)

	router.Post("/create", middleware.VerifySignatureActivateKey, freshdeskMessage.NewTickets)
	router.Post("/reply", middleware.VerifySignatureActivateKey, freshdeskMessage.ReplytoLine)
	router.Static("/assets/resize", "./assets/resize")
	router.Post("/appInstall", middleware.VerifySignatureActivateKey, freshdeskMessage.AppInstall)
	router.Post("/appUninstall", middleware.VerifySignatureActivateKey, freshdeskMessage.AppUninstall)
}
