package models

import (
	"fli-gateway-api/lib"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

const USER_DB string = "USER"
const USER_DETAIL_DB string = "USER_DETAIL"

type UserModel struct{}

type User struct {
	bongo.DocumentBase `bson:",inline"`
	Mail               string    `bson:"mail" json:"mail"`
	MaxLineAccount     int       `bson:"max_line_account" json:"max_line_account"`
	FreshdeskID        string    `bson:"freshdesk_id" json:"freshdesk_id"`
	SubDomain          string    `bson:"sub_domain" json:"sub_domain"`
	ApiKey             string    `bson:"api_key" json:"api_key"`
	ExpirationDate     time.Time `bson:"expiration_date" json:"expiration_date" `
}

type UserDetail struct {
	bongo.DocumentBase `bson:",inline"`
	UserID             string `bson:"user_id" json:"user_id"`
	Company            string `json:"company" json:"company"`
	Name               string `json:"name" json:"name"`
	Phone              string `json:"phone" json:"phone"`
}

func (UserModel) CreateUser(mail, expirationDATE string, maxLINEAccount int) (*User, error) {
	connection := lib.DBConnection()

	generateActivateKey := new(ActivateKeyModel)
	userDetail := &User{}
	freshdeskID := ""

	for {
		freshdeskID = generateActivateKey.GenerateUserID()
		connection.Collection(USER_DB).FindOne(bson.M{"freshdesk_id": freshdeskID}, userDetail)
		if userDetail.DocumentBase.Id == "" && freshdeskID != "" {
			break
		}
	}

	expirationDateSplit := strings.Split(expirationDATE, "/")
	year, _ := strconv.Atoi(string(expirationDateSplit[2]))
	month, _ := strconv.Atoi(string(expirationDateSplit[1]))
	day, _ := strconv.Atoi(string(expirationDateSplit[0]))

	createUser := &User{
		Mail:           mail,
		MaxLineAccount: maxLINEAccount,
		FreshdeskID:    freshdeskID,
		ExpirationDate: time.Now().AddDate(year, month, day),
	}

	err := connection.Collection(USER_DB).Save(createUser)
	return createUser, err
}

func (UserModel) CreateUserTrial(mail string) (*User, error) {
	connection := lib.DBConnection()

	generateActivateKey := new(ActivateKeyModel)
	userDetail := &User{}
	freshdeskID := ""

	for {
		freshdeskID = generateActivateKey.GenerateUserID()
		connection.Collection(USER_DB).FindOne(bson.M{"freshdesk_id": freshdeskID}, userDetail)
		if userDetail.DocumentBase.Id == "" && freshdeskID != "" {
			break
		}
	}
	configDetail := &ConfigDetail{}
	err := connection.Collection(CONFIG_DB).FindOne(bson.M{"version": "trial"}, configDetail)
	if err != nil {
		return nil, err
	}

	createUser := &User{
		Mail:           mail,
		MaxLineAccount: configDetail.MaxLINEAccount,
		FreshdeskID:    freshdeskID,
		ExpirationDate: time.Now().AddDate(0, 0, configDetail.ExpirationDate),
	}

	err = connection.Collection(USER_DB).Save(createUser)
	return createUser, err
}

func (UserModel) AddUserDetail(userID, company, firstName, lastName, phone string) (*UserDetail, error) {
	connection := lib.DBConnection()

	userDetail := &UserDetail{
		UserID:  userID,
		Company: company,
		Name:    firstName + " " + lastName,
		Phone:   phone,
	}

	errsave := connection.Collection(USER_DETAIL_DB).Save(userDetail)
	return userDetail, errsave
}

func (UserModel) AddFreshdeskDetail(userID, subDomain, apiKey string) (*User, error, error) {
	connection := lib.DBConnection()

	if subDomain == "" || apiKey == "" {
		return nil, nil, nil
	}

	user := &User{}
	errfind := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), user)

	user.SubDomain = subDomain
	user.ApiKey = apiKey

	errsave := connection.Collection(USER_DB).Save(user)
	return user, errfind, errsave
}

func (UserModel) VerifyActivateKey(key string) (*ActivateKey, error) {
	connection := lib.DBConnection()

	activateKey := &ActivateKey{}
	err := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"activate_key": key}, activateKey)

	return activateKey, err
}

func (UserModel) DeleteUser(userID string) error {
	connection := lib.DBConnection()

	err := connection.Collection(USER_DB).DeleteOne(bson.M{"_id": bson.ObjectIdHex(userID)})
	connection.Collection(USER_DETAIL_DB).DeleteOne(bson.M{"user_id": userID})
	connection.Collection(ActivateKEY_DB).DeleteOne(bson.M{"user_id": userID})
	connection.Collection(LINE_OA_DB).Delete(bson.M{"user_id": userID})
	return err
}

func (UserModel) RegenerateActivateKey(userID string) error {
	connection := lib.DBConnection()

	activateKey := &ActivateKey{}
	err := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": userID}, activateKey)

	activateKey.Status = false

	generateActivateKey := new(ActivateKeyModel)
	NewactivateKey := generateActivateKey.GenerateActivateKey()

	activateKey.ActivateKey = NewactivateKey

	connection.Collection(ActivateKEY_DB).Save(activateKey)
	return err
}

func (UserModel) ChangeStatusKeygen(key string, status bool) (*ActivateKey, error) {
	connection := lib.DBConnection()

	activateKey := &ActivateKey{}
	err := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"activate_key": key}, activateKey)

	activateKey.Status = status
	connection.Collection(ActivateKEY_DB).Save(activateKey)

	return activateKey, err
}

func (UserModel) ChangeMaxLINEAccount(userID string, maxLineAccount int) (*User, error) {
	connection := lib.DBConnection()

	user := &User{}
	err := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), user)

	user.MaxLineAccount = maxLineAccount

	connection.Collection(USER_DB).Save(user)
	return user, err
}

func (UserModel) ValidateFreshdeskDomain(subdomain, apikey string) (bool, []byte) {

	Url := subdomain + "api/v2/contacts"

	req, _ := http.NewRequest("GET", Url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 400 {
		return true, body
	} else {
		return false, body
	}

}

func (UserModel) ChangeExpiredDate(userID, expirationDate string) (*User, error, error) {
	connection := lib.DBConnection()

	expirationDateSplit := strings.Split(expirationDate, "/")
	year, _ := strconv.Atoi(string(expirationDateSplit[2]))
	month, _ := strconv.Atoi(string(expirationDateSplit[1]))
	day, _ := strconv.Atoi(string(expirationDateSplit[0]))

	user := &User{}
	errUserDB := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), user)

	activateKey := &ActivateKey{}
	errActivateKeyDB := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": userID}, activateKey)

	expiredDay := (user.ExpirationDate.Sub(time.Now()).Hours() / 24) + 1

	if expirationDate != "00/00/00" {
		if expiredDay > 0 {
			user.ExpirationDate = user.ExpirationDate.AddDate(year, month, day)
		} else {
			user.ExpirationDate = time.Now().AddDate(year, month, day)
		}
		if !activateKey.Status {
			activateKey.Status = true
		}
	}
	connection.Collection(USER_DB).Save(user)
	connection.Collection(ActivateKEY_DB).Save(activateKey)
	return user, errUserDB, errActivateKeyDB
}

func (UserModel) CheckExpiredDate(userID string) (int, error) {
	connection := lib.DBConnection()

	user := &User{}
	errUserDB := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), user)

	expiredDay := (user.ExpirationDate.Sub(time.Now()).Hours() / 24) + 1

	return int(expiredDay), errUserDB
}
