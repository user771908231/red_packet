package keysModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"

)

type Keys struct {
	ObjId 	bson.ObjectId	`bson:"_id"`
	Keys	string
	Remarks	string
	Status	int	//0 关 1 开
	Time 	time.Time
}

func (K *Keys)  Insert() error {
	K.ObjId = bson.NewObjectId()
	K.Status = 0
	K.Time = time.Now()
	err := db.C(tableName.DB_KEYS_LISTS).Insert(K)
	return err
}

func (K *Keys)  Update() error {
	err := db.C(tableName.DB_KEYS_LISTS).Update(bson.M{"_id":K.ObjId},K)
	return err
}

func (K *Keys)  Del() error {
	err := db.C(tableName.DB_KEYS_LISTS).Remove(bson.M{"_id":K.ObjId})
	return err
}
func GetKeysAll(query bson.M,page int,number int) (int,[]*Keys){
	list := []*Keys{}
	_,count := db.C(tableName.DB_KEYS_LISTS).Page(query, &list, "-requesttime", page, number)
	return count,list
}

func IdKeyRow(string string) *Keys {
	row := new(Keys)
	err := db.C(tableName.DB_KEYS_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},row)
	if err != nil {
		return nil
	}
	return row
}

func KeysStatus(string string,status int) error {
	row := GetkeysId(string)
	row.Status = status
	err :=row.Update()
	return err
}

func GetkeysId(string string) *Keys {
	row := new(Keys)
	err := db.C(tableName.DB_KEYS_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},row)
	if err != nil {
		return nil
	}
	return row
}

