package controller

import (
	"fli-gateway-api/lib"
	"fli-gateway-api/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserID         string `json:"user_id"`
	MaxLineAccount int    `json:"max_line_account"`
	ActivateKey    string `json:"activate_key"`
	Status         bool   `json:"status"`
	SubDomain      string `json:"sub_domain"`
	ApiKey         string `json:"api_key"`
	ExpirationDate string `json:"expiration_date"`
	Company        string `json:"company"`
	Mail           string `json:"mail"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	Version        string `json:"version"`
}

func (UserController) Create(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	userModel := new(models.UserModel)

	if Data.Version == "trial" {
		user, err := userModel.CreateUserTrial(Data.Mail)

		userDetail, err := userModel.AddUserDetail(user.Id.Hex(), Data.Company, Data.FirstName, Data.LastName, Data.Phone)

		activateModel := new(models.ActivateKeyModel)
		activateKeyDetail, err := activateModel.CreateActivateKey(user.Id.Hex())
		if err != nil {
			_logger.Error(fmt.Sprintf("%v", err), false)
			return c.JSON(map[string]string{
				"code":    "unknown-error",
				"message": "Something went wrong. Please try again.",
			})
		}

		mail := new(models.SendMail)
		fmt.Println("Create user trial. sending mail.")
		mail.SendActivateKeyTrialVersion(userDetail.Name, activateKeyDetail.ActivateKey, user.Mail)

		return c.JSON(activateKeyDetail.ActivateKey)
	}

	if Data.MaxLineAccount <= 1 {
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	user, err := userModel.CreateUser(Data.Mail, Data.ExpirationDate, Data.MaxLineAccount)

	userModel.AddUserDetail(user.Id.Hex(), Data.Company, Data.FirstName, Data.LastName, Data.Phone)

	activateModel := new(models.ActivateKeyModel)
	activateModel.CreateActivateKey(user.Id.Hex())

	aggregationModel := new(models.AggregationData)
	dataSummary := aggregationModel.CreateSummaryUser(user.Id.Hex())

	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	return c.JSON(dataSummary)
}

func (UserController) Verify(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	userModel := new(models.UserModel)
	activateKey, err := userModel.VerifyActivateKey(Data.ActivateKey)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	} else if activateKey.Status == true && activateKey.Install == false {
		verifyFreshdesk, _ := userModel.ValidateFreshdeskDomain(Data.SubDomain, Data.ApiKey)
		if verifyFreshdesk == true {
			userDetail, errfind, errsave := userModel.AddFreshdeskDetail(activateKey.UserID, Data.SubDomain, Data.ApiKey)
			if errfind != nil || errsave != nil {
				if errfind != nil {
					_logger.Error(fmt.Sprintf("%v", errfind), false)
				} else {
					_logger.Error(fmt.Sprintf("%v", errsave), false)
				}
				return c.JSON(map[string]string{
					"code":    "unknown-error",
					"message": "Something went wrong. Please try again.",
				})
			} else {
				return c.JSON(map[string]interface{}{
					"install":      false,
					"freshdesk_id": userDetail.FreshdeskID,
					"activate_key": true,
					"domain":       true,
					"message":      "pass validation",
				})
			}
		} else {
			return c.JSON(map[string]interface{}{
				"install":      false,
				"activate_key": true,
				"domain":       false,
				"message":      "Domain or api key is not valid",
			})
		}
	} else if activateKey.Status == true && activateKey.Install == true {
		return c.JSON(map[string]interface{}{
			"install":      true,
			"activate_key": true,
			"domain":       false,
			"message":      "The activate key has been activated with another Freshdesk.",
		})
	} else if activateKey.Status == false {
		return c.JSON(map[string]interface{}{
			"install":      false,
			"activate_key": false,
			"domain":       false,
			"message":      "User non Activate",
		})
	}

	return nil
}

func (UserController) Delete(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	userModel := new(models.UserModel)
	if Data.UserID == "" {
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	err := userModel.DeleteUser(Data.UserID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	return c.JSON(map[string]string{
		"message": "Successfully deleted.",
	})
}

func (UserController) Regenerate(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}

	userModel := new(models.UserModel)
	err := userModel.RegenerateActivateKey(Data.UserID)
	aggregationModel := new(models.AggregationData)
	dataSummary := aggregationModel.CreateSummaryUser(Data.UserID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	return c.JSON(dataSummary)
}

func (UserController) ChangeMaxLineAccount(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userModel := new(models.UserModel)
	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}
	activateKey, err := userModel.VerifyActivateKey(Data.ActivateKey)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "Can not find Activate Key. Please try again.",
		})
	} else {
		if Data.MaxLineAccount >= 0 {
			user, err := userModel.ChangeMaxLINEAccount(activateKey.UserID, Data.MaxLineAccount)
			if user == nil {
				_logger.Error(fmt.Sprintf("%v", err), false)
				return c.JSON(map[string]string{
					"code":    "unknown-error",
					"message": "MaxLineAccount is not correct",
				})
			} else if err != nil {
				_logger.Error(fmt.Sprintf("%v", err), false)
				return c.JSON(map[string]string{
					"code":    "unknown-error",
					"message": "Something went wrong. Please try again.",
				})
			}
			return c.JSON(user)
		} else {
			return c.JSON(map[string]string{
				"code":    "unknown-error",
				"message": "Can not find Activate Key. Please try again.",
			})
		}
	}
}

func (UserController) ChangeStatus(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userModel := new(models.UserModel)
	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}
	user, err := userModel.ChangeStatusKeygen(Data.ActivateKey, Data.Status)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
		return c.JSON(map[string]interface{}{
			"Result":  false,
			"message": "Can't Verify active Key. Please try again.",
		})
	} else if user.Status == true {
		return c.JSON(map[string]interface{}{
			"Result":  user.Status,
			"message": "",
		})
	} else {
		return c.JSON(map[string]interface{}{
			"Result":  user.Status,
			"message": "User non Activate",
		})
	}
}

func (UserController) ChangeExpirationActivateKey(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userModel := new(models.UserModel)
	Data := new(UserController)
	if err := c.BodyParser(Data); err != nil {
		return err
	}
	user, errUserDB, errActivateKeyDB := userModel.ChangeExpiredDate(Data.UserID, Data.ExpirationDate)
	if errUserDB != nil || errActivateKeyDB != nil {
		if errUserDB != nil {
			_logger.Error(fmt.Sprintf("%v", errUserDB), false)
			return c.JSON(map[string]interface{}{
				"Result":  false,
				"message": "Can't Verify active Key. Please try again.",
			})
		} else {
			_logger.Error(fmt.Sprintf("%v", errActivateKeyDB), false)
			return c.JSON(map[string]interface{}{
				"Result":  false,
				"message": "Can't Verify active Key. Please try again.",
			})
		}
	} else {
		return c.JSON(map[string]interface{}{
			"Your expiration date is": user.ExpirationDate.Format("2006-01-02 15:04:05"),
			"message":                 "Save data successfully.",
		})
	}
}

func (UserController) CheckExpirationActivateKey(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userModel := new(models.UserModel)

	UserID := ""

	if c.Locals("userID") != nil {
		UserID = fmt.Sprintf("%v", c.Locals("userID"))
	} else if UserID == "" {
		Data := new(UserController)
		if err := c.BodyParser(Data); err != nil {
			return err
		}
		UserID = Data.UserID
	}

	expirationDate, errUserDB := userModel.CheckExpiredDate(UserID)
	if errUserDB != nil {
		_logger.Error(fmt.Sprintf("%v", errUserDB), false)
		return c.JSON(map[string]interface{}{
			"Result":  false,
			"message": "Can't Verify active Key. Please try again.",
		})
	} else {
		return c.JSON(map[string]interface{}{
			"version":        "trial",
			"expirationDate": expirationDate,
		})
	}
}
