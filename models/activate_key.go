package models

import (
	"crypto/rand"
	"encoding/hex"
	"fli-gateway-api/lib"
	"fmt"
	"strings"
	"time"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

const ActivateKEY_DB string = "ACTIVATE_KEY"

type ActivateKeyModel struct{}

type ActivateKey struct {
	bongo.DocumentBase `bson:",inline"`
	UserID             string `bson:"user_id" json:"user_id"`
	ActivateKey        string `bson:"activate_key" json:"activate_key"`
	Status             bool   `bson:"status" json:"status"`
	Install            bool   `bson:"install" json:"install"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func splitBy(s string, n int) []string {
	var ss []string
	for i := 1; i < len(s); i++ {
		if i%n == 0 {
			ss = append(ss, s[:i])
			s = s[i:]
			i = 1
		}
	}
	ss = append(ss, s)
	return ss
}

func (ActivateKeyModel) GenerateActivateKey() string {

	DateToHex := fmt.Sprintf("%02x", time.Now().Unix())
	HexToSplit := splitBy(DateToHex, 4)

	key1, _ := randomHex(2)
	key2 := HexToSplit[0]
	key3, _ := randomHex(2)
	key4 := HexToSplit[1]

	ActivateKey := make([]string, 0)
	ActivateKey = append(ActivateKey, key1+"-")
	ActivateKey = append(ActivateKey, key2+"-")
	ActivateKey = append(ActivateKey, key3+"-")
	ActivateKey = append(ActivateKey, key4)

	Key := strings.ToUpper(strings.Join(ActivateKey, ""))

	return Key
}

func (ActivateKeyModel) CreateActivateKey(userID string) (*ActivateKey, error) {
	connection := lib.DBConnection()

	activateKey := new(ActivateKeyModel)
	keygen := activateKey.GenerateActivateKey()

	createActivateKey := &ActivateKey{
		UserID:      userID,
		ActivateKey: keygen,
		Status:      false,
		Install:     false,
	}

	userDetail := &User{}
	connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), userDetail)

	expiredDate := userDetail.ExpirationDate.Sub(time.Now()).Hours()/24 + 1
	if expiredDate > 0 {
		createActivateKey.Status = true
	}

	err := connection.Collection(ActivateKEY_DB).Save(createActivateKey)

	return createActivateKey, err
}

func (ActivateKeyModel) LoadDataActivateKeyformDB(clientID string) (*ActivateKey, error, error) {
	connection := lib.DBConnection()

	detail := new(LineModels)
	detail_lineOA, errDB := detail.GetLINE(clientID)
	if errDB != nil {
		return nil, errDB, nil
	}

	activateKey := &ActivateKey{}
	err := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": detail_lineOA.UserID}, activateKey)
	return activateKey, nil, err
}

func (ActivateKeyModel) GenerateUserID() string {

	key, _ := randomHex(12)

	return key
}
