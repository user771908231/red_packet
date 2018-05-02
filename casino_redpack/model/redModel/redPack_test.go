package redModel

import (
	"testing"

)

func TestGetOpenRedMoney(t *testing.T) {
	var ren int = 5
	var num int = 5
	var money int = 1000
	var zong float64
	for i := 0;i < ren;i++{
		val := getOpenRedMoney(money,num)
		money = money - val
		t.Log("还剩余：",float64(money)/100)
		t.Log("money:",float64(val)/100)
		zong = zong + float64(val)/100
		num --
	}
	t.Log("zong",zong)
}

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


