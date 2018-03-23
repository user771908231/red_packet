package redpack

import (
	"casino_redpack/modules"
)

//充值
func RechargeHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/recharge")
}
//
//func  RechargeConfirmHandler(ctx *modules.Context)  {
//	list := bson.M{
//		"code": 0,
//		"message": "faild",
//		"request": bson.M{},
//	}
//	valtype := ctx.QueryInt("totalFee")
//	data := weixin.ReturnData(float64(valtype),ctx.IsLogin().Id)
//	if data != nil {
//		list["code"] = 1
//		list["request"] = data
//		list["message"] = "success"
//	}
//	lists,_ := json.Marshal(list)
//	ctx.Write([]byte(lists))
//
//}

func PostRechargeHandler(ctx *modules.Context){

}