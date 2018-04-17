package adminModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
)

//权限
type Power struct {
	ObjIds 	bson.ObjectId 	`bson:"_id"`
	PowerName 	string		//权限名称
	Status 		int	// 0 关闭	1开启
	Time    	time.Time
}

func (P *Power) Insert() error {
	P.ObjIds = bson.NewObjectId()
	P.Status = 0
	P.Time = time.Now()
	err := db.C(tableName.DB_LINKS_POWER_INFO).Insert(P)
	return err
}

//删除
func DelPower(s string) error {
	err := db.C(tableName.DB_LINKS_POWER_INFO).Remove(bson.M{"_id":bson.ObjectIdHex(s)})
	return err
}

//修改
func (P *Power) Save() error {
	err := db.C(tableName.DB_LINKS_POWER_INFO).Update(bson.M{"_id": P.ObjIds}, P)
	return err
}