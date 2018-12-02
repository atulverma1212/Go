package dao

import (
	"gopkg.in/mgo.v2/bson"
)

func GetMap(distt string) (map [string] interface{}, error) {
	var data map[string] interface{}
	if err := db.C(CollectionMap).Find(bson.M{distt:bson.M{"$exists":true}}).One(&data); err!=nil {
		return nil, err
	}
	return data[distt].(map[string] interface{}), nil
}
