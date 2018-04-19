package redpack

import (
	"testing"
	"casino_redpack/model/userModel"
)

func TestJudgeInMine(t *testing.T) {
	var open_tail_num int = 7
	var tailnumber int = 7
	var red_money float64 = 500
	var money float64 = 200.07
	var number int = 7
	var ThisUserID uint32 = 10023
	var SendPacketUserId uint32 = 10024

	var Odds float64
	switch number {
	case 7:
		Odds = 1.6
	case 8:
		Odds = 1.4
	case 9:
		Odds = 1.2
	case 10:
		Odds = 1
	default:
		Odds = 1
	}
	if open_tail_num == tailnumber {
		t.Log("中雷用户：%d",ThisUserID)
		this_user := userModel.GetUserById(ThisUserID)
		if this_user == nil {
			t.Log("没有找到中雷的用户")
			return
		}
		t.Log("原来的金币数：%d",money)
		money1 := money - FloatValue(money * 0.03, 2)
		t.Log("扣除的费率的到金币数:%d",money1)
		err := this_user.CapitalUplete("+",money1,"开红包")
		if err != nil {
			t.Log("开红包的用户加金币失败 err:%s",err)
			return
		}
		money0 := FloatValue(red_money * Odds,2)
		t.Log("中雷赔的金币数：%d",money0)

		//获取发包人的信息
		SendUser := userModel.GetUserById(SendPacketUserId)
		if SendUser == nil {
			t.Log("没有找到发红包的用户")
			return
		}
		err1 := SendUser.CapitalUplete("+",money0,"赢")
		if err1 != nil {
			t.Log("开红包的用户加金币失败 err:%s",err)
			return
		}
		t.Log("执行成功")
	}else {
		//--没中雷---
		t.Log("没中雷")
		//open_record_multiple = 0.03 //百分之三的扣费率
		this_user := userModel.GetUserById(ThisUserID)
		if this_user == nil {
			t.Log("没有找到此用户！")
			return
		}
		money0 := money - FloatValue(money * 0.03,2)
		err := this_user.CapitalUplete("+",money0,"开红包")
		//---end
		if err != nil {

			t.Log("➕操作错误")
			return
		}
		t.Log("执行成功")
		return
	}
}