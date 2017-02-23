package agentModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"time"
)

//充值记录表
type RechargeLog struct {
	Id bson.ObjectId `bson:"_id"` //编号
	GoodsId int32  //商品Id
	GoodsNum int64  //商品数量
	AgentId uint32  //代理商id
	Amount float32  //总价格
	IsPay bool  //是否已支付
	AddTime time.Time   //充值时间
	DoneTime time.Time   //付款时间
}

//表更新
func (r *RechargeLog) Save() error {
	err := db.C(tableName.DBT_AGENT_RECHARGE_LOG).Update(bson.M{
		"_id": r.Id,
	},r)
	return err
}

//新增一个代理商充值订单
func AddNewRechargeLog(agent_id uint32,goods_id int32, goods_num int64) *RechargeLog {
	goods := GetGoodsInfoById(goods_id)
	if goods == nil {
		return nil
	}
	newOrder := &RechargeLog{
		Id:bson.NewObjectId(),
		GoodsId:goods_id,
		GoodsNum:goods_num,
		AgentId:agent_id,
		Amount:goods.Price * float32(goods_num),
		AddTime:time.Now(),
		DoneTime:time.Time{},
		IsPay:false,
	}
	err := db.C(tableName.DBT_AGENT_RECHARGE_LOG).Insert(newOrder)
	if err != nil {
		return nil
	}
	return newOrder
}
