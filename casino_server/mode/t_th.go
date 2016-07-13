package mode

import (
	"gopkg.in/mgo.v2/bson"
)


//一局德州游戏是一条数据
type T_th struct {
	Mid      bson.ObjectId		`json:"mid" bson:"_id"`
	Id       uint32   	//id
	DeskId   uint32		//桌子号
}