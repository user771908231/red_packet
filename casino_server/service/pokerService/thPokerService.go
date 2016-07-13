package pokerService

import (
	"casino_server/utils"
	"casino_server/msg/bbprotogo"
)

//德州扑克的纸牌

//config

//返回牌的列表
func RandomTHPorkCards(total int ) []*bbproto.Pai{
	result := make([]*bbproto.Pai,total)	//返回值
	indexs := RandomTHPorkIndex(0,52,total)
	for i := 0;i< total;i++ {
		result[i] =  bbproto.CreatePorkByIndex(indexs[i])
	}
	return result
}

//随机的德州牌的坐标
func RandomTHPorkIndex(min, max,total int) []int32 {
	result := make([]int32,total);
	count := 0;
	for count < total {
		num := utils.Rand(int32(min),int32(max))
		flag := true;
		for j := 0; j < total; j++ {
			if num == result[j] {
				flag = false;
				break;
			}
		}
		if flag {
			result[count] = num;
			count++;
		}
	}
	return result;
}