package mode

import (
	"gopkg.in/mgo.v2/bson"
)


//一局德州游戏是一条数据
type T_th struct {
	Mid      bson.ObjectId		`json:"mid" bson:"_id"`
	Id       uint32   	//id
	DeskId   int32		//桌子号
	UserId	 uint32		//用户Id
	BetAmount	int64	//押注了多少钱
	WinAmount	int64	//赢了多少钱
	Blance		int64	//账户余额

}