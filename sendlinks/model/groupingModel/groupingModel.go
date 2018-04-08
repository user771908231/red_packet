package groupingModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
)
type Grouping struct {
	ObjId bson.ObjectId		`bson:"_id"`
	GroupName string
	Time time.Time
}

func (G Grouping) Insert() error {
	G.ObjId = bson.NewObjectId()
	G.Time = time.Now()
	err := db.C(tableName.DB_GROUPING_LISTS).Insert(G)
	return err
}
