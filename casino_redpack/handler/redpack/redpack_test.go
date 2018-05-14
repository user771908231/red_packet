package redpack

import (
	"testing"
	"casino_redpack/model/userModel"
	"fmt"
	"strings"
	"unicode/utf8"


	"casino_redpack/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
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

func TestAgentRebate(t *testing.T) {
	var money float64 = 1000
	U := userModel.GetUserById(10024)
	err1 := U.CapitalUplete("+", FloatValue(money*0.3, 2), "用户返利")
	if err1 != nil {
		t.Log("用户的一级代理返利失败！ error：%s", err1)
		return
	}
	//判断代理是否有上级
	level_two := userModel.GetUserById(uint32(U.ExtensionId))
	if level_two != nil {
		//找到代理的上级
		err2 := level_two.CapitalUplete("+", FloatValue(money*0.07, 2), "下级代理返利")
		if err2 != nil {
			t.Log("代理的下一级代理返利失败！ error：%s", err2)
			return
		}
		level_three := userModel.GetUserById(uint32(level_two.ExtensionId))
		if level_three != nil {
			err3 := level_three.CapitalUplete("+", FloatValue(money*0.07, 2), "下级代理返利")
			if err3 != nil {
				t.Log("代理的下一级代理返利失败！ error：%s", err3)
				return
			}
			t.Log("成功退出！")
			return
		}
		t.Log("成功退出！")
		return
	}
	t.Log("成功退出！")
	return
}

func TestJoinWurenRedPacketHandler(t *testing.T) {
	var val float64 = 19.83
	t.Log("得到的值%f",val)
	s :=fmt.Sprintf("%.2f", val)
	t.Log("得到的值字符串",s)
	ss := strings.Split(s,".")
	t.Log("得到的值字符串去除小数点byte",ss)

	sss := strings.Join(ss,"")
	t.Log("得到的去除小数点字符串",sss)
	by := []byte(sss)

	for i,_ := range by {
		t.Log(string(by[i]))
	}

	slengt := utf8.RuneCountInString(sss)
	t.Log("得到的去除小数点字符串的长度",slengt)
}


func TestShouxi(t *testing.T) {
	log := agentModel.GetAgentRebateLogByIdList(10004,agentModel.TimeObject())
	list := []bson.M{}
	for _,item := range log {
		row := bson.M{
			"name":IsName(item.RebateId),
			"money":item.RebateMoeny,
		}
		list = append(list,row)
	}
	t.Log(list)

}

func TestAgentHandler(t *testing.T) {
	val := 123
	s :=fmt.Sprintf("%d",val)
	lengt := len(string(s))
	srt := strings.Replace(s,s[:lengt],"4",lengt)
	t.Log(srt)
}


