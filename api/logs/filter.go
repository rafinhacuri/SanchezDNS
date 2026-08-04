package logs

import (
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Filter(search string) bson.M {
	if search == "" {
		return bson.M{}
	}

	busca := bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"}

	return bson.M{"$or": []bson.M{
		{"username": busca},
		{"action": busca},
		{"details": busca},
		{"zone": busca},
	}}
}
