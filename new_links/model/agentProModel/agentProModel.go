package agentProModel

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_redpack/model/weixinModel"
	"casino_redpack/modules"
	"math"
	"casino_redpack/model/userModel"
	"fmt"
	"casino_redpack/handler/weixin"
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

func GetWithdrawalsList(ctx *modules.Context,query bson.M,page int) bson.M {

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
			"Id":item.ObjId.Hex(),
			"username":user.NickName,
			"Time": item.Time.Format("2006-01-02 15:04:05"),
			"Number":item.Number,
			"Status":item.Status,
		}
		listData = append(listData,NewList)
	}
	list["data"] = listData
	//ctx.Data["user"] = UserList
	return list

}

func GetOrderLists(query bson.M,page int)  bson.M {
	list := bson.M{
		"code": 0,
		"message": "fail",
		"request": bson.M{},
	}
	Data := []*weixin.RechargeOrder{}
	_,count := db.C(tableName.TABLE_NAME_OPEN_PACKET_LISTS).Page(query,&Data, "-endtime", page, 20)

	if count != 0{
		list["code"] = 1
		list["message"] = "success"
		list["count"] = count
		list["request"] = bson.M{}
	}
	return list
}

