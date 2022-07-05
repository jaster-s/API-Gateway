package router

import (
	"fli-gateway-api/controller"
	"fli-gateway-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func LineOA(router fiber.Router) {
	lineOACtrl := new(controller.LineOAController)

	router.Get("/list/LineOA", middleware.VerifySignatureListLINE_OA, lineOACtrl.ListLineOA)
	router.Post("/add/LineOA", middleware.VerifySignatureActivateKey, lineOACtrl.AddLINE_OA)
	router.Post("/remove/LineOA", middleware.VerifySignatureActivateKey, lineOACtrl.DelateLINE_OA)
	router.Post("/edit/LineOA", middleware.VerifySignatureActivateKey, lineOACtrl.EditLINE_OA)
	router.Post("/webhook", middleware.VerifySignatureActivateKey, lineOACtrl.Webhook)
}
