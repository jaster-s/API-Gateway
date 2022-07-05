package models

import (
	"fli-gateway-api/lib"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

const CONFIG_DB string = "CONFIG"

type ConfigModel struct{}

type ConfigDetail struct {
	bongo.DocumentBase `bson:",inline"`
	Version            string `bson:"version" json:"version"`
	ExpirationDate     int    `bson:"expiration_date" json:"expiration_date"`
	MaxLINEAccount     int    `bson:"max_line_account" json:"max_line_account"`
}

func (ConfigModel) AddConfig(day, Account int) error {
	connection := lib.DBConnection()

	configDetail := &ConfigDetail{}
	err := connection.Collection(CONFIG_DB).FindOne(bson.M{"version": "trial"}, configDetail)
	if err != nil {
		configDetail = &ConfigDetail{
			Version:        "trial",
			ExpirationDate: day,
			MaxLINEAccount: Account,
		}
	} else {
		configDetail.ExpirationDate = day
		configDetail.MaxLINEAccount = Account
	}

	err = connection.Collection(CONFIG_DB).Save(configDetail)

	return err
}
