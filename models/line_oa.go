package models

import (
	"encoding/json"
	"fli-gateway-api/lib"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

const LINE_OA_DB string = "LINE_OA"

type LineOAModel struct{}

type LineOA struct {
	bongo.DocumentBase `bson:",inline"`
	UserID             string `bson:"user_id" json:"user_id"`
	LINEname           string `bson:"line_name" json:"line_name"`
	ClientID           string `bson:"client_id" json:"client_id"`
	ClientSecret       string `bson:"client_secret" json:"client_secret"`
	ResponderID        int    `bson:"responder_id" json:"responder_id"`
	GroupID            int    `bson:"group_id" json:"group_id"`
	RecoveryMail       string `bson:"recovery_mail" json:"recovery_mail"`
	AutoMessageReply   string `bson:"auto_message_reply" json:"auto_message_reply"`
}

func (LineOAModel) List(userID string) (*[]LineOA, error) {
	connection := lib.DBConnection()

	lineOA := &[]LineOA{}

	result := connection.Collection(LINE_OA_DB).Find(bson.M{"user_id": userID})
	err := result.Query.All(lineOA)

	return lineOA, err
}

func (LineOAModel) AddLINEAccount(userID, name, clientID, clientSecret, mail, autoMessageReply string, groupID, responderID int) (*LineOA, error) {
	connection := lib.DBConnection()

	lineOA := &LineOA{}
	errfind := connection.Collection(LINE_OA_DB).FindOne(bson.M{"client_id": clientID}, lineOA)
	if errfind == nil {
		return nil, nil
	}

	createUser := &LineOA{
		UserID:           userID,
		LINEname:         name,
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		GroupID:          groupID,
		ResponderID:      responderID,
		RecoveryMail:     mail,
		AutoMessageReply: autoMessageReply,
	}
	err := connection.Collection(LINE_OA_DB).Save(createUser)
	return createUser, err
}

func (LineOAModel) DeleteLINEAccount(lineID string) error {
	connection := lib.DBConnection()

	err := connection.Collection(LINE_OA_DB).DeleteOne(bson.M{"_id": bson.ObjectIdHex(lineID)})
	return err
}

func (LineOAModel) EditAccount(lineID, mail, autoMessageReply string, groupID, responderID int) (*LineOA, string, error) {
	connection := lib.DBConnection()

	lineOA := &LineOA{}

	err := connection.Collection(LINE_OA_DB).FindById(bson.ObjectIdHex(lineID), lineOA)

	oldMail := lineOA.RecoveryMail
	if mail != "" && lineOA.RecoveryMail != mail {
		lineOA.RecoveryMail = mail
	}

	if autoMessageReply != "" && lineOA.AutoMessageReply != autoMessageReply {
		lineOA.AutoMessageReply = autoMessageReply
	}

	if groupID != 0 {
		lineOA.GroupID = groupID
	} else {
		lineOA.GroupID = 0
	}

	if responderID != 0 {
		lineOA.ResponderID = responderID
	} else {
		lineOA.ResponderID = 0
	}
	connection.Collection(LINE_OA_DB).Save(lineOA)
	return lineOA, oldMail, err
}

func (LineOAModel) Count(userID string) int {
	connection := lib.DBConnection()

	lineOA := &LineOA{}

	results := connection.Collection(LINE_OA_DB).Find(bson.M{"user_id": userID})
	count := 0

	for results.Next(lineOA) {
		count++
	}

	return count
}

func (LineOAModel) ValidateLINE(clientID, clientSecret string) (bool, []byte) {

	Url := "https://api.line.me/v2/oauth/accessToken"

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	encodedData := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", Url, encodedData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 400 {
		token_detail := AccessTokenLine{}
		json.Unmarshal(body, &token_detail)

		revokeAccessToken := new(LineModels)
		revokeAccessToken.RevokeAccessToken(token_detail.AccessToken)
		return true, body
	} else {
		return false, body
	}
}

func (LineOAModel) Webhook(lineID string) (string, error) {
	connection := lib.DBConnection()

	lineOA := &LineOA{}

	err := connection.Collection(LINE_OA_DB).FindById(bson.ObjectIdHex(lineID), lineOA)

	Url := os.Getenv("DOMAIN_WEBHOOK") + "/line_webhook/push/" + lineOA.ClientID

	return Url, err
}
