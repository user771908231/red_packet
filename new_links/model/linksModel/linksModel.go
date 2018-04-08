package linksModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
	"time"
)

type Links struct {
	ObjId 	bson.ObjectId		`bson:"_id"`
	Id 		uint32
	LinkName	string			//链接
	Groupings	bson.ObjectId	//分组
	Time time.Time
}

func (L Links) Insert() error {
	L.ObjId = bson.NewObjectId()
	L.Time = time.Now()
	err := db.C(tableName.DB_LINKS_LISTS).Insert(L)
	return err
}

func (L Links) Update() error{
	err :=db.C(tableName.DB_LINKS_LISTS).Update(bson.M{"_id":L.ObjId},L)
	return err
}

func LinksIdDel(Id bson.ObjectId) error {
	err := db.C(tableName.DB_LINKS_LISTS).Remove(bson.M{"_id":Id})
	return err
}

func GetLinksAll(query bson.M,page int,number int) (int,[]*Links){
	list := []*Links{}
	_,count := db.C(tableName.DB_LINKS_LISTS).Page(query, &list, "-requesttime", page, number)
	return count,list
}