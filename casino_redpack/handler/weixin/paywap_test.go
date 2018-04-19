package weixin

import (
	"testing"
	"casino_common/utils/numUtils"
	"fmt"
	"new_links/model/googsModel"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/userModel"
	"regexp"
	"unicode/utf8"
	"strings"
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

func TestGenerateOtder(t *testing.T) {

	var D float64 = 9.99
	t.Log("得到的值%f",D)
	s :=fmt.Sprintf("%.2f", D)
	t.Log("得到的值字符串",s)
	ss := strings.Split(s,".")
	t.Log("得到的值字符串去除小数点byte",ss)
	sss := strings.Join(ss,"")
	t.Log("得到的去除小数点字符串",sss)
	slengt := utf8.RuneCountInString(sss)
	t.Log("得到的去除小数点字符串的长度",slengt)
	switch slengt {
	case 3:
		//顺子
		match1, _ := regexp.MatchString("^(?:(123)||(234)||(345)||(456)||(567)||(678)||(789)){3}$", sss)
		//豹子
		match2, _ := regexp.MatchString("^(?:(1)||(2)||(3)||(4)||(5)||(6)||(7)||(8)||(9)){3}$", sss)
		//520
		match7, _ := regexp.MatchString("^(?:(520){3})$", sss)
		//001
		match9, _ := regexp.MatchString("^(?:(001)){3}$", sss)
		if match1 {
			t.Log("3顺子",6.88)
			return
		}else if match2 {
			t.Log("3豹子",10)
			return
		}else if match7 {
			t.Log("520",18.88)
			return
		}else if match9 {
			t.Log("001",18.88)
			return
		}
		t.Log("什么都不是",sss)
		break;
	case 4:
		match3, _ := regexp.MatchString("^(?:(1234)||(2345)||(3456)||(4567)||(5678)||(6789)){4}$", sss)

		match4, _ := regexp.MatchString("^(?:(1)||(2)||(3)||(4)||(5)||(6)||(7)||(8)||(9)){4}$", sss)

		match8, _ := regexp.MatchString("^(?:(1314)){4}$", sss)

		match0, _ := regexp.MatchString("^(?:(2018)){4}$", sss)

		if match3 {
			t.Log("4顺子",38.88)
			return
		}else if match4 {
			t.Log("4豹子",58.88)
			return
		}else if match8 {
			t.Log("1314",18.88)
			return
		}else if match0 {
		t.Log("2018",18.88)
		return
		}
		t.Log("什么都不是",sss)
		break;
	case 5:
		match5, _ := regexp.MatchString("^(?:(12345)||(23456)||(34567)||(45678)||(56789)){5}$", sss)
		t.Log("判断4位顺子",match5)
		match6, _ := regexp.MatchString("^(?:(1)||(2)||(3)||(4)||(5)||(6)||(7)||(8)||(9)){5}$", sss)
		t.Log("判断4位豹子",match6)
		if match5 {
			t.Log("5顺子",388)
			return
		}else if match6 {
			t.Log("5豹子",588)
			return
		}
		t.Log("什么都不是",sss)
		break;
	//case 6:
	//	match5, _ := regexp.MatchString("^(?:[1-6]|[2-7]|[3-8]|[4-9]){6}$", "3456")
	//	t.Log("判断6位顺子",match5)
	//	match6, _ := regexp.MatchString("^(?:[1]|[2]|[3]|[4]|[5]|[6]|[7][8][9]){6}$", "222")
	//	t.Log("判断6位豹子",match6)
	//	if match5 {
	//		t.Log("6顺子",6.88)
	//		return
	//	}else if match6 {
	//		t.Log("6豹子",6.88)
	//		return
	//	}
	//	break;
	default:
		t.Log("什么都不是",sss)
		break;
	}

}


