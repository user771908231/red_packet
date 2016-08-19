package mode

import "gopkg.in/mgo.v2/bson"

type BaseMode interface {
	GetMid()	bson.ObjectId
}

