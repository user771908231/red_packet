package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


//一局德州游戏是一条数据
type T_th_record struct {
	Mid      bson.ObjectId		`json:"mid" bson:"_id"`
	Id       int32   	//id
	DeskId   int32		//桌子号
	UserId	 uint32		//用户Id
	BetAmount	int64	//押注了多少钱
	WinAmount	int64	//赢了多少钱
	Blance		int64	//账户余额
	BeginTime	time.Time	//游戏开始时间
	GameNumber	int32	//游戏编号
}