package agentProModel

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_redpack/model/weixinModel"
	"math"
	"casino_redpack/model/userModel"
	"fmt"
	"casino_redpack/handler/weixin"
	"casino_redpack/handler/redpack"
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

func GetWithdrawalsList(query bson.M,page int) bson.M {

	Data := []*weixinModel.Withdrawals{}
	_,count := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Page(query,&Data, "endtime", page, 20)

	list := bson.M{
		"count":      count,
		"list_count": len(Data),
		"limit":      20,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(20)),
		"data":[]bson.M{},
	}
	listData := []bson.M{}
	for i,item := range Data {
		if i == count {
			continue
		}

		user := new(userModel.User)
		err := db.C(userModel.USER_TABLE_NAME).Find(bson.M{"id":item.UserId,},user)
		if err != nil {
			fmt.Println("null",err)
		}
		NewList := bson.M{
			"id":item.ObjId.Hex(),
			"username":userModel.GetUsernicknameById(item.UserId),
			"Time": item.Time.Format("2006-01-02 15:04:05"),
			"Number":item.Number,
			"Status":item.Status,
			"nikename":user.AccountNumber,
			"banck":userBanck(item.UserId),
		}
		listData = append(listData,NewList)
	}
	list["data"] = listData
	//ctx.Data["user"] = UserList
	return list

}

func userBanck(uid uint32) []bson.M {
	data := redpack.GetUserBanck(uid)
	list := []bson.M{}
	for _,item := range data {
		row := bson.M{
			"acct_no": item.AcctNo,
			"acct_name": item.AcctName,
			"acct_bank_name": item.AcctBankName,
			"rec_bank_name": item.RecBankName,
		}
		list = append(list,row)
	}
	return list
}

func GetOrderLists(query bson.M,page int)  bson.M {
	list := bson.M{
		"code": 0,
		"message": "fail",
		"request": bson.M{},
	}
	Data := []*weixin.RechargeOrder{}
	_,count := db.C(tableName.TABLE_ORDER_LISTS).Page(query,&Data, "-ordertime", page, 20)
	data := []bson.M{}
	for _,item := range Data {
		row := bson.M{
			"id":item.ObjId.Hex(),
			"time":item.OrderTime.Format("2006-01-02 15:04:05"),
			"userid":item.UserId,
			"nickname":userModel.GetUsernicknameById(item.UserId),
			"type":item.OrderType,
			"number":item.GoodsNunber,
			"status":item.OrderStatus,
			"googs":"金币",
		}
		data = append(data,row)
	}
	if count != 0{
		list["code"] = 1
		list["message"] = "success"
		list["count"] = count
		list["request"] = data
	}
	return list
}

