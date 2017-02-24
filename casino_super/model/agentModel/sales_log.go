package agentModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"time"
)

//销售记录表
type SalesLog struct {
	Id bson.ObjectId `bson:"_id"` //销售编号
	AgentId uint32  //代理商id
	UserId uint32  //用户id
	Type GoodsType //充值类型 金币 钻石 房卡 等
	Num int64  //数量
	Money float64  //金额
	Remark string  //备注
	AddTime time.Time  //销售时间
}

//新增代理商销售记录
func AddNewSalesLog(agent_id uint32, user_id uint32, goods_type GoodsType, num int64,money float64, remark string) (err error) {
	err = db.C(tableName.DBT_AGENT_SALES_LOG).Insert(SalesLog{
		Id: bson.NewObjectId(),
		AgentId:agent_id,
		UserId:user_id,
		Type:goods_type,
		Num: num,
		Money: money,
		Remark:remark,
		AddTime: time.Now(),
	})

	return err
}
