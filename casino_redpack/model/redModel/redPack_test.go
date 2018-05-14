package redModel

import (
	"testing"

	"math/rand"
	"casino_common/common/log"
	"fmt"
	"casino_redpack/model/userModel"
	"time"
)

//func TestGetOpenRedMoney(t *testing.T) {
//	var ren int = 5
//	var num int = 5
//	var money int = 1000
//	var zong float64
//	for i := 0;i < ren;i++{
//		val := getOpenRedMoney(money,num)
//		money = money - val
//		t.Log("还剩余：",float64(money)/100)
//		t.Log("money:",float64(val)/100)
//		zong = zong + float64(val)/100
//		num --
//	}
//	t.Log("zong",zong)
//}

////拆红包算法(剩余的钱、剩余的人)
//func getOpenRedMoney(lost_money int, lost_person int,) int {
//	//参数合法性验证
//	if lost_money < 1 || lost_person <= 0 {
//		return 0
//	}
//
//	//只有一个人，则把钱全给他
//	if lost_person == 1 {
//		 return lost_money
//	}
//
//	//取0.01 - 平均金额 * 2 的值
//	lost_score := lost_money
//	avg_score := lost_score / lost_person
//	res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
//	if res_score == 0 {
//		res_score = 1
//	}
//
//	res_money := res_score
//
//	return res_money
//}

func TestGetRedPacketRecord(t *testing.T) {
	//拆红包算法(剩余的钱、剩余的人)
	var lost_money int = 100
	var lost_person int =7
	var id uint32 = 10117
	var L int = 2
	var u = &userModel.User{
		Id:10118,
	}
		//参数合法性验证
	if lost_money < 1 || lost_person <= 0 {
		t.Log("NOT FUNT")
	}

		//只有一个人，则把钱全给他
	if lost_person == 1 {
		t.Log("onre pepoer")
	}
	//取0.01 - 平均金额 * 2 的值
		lost_score := lost_money
		avg_score := lost_score / lost_person
		res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
		weishu := GetWeishu(float64(res_score)/100)
		res_money := res_score
	if u.Id == 10117 {

	 if id == 10117 {
		 log.T("米啊嘛%d", u.Id)
		 log.T("尾数:%d", weishu)
		 if weishu == L {
			 fmt.Println("开红包算法------------------------------------------------------%d", u.Id)
			 val := getOpenRedMoney(lost_money, lost_person, id, L, u)
			 t.Log(val)
		 }
	 }
		log.T("米啊嘛%d",u.Id)
		log.T("尾数:%d",weishu)
		if weishu != L {
			fmt.Println("开红包算法------------------------------------------------------%d",u.Id)
			val := getOpenRedMoney(lost_money, lost_person,id ,L ,u )
			t.Log(val)
		}

	}else{
		if res_score == 0 {
		res_score = 1
		}
		t.Log(res_money)
	}

}


func TestGetRedPacketRecord2(t *testing.T) {
	//拆红包算法(剩余的钱、剩余的人)
	var lost_money int = 100
	var lost_person int =7
	var id uint32 = 10115
	var L int = 2
	var u = &userModel.User{
		Id:10117,
	}
	for i := lost_person;i> 0;i-- { 

		//参数合法性验证
		if lost_money < 1 || lost_person <= 0 {
			t.Log("NOT FUNT")
		}

		//只有一个人，则把钱全给他
		if lost_person == 1 {
			t.Log("onre pepoer")
		}
		//取0.01 - 平均金额 * 2 的值
		lost_score := lost_money
		avg_score := lost_score / lost_person
		res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score * 2)
		weishu := GetWeishu(float64(res_score) / 100)
		res_money := res_score
		if u.Id == 10117 {

			if id == 10117 {
				log.T("米啊嘛%d", u.Id)
				log.T("尾数:%d", weishu)
				if weishu == L {
					fmt.Println("开红包算法------------------------------------------------------%d", u.Id)
					val := getOpenRedMoney(lost_money, lost_person, id, L, u)
					t.Log(val)
				}
			}
			log.T("米啊嘛%d", u.Id)
			log.T("尾数:%d", weishu)
			if weishu != L {
				fmt.Println("开红包算法------------------------------------------------------%d", u.Id)
				val := getOpenRedMoney(lost_money, lost_person, id, L, u)
				t.Log(val)
			}

		} else {
			if res_score == 0 {
				res_score = 1
			}
			t.Log(res_money)
		}
	}
}

