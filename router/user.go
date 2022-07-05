package router

import (
	"fli-gateway-api/controller"
	"fli-gateway-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func User(router fiber.Router) {
	userController := new(controller.UserController)

	router.Post("/create", middleware.VerifySignatureAdmin, userController.Create)
	router.Post("/create/trial", middleware.VerifySignatureCreateTrial,userController.Create)
	router.Post("/activate_key/verify", middleware.VerifyActivateKey, userController.Verify)
	router.Post("/remove_user", middleware.VerifySignatureAdmin, userController.Delete)
	router.Post("/activate_key/regenerate", middleware.VerifySignatureAdmin, userController.Regenerate)
	router.Post("/activate_key/status", middleware.VerifySignatureAdmin, userController.ChangeStatus)
	router.Post("/change/MaxLineAccount", middleware.VerifySignatureAdmin, userController.ChangeMaxLineAccount)
	router.Post("/change/Expiration_ActivateKey", middleware.VerifySignatureAdmin, userController.ChangeExpirationActivateKey)
	router.Get("/check/Expiration_ActivateKey", middleware.VerifySignatureListLINE_OA, userController.CheckExpirationActivateKey)

}
