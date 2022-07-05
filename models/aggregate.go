package models

import (
	"fli-gateway-api/lib"

	"gopkg.in/mgo.v2/bson"
)

type AggregationData struct{}
type SummaryUser struct {
	User        []User        `json:"user"`
	Detail      []UserDetail  `json:"detail"`
	ActivateKey []ActivateKey `json:"activatekey"`
}

func (AggregationData) CreateSummaryUser(userID string) *SummaryUser {
	connection := lib.DBConnection()

	activateKey := &ActivateKey{}
	connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": userID}, activateKey)

	user := &User{}
	connection.Collection(USER_DB).FindById(bson.ObjectIdHex(userID), user)

	userDeatail := &UserDetail{}
	connection.Collection(USER_DETAIL_DB).FindOne(bson.M{"user_id": userID}, userDeatail)

	detail := &SummaryUser{
		User: []User{
			*user,
		},
		Detail: []UserDetail{
			*userDeatail,
		},
		ActivateKey: []ActivateKey{
			*activateKey,
		},
	}

	return detail
}
