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
	Remarks string
	Time time.Time
	Status int // 0 不启用 1 启用
}

func (G Grouping) Insert() error {
	G.ObjId = bson.NewObjectId()
	G.Time = time.Now()
	G.Status = 0
	err := db.C(tableName.DB_GROUPING_LISTS).Insert(G)
	return err
}

func (G Grouping) Update() error{
	err := db.C(tableName.DB_GROUPING_LISTS).Update(bson.M{"_id":G.ObjId},G)
	return err
}

func GetGroup(query bson.M) []*Grouping {
	list := []*Grouping{}
	err := db.C(tableName.DB_GROUPING_LISTS).FindAll(query,&list)
	if err != nil {
		return nil
	}
	return list
}

func DelGroup(string string) error {
	err := db.C(tableName.DB_GROUPING_LISTS).Remove(bson.M{"_id":bson.ObjectIdHex(string)})
	return err
}

func GroupStatus(string string,status int) error {
	row := GetGroupId(string)
	row.Status = status
	err :=row.Update()
	return err
}

func GetGroupId(string string) *Grouping {
	row := new(Grouping)
	err := db.C(tableName.DB_GROUPING_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},row)
	if err != nil {
		return nil
	}
	return row
}

func GetGroupObjId(string bson.ObjectId) *Grouping {
	row := new(Grouping)
	err := db.C(tableName.DB_GROUPING_LISTS).Find(bson.M{"_id":string},row)
	if err != nil {
		return nil
	}
	return row
}

func GetGroupHost(string string) *Grouping {
	row := new(Grouping)
	err := db.C(tableName.DB_GROUPING_LISTS).Find(bson.M{"groupname":string},row)
	if err != nil {
		return nil
	}
	return row
}