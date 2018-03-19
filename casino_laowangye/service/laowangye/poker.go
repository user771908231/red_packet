package laowangye

import (
	"casino_common/proto/ddproto"
	"fmt"
	"errors"
)

//摇色子
func ParseShaiziType(list []int32) (shaiziType ddproto.LwyShaiziType, err error) {
	if len(list) != 4 {
		shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_LAN_DIAN
		err = errors.New("点数非法")
		return
	}
	//验证
	for _,num := range list {
		if num < 1 || num > 6 {
			shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_LAN_DIAN
			err = errors.New("点数非法")
			return
		}
	}

	//给色子排序
	for i:=0; i<3; i++ {
		for j:=i+1; j<4; j++ {
			if list[i] > list[j] {
				list[j], list[i] = list[i], list[j]
			}
		}
	}

	//色子类型
	switch fmt.Sprintf("%d%d%d%d", list[0], list[1], list[2], list[3]) {
	case "1234",  //顺子
		 "1222", "1333", "1444", "1555", "1666": //3个一样点数加点数1
		shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_7_DIAN
	case "3456",//顺子
		 "1111", "1112", "1113", "1114", "1115", "1116", //3个一样点数加点数2 3 4 5 6
		 "2222", "2223", "2224", "2225", "2226",
		 "3333", "2333", "3334", "3335", "3336",
		 "4444", "2444", "3444", "4445", "4446",
		 "5555", "2555", "3555", "4555", "5556",
		 "6666", "2666", "3666", "4666", "5666":
		shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_8_DIAN
	default:
		//最后是对子类型
		num_map := map[int32]int32{}
		for _,num := range list {
			if _,ok := num_map[num];ok {
				num_map[num]++
			}else {
				num_map[num] = 1
			}
		}

		if len(num_map) == 2 {
			//两对子如2255 3366，为8点
			shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_8_DIAN
			return
		}

		if len(num_map) == 4 {
			//四个不同的牌，则为烂点
			shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_LAN_DIAN
			return
		}

		//对子的情况
		var sum int32 = 0
		for _,num := range list {
			if num_item,ok := num_map[num];ok && num_item == 1 {
				sum += num
			}
		}
		if sum <= 7 {
			shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_7_DIAN
		}else {
			shaiziType = ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_8_DIAN
		}
	}

	return
}
