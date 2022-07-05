package controller

import (
	"fli-gateway-api/lib"
	"fli-gateway-api/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type LineOAController struct {
	UserID           string `json:"user_id"`
	LINEObjID        string `json:"line_obj_id"`
	Name             string `json:"name"`
	ClientID         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	Responder_id     int    `json:"responder_id"`
	Group_id         int    `json:"group_id"`
	RecoveryMail     string `json:"recovery_mail"`
	AutoMessageReply string `json:"auto_message_reply"`
}

func (LineOAController) ListLineOA(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	UserID := fmt.Sprintf("%v", c.Locals("userID"))

	lineOAModel := new(models.LineOAModel)
	lineOA, err := lineOAModel.List(UserID)

	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)

		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}

	return c.JSON(lineOA)
}

func (LineOAController) AddLINE_OA(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(LineOAController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	UserID := fmt.Sprintf("%v", c.Locals("userID"))
	maxLineAccount := c.Locals("maxLineAccount").(int)

	lineOAModel := new(models.LineOAModel)
	count_lineOA := lineOAModel.Count(UserID)
	if maxLineAccount >= count_lineOA {
		verifyLINEOA, _ := lineOAModel.ValidateLINE(Data.ClientID, Data.ClientSecret)
		if verifyLINEOA == true {
			lineOA, err := lineOAModel.AddLINEAccount(UserID, Data.Name, Data.ClientID, Data.ClientSecret, Data.RecoveryMail, Data.AutoMessageReply, Data.Group_id, Data.Responder_id)
			if lineOA == nil {
				return c.JSON(map[string]interface{}{
					"Result":  false,
					"LINEOA":  true,
					"message": "Can't add clientID because this clientID is already in use. Please try again.",
				})
			} else if err != nil {
				_logger.Error(fmt.Sprintf("%v", err), false)
				return c.JSON(map[string]interface{}{
					"Result":  false,
					"LINEOA":  false,
					"message": "Something went wrong. Please try again.",
				})
			} else {
				mail := new(models.SendMail)
				fmt.Println("Add LINE Office Account")
				fmt.Println("Add recovery email, sending mail.")
				getDetailLINEOA := new(models.LineModels)
				detail_LINEOA := getDetailLINEOA.GetProfileLINEOA(Data.ClientID)
				mail.Add_recoveryEmail(lineOA.RecoveryMail, detail_LINEOA.DisplayName)
				return c.JSON(map[string]interface{}{
					"Result":  true,
					"LINEOA":  true,
					"message": "",
				})
			}
		} else {
			return c.JSON(map[string]interface{}{
				"Result":  true,
				"LINEOA":  false,
				"message": "ClientID is invalid.",
			})
		}
	} else {
		return c.JSON(map[string]interface{}{
			"Result":  false,
			"LINEOA":  false,
			"message": "Adding LINE Official Account of your limit.",
		})
	}
}

func (LineOAController) DelateLINE_OA(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(LineOAController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	lineOAModel := new(models.LineOAModel)
	err := lineOAModel.DeleteLINEAccount(Data.LINEObjID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]interface{}{
			"status":  false,
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	} else {
		return c.JSON(map[string]interface{}{
			"status":  true,
			"message": "Successfully deleted.",
		})
	}
}

func (LineOAController) EditLINE_OA(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(LineOAController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	lineOAModel := new(models.LineOAModel)
	detailLineOA, old_email, err := lineOAModel.EditAccount(Data.LINEObjID, Data.RecoveryMail, Data.AutoMessageReply, Data.Group_id, Data.Responder_id)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]interface{}{
			"status":  false,
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	} else {
		if detailLineOA.RecoveryMail != old_email {
			mail := new(models.SendMail)
			fmt.Println("edit recovery email, sending mail.")
			getDetailLINEOA := new(models.LineModels)
			detail_LINEOA := getDetailLINEOA.GetProfileLINEOA(detailLineOA.ClientID)
			mail.Add_recoveryEmail(detailLineOA.RecoveryMail, detail_LINEOA.DisplayName)
		}
		return c.JSON(map[string]interface{}{
			"status":  true,
			"message": "Edited successfully.",
		})
	}
}

func (LineOAController) Webhook(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(LineOAController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	lineOAModel := new(models.LineOAModel)
	webhookURL, err := lineOAModel.Webhook(Data.LINEObjID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return nil
	} else {
		return c.JSON(map[string]interface{}{
			"webhook": webhookURL,
		})
	}

}
