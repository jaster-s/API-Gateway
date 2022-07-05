package router

import (
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Gateway API")
	})

	lineOA := app.Group("/line_oa")
	LineOA(lineOA)

	user := app.Group("/user")
	User(user)

	webhookLine := app.Group("/line_webhook")
	Line(webhookLine)

	webhookFreshdesk := app.Group("/freshdesk_webhook")
	Freshdesk(webhookFreshdesk)

	config := app.Group("/config")
	Config(config)
}
