package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fli-gateway-api/lib"
	"fmt"
	"strings"
	"time"

	"github.com/go-bongo/bongo"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

const USER_DB string = "USER"
const ActivateKEY_DB string = "ACTIVATE_KEY"
const LINE_OA_DB string = "LINE_OA"

type ActivateKey struct {
	bongo.DocumentBase `bson:",inline"`
	UserID             string `bson:"user_id" json:"user_id"`
	ActivateKey        string `bson:"activate_key" json:"activate_key"`
	Status             bool   `bson:"status" json:"status"`
}

type User struct {
	bongo.DocumentBase `bson:",inline"`
	MaxLineAccount     int       `bson:"max_line_account" json:"max_line_account"`
	ExpirationDate     time.Time `bson:"expiration_date" json:"expiration_date" `
}

type LINEDetailMiddleware struct {
	ClientSecret string `bson:"client_secret" json:"client_secret"`
}

type VerifySignaturesLINE struct {
	Content_Type  string `json:"content_type"`
	Authorization string `json:"authorization"`
}

func VerifySignatureAdmin(c *fiber.Ctx) error {
	if c.Get("x-signature") == "fGPYMIyH3JO8W86vZo" {
		return c.Next()
	} else {
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	}
}

func VerifySignatureCreateTrial(c *fiber.Ctx) error {
	if c.Get("x-signature") == "w3vv1dmM0aYcN7KNW2ep" {
		return c.Next()
	} else {
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	}
}

func VerifySignatureActivateKey(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	if c.Get("x-signature") == "fGPYMIyH3JO8W86vZo" {
		return c.Next()
	}

	decoded, err := base64.StdEncoding.DecodeString(c.Get("x-signature"))
	if err != nil {
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	}
	split_signature := strings.Split(string(decoded), ".")

	connection := lib.DBConnection()
	user_detail := &User{}
	erruser_deatil := connection.Collection(USER_DB).FindOne(bson.M{"freshdesk_id": split_signature[0]}, user_detail)
	if erruser_deatil != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", erruser_deatil))
		return c.Status(401).JSON(map[string]interface{}{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	activatekey := &ActivateKey{}
	errdetail_key := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": user_detail.DocumentBase.Id.Hex()}, activatekey)
	if errdetail_key != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", errdetail_key))
		return c.Status(401).JSON(map[string]interface{}{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}

	expirationDate := (user_detail.ExpirationDate.Sub(time.Now()).Hours() / 24) + 1
	if expirationDate > 0 {
		if !activatekey.Status {
			activatekey.Status = true
			connection.Collection(ActivateKEY_DB).Save(activatekey)
		}
		if !validateSignatureLINE(activatekey.ActivateKey, split_signature[1], c.Body()) {
			return c.Status(401).JSON(map[string]string{
				"Validate Signature": "false",
			})
		} else {
			c.Locals("userID", user_detail.DocumentBase.Id.Hex())
			c.Locals("maxLineAccount", user_detail.MaxLineAccount)
			return c.Next()
		}
	} else {
		activatekey.Status = false
		connection.Collection(ActivateKEY_DB).Save(activatekey)
		return c.Status(402).JSON(map[string]string{
			"Expiration Date": "false",
		})
	}
}

func VerifySignatureListLINE_OA(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	if c.Get("x-signature") == "fGPYMIyH3JO8W86vZo" {
		return c.Next()
	}

	decoded, err := base64.StdEncoding.DecodeString(c.Get("x-signature"))
	if err != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", err))
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	}
	split_signature := strings.Split(string(decoded), ".")

	connection := lib.DBConnection()
	user_detail := &User{}
	erruser_deatil := connection.Collection(USER_DB).FindOne(bson.M{"freshdesk_id": split_signature[0]}, user_detail)
	if erruser_deatil != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", erruser_deatil))
		return c.Status(401).JSON(map[string]interface{}{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	activatekey := &ActivateKey{}
	errdetail_key := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": user_detail.DocumentBase.Id.Hex()}, activatekey)
	expirationDate := user_detail.ExpirationDate.Sub(time.Now()).Hours() / 24
	if expirationDate <= 0 {
		activatekey.Status = false
		connection.Collection(ActivateKEY_DB).Save(activatekey)
	}
	if errdetail_key != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", errdetail_key))
		return c.Status(401).JSON(map[string]interface{}{
			"code":    "unknown-error",
			"message": "Something went wrong. Please try again.",
		})
	}
	if !validateSignatureLINE(activatekey.ActivateKey, split_signature[1], c.Body()) {
		_logger.InfoMiddleware(fmt.Sprintf("%v", "Validate Signature: false"))
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	} else {
		c.Locals("userID", user_detail.DocumentBase.Id.Hex())
		c.Locals("maxLineAccount", user_detail.MaxLineAccount)
		return c.Next()
	}

}

func VerifyActivateKey(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	if c.Get("x-signature") == "fGPYMIyH3JO8W86vZo" {
		return c.Next()
	}

	fmt.Println(string(c.Body()))

	connection := lib.DBConnection()
	activateKey := &ActivateKey{}
	errkey := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"activate_key": c.Get("ActivateKey")}, activateKey)
	if errkey != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", errkey))
		return c.JSON(map[string]interface{}{
			"ActivateKey": false,
			"code":        "unknown-error",
			"message":     "This activate key isn't valid",
		})
	}

	UserDetail := &User{}
	erruser := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(activateKey.UserID), UserDetail)
	if erruser != nil {
		_logger.InfoMiddleware(fmt.Sprintf("%v", erruser))
		return c.JSON(map[string]interface{}{
			"ActivateKey": false,
			"code":        "unknown-error",
			"message":     "Something went wrong. Please try again.",
		})
	} else if activateKey.Status == false {
		return c.JSON(map[string]interface{}{
			"ActivateKey": false,
			"code":        "unknown-error",
			"message":     "This activate key unavailable.",
		})
	} else if activateKey.Status == true {
		return c.Next()
	}
	return c.Status(401).JSON(map[string]string{
		"Validate Signature": "false",
	})
}

func VerifyingSignaturesLINE(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	connection := lib.DBConnection()
	Detail_Line := &LINEDetailMiddleware{}
	errDB := connection.Collection(LINE_OA_DB).FindOne(bson.M{"client_id": c.Params("*")}, Detail_Line)
	if errDB != nil {
		_logger.InfoLine(fmt.Sprintf("%v", errDB))
		return c.JSON(map[string]string{
			"code":    "unknown-error",
			"message": "This activate key isn't valid",
		})
	}

	if !validateSignatureLINE(Detail_Line.ClientSecret, c.Get("X-Line-Signature"), c.Body()) {
		return c.Status(401).JSON(map[string]string{
			"Validate Signature": "false",
		})
	} else {
		return c.Next()
	}
}

func validateSignatureLINE(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}
	return hmac.Equal(decoded, hash.Sum(nil))
}
