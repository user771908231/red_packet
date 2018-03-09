package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"math/rand"
	"sync"
	"casino_redpack/model/userModel"
	"casino_common/utils/db"
)

//红包详情表
var TABLE_NAME_REDPACK_INFO = "t_redpack_info"

//单个红包
type Redpack struct {
	ObjId       bson.ObjectId `bson:"_id"`
	CreatorUser uint32            //发包的人id
	CreatorName string            //发包人昵称
	CreatorHead string            //发包人头像
	Id          int32             //红包id
	Money       float64           //红包大小
	Lost        float64           //剩余大小
	Piece       int             //分成几份
	Type        RoomType          //红包类型
	TailNumber  int             //尾数 0-9, 默认-1
	RoomId      int32             //房间id
	OpenRecord  []*OpenRecordItem //开红包记录
	Time        time.Time         //发红包时间
}

//开红包记录
type OpenRecordItem struct {
	UserId uint32  //领红包的人id
	NickName string  //领红包的人的昵称
	Head string  //领红包的人的头像
	Money float64  //领了多少钱
	Time time.Time  //时间
}

//红包锁
var redLock sync.Map

//锁
func Lock(redId int32) {
	lock,_ := redLock.LoadOrStore(redId, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
}

//解锁
func UnLock(redId int32) {
	lock,ok := redLock.Load(redId)
	if ok {
		lock.(*sync.Mutex).Unlock()
	}
}

//保存到mongo数据库
func (redInfo *Redpack) Upsert() error {
	return db.C(TABLE_NAME_REDPACK_INFO).Upsert(bson.M{
		"_id": redInfo.ObjId,
	}, redInfo)
}

//拆红包
func (redInfo *Redpack) Open(user *userModel.User) float64 {
	//加逻辑锁，保证线程安全
	Lock(redInfo.Id)
	defer UnLock(redInfo.Id)

	//是否还有余额
	if redInfo.Lost <= 0 || redInfo.Piece <= len(redInfo.OpenRecord) {
		return 0
	}

	//是否已经拆过红包了
	for _,item := range redInfo.OpenRecord {
		if item.UserId == user.Id {
			return item.Money
		}
	}

	//开始拆红包
	open_money := getOpenRedMoney(redInfo.Lost, redInfo.Piece - len(redInfo.OpenRecord))
	//更新红包余额
	redInfo.Lost -= open_money
	//更新开包记录
	redInfo.OpenRecord = append(redInfo.OpenRecord, &OpenRecordItem{
		UserId: user.Id,
		NickName: user.NickName,
		Head: user.HeadUrl,
		Money: open_money,
		Time: time.Now(),
	})

	//更新到mongodb
	redInfo.Upsert()

	return open_money
}

//拆红包算法(剩余的钱、剩余的人)
func getOpenRedMoney(lost_money float64, lost_person int) float64 {
	//参数合法性验证
	if lost_money < 0.1 || lost_person <= 0 {
		return 0
	}

	//只有一个人，则把钱全给他
	if lost_person == 1 {
		return float64(int(lost_money * 100))/100
	}

	//取0.01 - 平均金额 * 2 的值
	lost_score := int(lost_money * 100)
	avg_score := lost_score / lost_person
	res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
	if res_score == 0 {
		res_score = 1
	}

	return float64(res_score)/100
}
