package Control

import (
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"errors"
)

type Control struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32
	Status Status //0雷  1红包雷
}

type Status struct {
	Open int
	Send int
}

func (C *Control) Isert() error {
	C.Id = bson.NewObjectId()
	err := db.C(tableName.TABLE_REDPACK_CONTROL).Update(bson.M{"_id":C.Id},C)
	if err != nil {
		return errors.New("插入一条记录失败！")
	}
	return nil
}
func (C *Control) Del() error {
	err := db.C(tableName.TABLE_REDPACK_CONTROL).Remove(bson.M{"userid":C.Id})
	if err != nil {
		return errors.New("插入一条记录失败！")
	}
	return nil
}
func (C *Control) Update() error {
	err := db.C(tableName.TABLE_REDPACK_CONTROL).Update(bson.M{"_id":C.Id},C)
	if err != nil {
		return errors.New("插入一条记录失败！")
	}
	return nil
}

func GetControlDel(id bson.ObjectId) error {
	err := db.C(tableName.TABLE_REDPACK_CONTROL).Remove(bson.M{"_id":id})
	return err

}

func GetControlAll() []*Control {
	lisr := []*Control{}
	err := db.C(tableName.TABLE_REDPACK_CONTROL).FindAll(bson.M{},&lisr)
	if err != nil {
		return nil
	}
	return lisr
}

func GetFindById( id uint32) *Control {
	lisr := new(Control)
	err := db.C(tableName.TABLE_REDPACK_CONTROL).Find(bson.M{"userid":id},lisr)
	if err != nil {
		return nil
	}
	return lisr
}
