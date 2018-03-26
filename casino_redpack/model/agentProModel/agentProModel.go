package agentProModel

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_redpack/model/weixinModel"
	"encoding/json"
)

//代理申请记录
type AgentProRecordRow struct {
	ObjId bson.ObjectId  //id
	Name string  //姓名
	Telphone string  //电话
	Comment string  //留言
	Wxid string  //微信id
	Ip string  //当前ip
	AddTime  time.Time  //添加时间
}

//插入记录
func (t *AgentProRecordRow) Insert() error {
	return db.C(tableName.DBT_APPLY_AGENTPRO_RECORD).Insert(t)
}

func GetWithdrawalsList(query bson.M,page int) []byte {
	list := bson.M{
		"code": 0,
		"message": "fail",
		"request": bson.M{},
	}
	Data := []*weixinModel.Withdrawals{}
	_,count := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Page(query,&Data, "-endtime", page, 20)
	if count != 0{
		list["code"] = 1
		list["message"] = "success"
		list["count"] = count
		list["request"] = Data
	}
	Datas,_ := json.Marshal(list)
	return Datas

}

func GetOrderLists(query bson.M,page int)  []byte {
	list := bson.M{
		"code": 0,
		"message": "fail",
		"request": bson.M{},
	}
	Data := []*weixinModel.Withdrawals{}
	_,count := db.C(tableName.TABLE_ORDER_LISTS).Page(query,&Data, "-endtime", page, 20)
	if count != 0{
		list["code"] = 1
		list["message"] = "success"
		list["count"] = count
		list["request"] = Data
	}
	Datas,_ := json.Marshal(list)
	return Datas
}