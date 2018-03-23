package pay

import (
	"casino_redpack/modules"
	"casino_redpack/handler/weixin"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

func  PostRechargeHandler(ctx *modules.Context)  {
	list := bson.M{
		"code": 0,
		"message": "faild",
		"request": bson.M{},
	}
	valtype := ctx.QueryInt("totalFee")
	data := weixin.ReturnData(float64(valtype),ctx.IsLogin().Id)
	if data != nil {
		list["code"] = 1
		list["request"] = data
		list["message"] = "success"
	}
	lists,_ := json.Marshal(list)
	ctx.Write([]byte(lists))

}
//金币充值
//func GoldRechargeHandler(ctx *modules.Context)  {
//	list := bson.M{
//		"code": 0,
//		"message": "success",
//		"request": bson.M{},
//	}
//	val := ctx.QueryInt("totalFee")
//	switch val {
//	case 20:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println(value)
//	case 50:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("50")
//	case 100:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("100")
//	case 200:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("200")
//	case 500:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("500")
//	case 1000:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("1000")
//	case 2000:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("2000")
//	default:
//		lists,_ := json.Marshal(list)
//		ctx.Write([]byte(lists))
//	}
//}
//
//func DData (ctx *modules.Context)  {
//	list := bson.M{
//		"code": 0,
//		"message": "success",
//		"request": bson.M{},
//	}
//	val := ctx.QueryInt("totalFee")
//	switch val {
//	case 20:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println(value)
//	case 50:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("50")
//	case 100:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("100")
//	case 200:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("200")
//	case 500:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("500")
//	case 1000:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("1000")
//	case 2000:
//		err,value := weixin.PayWapMoney(float64(0.01),ctx)
//		if err != nil{
//			list["code"] = 1
//			list["request"] = value
//			data,_ := json.Marshal(list)
//			ctx.Write([]byte(data))
//		}
//		fmt.Println("2000")
//	default:
//		lists,_ := json.Marshal(list)
//		ctx.Write([]byte(lists))
//	}
//}


