package agentModel

import (
	"time"
	"casino_common/common/service/exchangeService"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
)

//代理商申请表
type ApplyRecord struct {
	Id bson.ObjectId `bson:"_id"`
	Name string  //真实姓名
	UserId uint32  //游戏内id
	InvitedId uint32  //推荐人id
	Phone string  //手机号
	Status exchangeService.ExchangeState  //审核状态
	RequestTime time.Time
	ProcessTime time.Time
}

//插入一张新记录
func (t *ApplyRecord) Insert() error {
	t.Id = bson.NewObjectId()
	t.RequestTime = time.Now()
	t.Status = exchangeService.PROCESS_ING
	return db.C(tableName.DBT_AGENT_APPLY_LOG).Insert(t)
}

//保存
func (r *ApplyRecord) Save() error {
	return db.C(tableName.DBT_AGENT_APPLY_LOG).Update(bson.M{
		"_id": r.Id,
	}, r)
}

//判断一个用户是否为代理商
func IsAgent(user_id uint32) bool {
	row := new(ApplyRecord)
	err := db.C(tableName.DBT_AGENT_APPLY_LOG).Find(bson.M{
		"userid": user_id,
		"status": exchangeService.PROCESS_TRUE,
	}, row)
	if err == nil {
		return true
	}
	return false
}

//通过id获取一条记录
func GetApplyRecordById(id string) *ApplyRecord {
	row := new(ApplyRecord)
	db.C(tableName.DBT_AGENT_APPLY_LOG).Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}, row)
	return row
}

