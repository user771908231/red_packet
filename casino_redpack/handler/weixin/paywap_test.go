package weixin

import (
	"testing"
	"casino_common/utils/numUtils"
	"fmt"
	"new_links/model/googsModel"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/userModel"
)

func TestCheckOrder(t *testing.T) {
	order := "2018041911274811000410004830"
	R := GetOrderId(order)
	if R == nil {
		t.Log("没找到订单")
		return
	}
	t.Log(R.OrderMoney)
	str := "200"

	if R.OrderMoney != numUtils.String2Float64(str) {
		t.Log("对应的的金额与支付金额不一致")
		return
	}
	t.Log("对应的的金额与支付金额一致")
	//判断是否是重复回调
	if R.OrderStatus == 1 {
		t.Log("tradeNo[%v]重复回调", order)
		return
	}

	t.Log("更新订单[%v]的回调信息，detail[%v]", order, R)
	//找到套餐
	G := googsModel.GetGoog(bson.ObjectIdHex(R.OrderGoods))
	User := userModel.GetUserById(R.UserId)
	//更新订单状态
	if User == nil {
		msg := fmt.Sprintf("没有在数据中找到用户ID：【%d】..", R.UserId)
		t.Log(msg)
		return
	}
	User.CapitalUplete("+",float64(G.Number),"充值")
	//更新订单状态
	R.OrderStatus = int64(1)
	err := R.Update()
	if err != nil {
		User.CapitalUplete("-",float64(G.Number),"")
		msg := fmt.Sprintf("更新订单tradeNo[%v]支付状态失败", order)
		t.Log(msg)
		return
	}
	t.Log("微信支付成功，为用户%d充值%d金币。", User.Id, int64(G.Number))

}
