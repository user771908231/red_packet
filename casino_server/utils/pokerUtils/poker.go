package pokerUtils

import (
	"strings"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	"casino_server/common/log"
)

var (
	FLOWER_DIAMOND int32 = 1        //方片
	FLOWER_CLUB int32 = 2        //美化
	FLOWER_HEART int32 = 3        //红桃
	FLOWER_SPADE int32 = 4        //黑桃
	FLOWER_REDJOKER int32 = 5        //王
	FLOWER_BLACKJOKER int32 = 6        //王
)

var PokerMap map[int32]string

func init() {
	PokerMap = make(map[int32]string, 54)
	PokerMap[0] = "POKER_diamond_14_A"
	PokerMap[1] = "POKER_diamond_2_2"
	PokerMap[2] = "POKER_diamond_3_3"
	PokerMap[3] = "POKER_diamond_4_4"
	PokerMap[4] = "POKER_diamond_5_5"
	PokerMap[5] = "POKER_diamond_6_6"
	PokerMap[6] = "POKER_diamond_7_7"
	PokerMap[7] = "POKER_diamond_8_8"
	PokerMap[8] = "POKER_diamond_9_9"
	PokerMap[9] = "POKER_diamond_10_10"
	PokerMap[10] = "POKER_diamond_11_J"
	PokerMap[11] = "POKER_diamond_12_Q"
	PokerMap[12] = "POKER_diamond_13_K"

	PokerMap[13] = "POKER_club_14_A"
	PokerMap[14] = "POKER_club_2_2"
	PokerMap[15] = "POKER_club_3_3"
	PokerMap[16] = "POKER_club_4_4"
	PokerMap[17] = "POKER_club_5_5"
	PokerMap[18] = "POKER_club_6_6"
	PokerMap[19] = "POKER_club_7_7"
	PokerMap[20] = "POKER_club_8_8"
	PokerMap[21] = "POKER_club_9_9"
	PokerMap[22] = "POKER_club_10_10"
	PokerMap[23] = "POKER_club_11_J	"
	PokerMap[24] = "POKER_club_12_Q"
	PokerMap[25] = "POKER_club_13_K"

	PokerMap[26] = "POKER_heart_14_A"
	PokerMap[27] = "POKER_heart_2_2"
	PokerMap[28] = "POKER_heart_3_3"
	PokerMap[29] = "POKER_heart_4_4"
	PokerMap[30] = "POKER_heart_5_5"
	PokerMap[31] = "POKER_heart_6_6"
	PokerMap[32] = "POKER_heart_7_7"
	PokerMap[33] = "POKER_heart_8_8"
	PokerMap[34] = "POKER_heart_9_9"
	PokerMap[35] = "POKER_heart_10_10"
	PokerMap[36] = "POKER_heart_11_J"
	PokerMap[37] = "POKER_heart_12_Q"
	PokerMap[38] = "POKER_heart_13_K"

	PokerMap[39] = "POKER_spade_14_A"
	PokerMap[40] = "POKER_spade_2_2"
	PokerMap[41] = "POKER_spade_3_3"
	PokerMap[42] = "POKER_spade_4_4"
	PokerMap[43] = "POKER_spade_5_5"
	PokerMap[44] = "POKER_spade_6_6"
	PokerMap[45] = "POKER_spade_7_7"
	PokerMap[46] = "POKER_spade_8_8"
	PokerMap[47] = "POKER_spade_9_9"
	PokerMap[48] = "POKER_spade_10_10"
	PokerMap[49] = "POKER_spade_11_J"
	PokerMap[50] = "POKER_spade_12_Q"
	PokerMap[51] = "POKER_spade_13_K"
	PokerMap[53] = "POKER_blackjoker_15_JOKER"
	PokerMap[54] = "POKER_redjoker_16_JOKER"
}


/**
MapKey           *int32
Mapdes           *string
Value            *int32
Flower           *string
Name             *string
 */

//返回值
func ParseByIndex(index int32) (int32, string, int32, string, string) {
	var rmapdes string = PokerMap[index]
	sarry := strings.Split(rmapdes, "_")
	var pvalue int32 = int32(numUtils.String2Int(sarry[2]))
	var pname string = sarry[3]
	var pflower string = sarry[1]
	return index, rmapdes, pvalue, pflower, pname
}

//随机一副扑克牌的index
func Xipai(paiCount int) []int32 {
	//初始化一个顺序的牌的集合
	pmap := make([]int32, paiCount)
	for i := 0; i < paiCount; i++ {
		pmap[i] = int32(i)
	}
	//打乱牌的集合
	pResult := make([]int32, paiCount)
	for i := 0; i < paiCount; i++ {
		rand := utils.Rand(int32(0), int32(paiCount - i))
		pResult[i] = pmap[rand]
		pmap = append(pmap[:rand], pmap[rand + 1:]...)
	}

	log.T("洗牌之后,得到的随机的index数组[%v]", pResult)
	return pResult
}
