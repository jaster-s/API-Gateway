package controller

import (
	"fli-gateway-api/lib"
	"fli-gateway-api/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ConfigController struct {
	ExpirationDate int `json:"expiration_date"`
	MaxLineAccount int `json:"max_line_account"`
}

func (ConfigController) AddConfig(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	Data := new(ConfigController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	if Data.ExpirationDate <= 0 || Data.MaxLineAccount <= 0 {
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}

	configModel := new(models.ConfigModel)
	err := configModel.AddConfig(Data.ExpirationDate, Data.MaxLineAccount)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}

	return c.JSON(map[string]string{
		"message": "Save successfully.",
	})
}
