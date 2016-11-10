package majiang

import (
	"strings"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	. "casino_majiang/msg/protogo"
	"casino_server/common/log"
	"errors"
)

//得到一副牌...


//**处理忙将牌相关的service

var mjpaiMap map[int]string
var clienMap map[int]int32

//const
var W int32 = 1        //万
var S int32 = 2        //条
var T int32 = 3        //筒
var MJPAI_COUNT int = 108        //牌的张数

//基本牌型番数
var FAN_PINGHU		int32	= 0 //平胡 0番
var FAN_DADUIZI		int32	= 1 //大对子 1番
var FAN_QINGYISE	int32	= 2 //平胡清一色 2番
var FAN_DAIYAOJIU	int32	= 2 //带幺九 2番
var FAN_QIDUI		int32	= 2 //七对 2番
var FAN_QINGDUI		int32	= 3 //清对 3番
var FAN_JIANGDUI	int32	= 3 //将对 3番
var FAN_LONGQIDUI	int32	= 4 //龙七对 4番
var FAN_QINGQIDUI	int32	= 4 //清七对 4番
var FAN_JIANGQIDUI	int32	= 4 //将七对 4番
var FAN_QINGYAOJIU	int32	= 4 //清幺九 4番
var FAN_TIAN_DI_HU	int32	= 5 //天地胡
var FAN_QINGLONGQIDUI	int32	= 5 //清龙七对

//加番
var FAN_ZIMO		int32	= 1 //自摸
var FAN_JINGOUDIAO	int32	= 1 //金钩钓
var FAN_MENQ_ZHONGZ	int32	= 1 //门清中张
var FAN_GANGSHANGHUA	int32	= 1 //杠上花
var FAN_GANGSHANGPAO	int32	= 1 //杠上炮
var FAN_HD_HUA		int32	= 1 //海底花
var FAN_HD_PAO		int32	= 1 //海底炮
var FAN_QIANGGANG	int32	= 1 //抢杠
var FAN_HD_GANGSHANGHUA	int32	= 1 //海底杠上花
var FAN_HD_GANGSHANGPAO	int32	= 2 //海底杠上炮

var FAN_TOP		int32	= 5 //顶番

//初始化麻将牌
func init() {
	mjpaiMap = make(map[int]string, 108)
	mjpaiMap[0] = "T_1"
	mjpaiMap[1] = "T_1"
	mjpaiMap[2] = "T_1"
	mjpaiMap[3] = "T_1"
	mjpaiMap[4] = "T_2"
	mjpaiMap[5] = "T_2"
	mjpaiMap[6] = "T_2"
	mjpaiMap[7] = "T_2"
	mjpaiMap[8] = "T_3"
	mjpaiMap[9] = "T_3"
	mjpaiMap[10] = "T_3"
	mjpaiMap[11] = "T_3"
	mjpaiMap[12] = "T_4"
	mjpaiMap[13] = "T_4"
	mjpaiMap[14] = "T_4"
	mjpaiMap[15] = "T_4"
	mjpaiMap[16] = "T_5"
	mjpaiMap[17] = "T_5"
	mjpaiMap[18] = "T_5"
	mjpaiMap[19] = "T_5"
	mjpaiMap[20] = "T_6"
	mjpaiMap[21] = "T_6"
	mjpaiMap[22] = "T_6"
	mjpaiMap[23] = "T_6"
	mjpaiMap[24] = "T_7"
	mjpaiMap[25] = "T_7"
	mjpaiMap[26] = "T_7"
	mjpaiMap[27] = "T_7"
	mjpaiMap[28] = "T_8"
	mjpaiMap[29] = "T_8"
	mjpaiMap[30] = "T_8"
	mjpaiMap[31] = "T_8"
	mjpaiMap[32] = "T_9"
	mjpaiMap[33] = "T_9"
	mjpaiMap[34] = "T_9"
	mjpaiMap[35] = "T_9"
	mjpaiMap[36] = "S_1"
	mjpaiMap[37] = "S_1"
	mjpaiMap[38] = "S_1"
	mjpaiMap[39] = "S_1"
	mjpaiMap[40] = "S_2"
	mjpaiMap[41] = "S_2"
	mjpaiMap[42] = "S_2"
	mjpaiMap[43] = "S_2"
	mjpaiMap[44] = "S_3"
	mjpaiMap[45] = "S_3"
	mjpaiMap[46] = "S_3"
	mjpaiMap[47] = "S_3"
	mjpaiMap[48] = "S_4"
	mjpaiMap[49] = "S_4"
	mjpaiMap[50] = "S_4"
	mjpaiMap[51] = "S_4"
	mjpaiMap[52] = "S_5"
	mjpaiMap[53] = "S_5"
	mjpaiMap[54] = "S_5"
	mjpaiMap[55] = "S_5"
	mjpaiMap[56] = "S_6"
	mjpaiMap[57] = "S_6"
	mjpaiMap[58] = "S_6"
	mjpaiMap[59] = "S_6"
	mjpaiMap[60] = "S_7"
	mjpaiMap[61] = "S_7"
	mjpaiMap[62] = "S_7"
	mjpaiMap[63] = "S_7"
	mjpaiMap[64] = "S_8"
	mjpaiMap[65] = "S_8"
	mjpaiMap[66] = "S_8"
	mjpaiMap[67] = "S_8"
	mjpaiMap[68] = "S_9"
	mjpaiMap[69] = "S_9"
	mjpaiMap[70] = "S_9"
	mjpaiMap[71] = "S_9"
	mjpaiMap[72] = "W_1"
	mjpaiMap[73] = "W_1"
	mjpaiMap[74] = "W_1"
	mjpaiMap[75] = "W_1"
	mjpaiMap[76] = "W_2"
	mjpaiMap[77] = "W_2"
	mjpaiMap[78] = "W_2"
	mjpaiMap[79] = "W_2"
	mjpaiMap[80] = "W_3"
	mjpaiMap[81] = "W_3"
	mjpaiMap[82] = "W_3"
	mjpaiMap[83] = "W_3"
	mjpaiMap[84] = "W_4"
	mjpaiMap[85] = "W_4"
	mjpaiMap[86] = "W_4"
	mjpaiMap[87] = "W_4"
	mjpaiMap[88] = "W_5"
	mjpaiMap[89] = "W_5"
	mjpaiMap[90] = "W_5"
	mjpaiMap[91] = "W_5"
	mjpaiMap[92] = "W_6"
	mjpaiMap[93] = "W_6"
	mjpaiMap[94] = "W_6"
	mjpaiMap[95] = "W_6"
	mjpaiMap[96] = "W_7"
	mjpaiMap[97] = "W_7"
	mjpaiMap[98] = "W_7"
	mjpaiMap[99] = "W_7"
	mjpaiMap[100] = "W_8"
	mjpaiMap[101] = "W_8"
	mjpaiMap[102] = "W_8"
	mjpaiMap[103] = "W_8"
	mjpaiMap[104] = "W_9"
	mjpaiMap[105] = "W_9"
	mjpaiMap[106] = "W_9"
	mjpaiMap[107] = "W_9"


	//0---27 万条筒
	clienMap = make(map[int]int32, 108)
	clienMap[0] = 19
	clienMap[1] = 19
	clienMap[2] = 19
	clienMap[3] = 19

	clienMap[4] = 20
	clienMap[5] = 20
	clienMap[6] = 20
	clienMap[7] = 20

	clienMap[8] = 21
	clienMap[9] = 21
	clienMap[10] = 21
	clienMap[11] = 21

	clienMap[12] = 22
	clienMap[13] = 22
	clienMap[14] = 22
	clienMap[15] = 22

	clienMap[16] = 23
	clienMap[17] = 23
	clienMap[18] = 23
	clienMap[19] = 23

	clienMap[20] = 24
	clienMap[21] = 24
	clienMap[22] = 24
	clienMap[23] = 24

	clienMap[24] = 25
	clienMap[25] = 25
	clienMap[26] = 25
	clienMap[27] = 25

	clienMap[28] = 26
	clienMap[29] = 26
	clienMap[30] = 26
	clienMap[31] = 26

	clienMap[32] = 27
	clienMap[33] = 27
	clienMap[34] = 27
	clienMap[35] = 27

	clienMap[36] = 10
	clienMap[37] = 10
	clienMap[38] = 10
	clienMap[39] = 10

	clienMap[40] = 11
	clienMap[41] = 11
	clienMap[42] = 11
	clienMap[43] = 11

	clienMap[44] = 12
	clienMap[45] = 12
	clienMap[46] = 12
	clienMap[47] = 12

	clienMap[48] = 13
	clienMap[49] = 13
	clienMap[50] = 13
	clienMap[51] = 13

	clienMap[52] = 14
	clienMap[53] = 14
	clienMap[54] = 14
	clienMap[55] = 14

	clienMap[56] = 15
	clienMap[57] = 15
	clienMap[58] = 15
	clienMap[59] = 15

	clienMap[60] = 16
	clienMap[61] = 16
	clienMap[62] = 16
	clienMap[63] = 16

	clienMap[64] = 17
	clienMap[65] = 17
	clienMap[66] = 17
	clienMap[67] = 17

	clienMap[68] = 18
	clienMap[69] = 18
	clienMap[70] = 18
	clienMap[71] = 18

	clienMap[72] = 1
	clienMap[73] = 1
	clienMap[74] = 1
	clienMap[75] = 1

	clienMap[76] = 2
	clienMap[77] = 2
	clienMap[78] = 2
	clienMap[79] = 2

	clienMap[80] = 3
	clienMap[81] = 3
	clienMap[82] = 3
	clienMap[83] = 3

	clienMap[84] = 4
	clienMap[85] = 4
	clienMap[86] = 4
	clienMap[87] = 4

	clienMap[88] = 5
	clienMap[89] = 5
	clienMap[90] = 5
	clienMap[91] = 5

	clienMap[92] = 6
	clienMap[93] = 6
	clienMap[94] = 6
	clienMap[95] = 6

	clienMap[96] = 7
	clienMap[97] = 7
	clienMap[98] = 7
	clienMap[99] = 7

	clienMap[100] = 8
	clienMap[101] = 8
	clienMap[102] = 8
	clienMap[103] = 8

	clienMap[104] = 9
	clienMap[105] = 9
	clienMap[106] = 9
	clienMap[107] = 9

	//番数 顶番5

}


//统计麻将牌
func GettPaiStats(pais []*MJPai) []int {
	//统计每张牌的重复次数
	//log.T("GettPaiStats : %v", pais)
	counts := make([]int, 27) //0~27
	for i := 0; i < len(pais); i++ {
		pai := pais[i]
		value := pai.GetValue() - 1
		flower := pai.GetFlower()    //flower=1,2,3
		//log.T("getValue(%v),pai.GetFlower(%v) ", value, pai.GetFlower())
		value += (flower - 1) * 9
		//log.T("value[%v]", value)
		//log.T("value,f")
		counts[ value ] ++
	}
	return counts
}

func is19(val int) bool {
	return (val % 9 == 0) || (val % 9 == 8)
}

//七对 龙七对牌型胡牌判断
func tryHU7(handPai *MJHandPai, handCounts[] int) (canHu bool, isAll19 bool) {
	canHu, isAll19 = false, false
	if IsQiDui(handPai, handCounts) || IsLongQiDui(handPai, handCounts) {
		canHu = true
	}
	return canHu, isAll19
}

//33332牌型胡牌的算法
func tryHU(count []int, len int) (result bool, isAll19 bool) {
	//log.T("开始判断tryHu(%v,%v)", count, len)
	isAll19 = true //全带幺
	result = false
	//递归完所有的牌表示 胡了
	if (len == 0) {
		log.T("len == 0")
		return true, isAll19
	}

	if (len % 3 == 2) {
		//log.T("if %v 取模 3 == 2", len)
		// 说明对牌出现
		for i := 0; i < 27; i++ {
			if (count[i] >= 2) {
				count[i] -= 2
				result, isAll19 = tryHU(count, len - 2)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if ! is19(i) {
						//不是幺九
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 2
			}
		}
	} else {
		//log.T("else %v", len)
		// 是否是顺子，这里应该分开判断
		for i := 0; i < 7; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !is19(i) && !is19(i + 1) && !is19(i + 2) {
						//不是幺九
						//log.T("branch 2 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 1
				count[i + 1] += 1
				count[i + 2] += 1
			}
		}

		for i := 9; i < 16; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1
				count[i + 1] -= 1
				count[i + 2] -= 1
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !is19(i) && !is19(i + 1) && !is19(i + 2) {
						//不是幺九
						//log.T("branch 3 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 1
				count[i + 1] += 1
				count[i + 2] += 1
			}
		}

		for i := 18; i < 25; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !is19(i) && !is19(i + 1) && !is19(i + 2) {
						//不是幺九
						//log.T("branch 4 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 1;
				count[i + 1] += 1;
				count[i + 2] += 1;
			}
		}

		// 三个一样的
		for i := 0; i < 27; i++ {
			if (count[i] >= 3) {
				count[i] -= 3
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !is19(i) {
						//不是幺九
						//log.T("branch 5 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 3
			}
		}

	}

	return false, isAll19
}


func CanHuPai(handPai *MJHandPai) (bool,bool) {
	//在所有的牌中增加 pai,判断此牌是否能和
	pais := []*MJPai{}
	pais = append(pais, handPai.Pais...)
	pais = append(pais, handPai.InPai)

	counts := GettPaiStats(pais)

	var canHu, isAll19 bool

	//七对 龙七对牌型 不带幺九
	canHu, isAll19 = tryHU7(handPai, counts)
	if canHu {
		return canHu,isAll19
	}

	//普通33332牌型
	canHu, isAll19 = tryHU(counts, len(pais))
	if canHu {
		log.T("牌= %v  可以胡! isAll19=%v", handPai.InPai.LogDes(), isAll19)
	} else {
		//log.T("牌= %v  不能胡! isAll19=%v", handPai.InPai.LogDes(), isAll19)
	}

	return canHu,isAll19
}

func GetHuScore(handPai *MJHandPai, isZimo bool, is19 bool, extraAct HuPaiType, roomInfo RoomTypeInfo, mjDesk *MjDesk) (fan int32, score int64, huCardStr[] string) {

	log.T("pai: %v", handPai.GetDes(), handPai.InPai.LogDes())

	//底分
	score = int64(*roomInfo.BaseValue)

	//取得番数
	huFan, huCardStr := getHuFan(handPai, isZimo, is19, extraAct, mjDesk)

	for i := int32(0); i < huFan; i++ {
		score *= 2
	}

	log.T("胡[%d]番 (%v)", huFan, huCardStr)

	if isZimo {
		if mjDesk.IsNeedZiMoJiaDi() {
			//自摸加底
			score += *roomInfo.BaseValue
		}
	}

	return huFan, score, huCardStr
}

//计算带几个"勾"
func getGou(handPai *MJHandPai, handCounts[] int) (gou int32) {
	// 已杠的牌
	gou = int32(len(handPai.GangPais))
	gou = gou / 4      //杠牌/4才是gou 的数目

	log.T("杠牌的勾:%v", gou)
	// 计算 碰牌+手牌 的勾数
	for _, pai := range handPai.Pais {
		for _, pengPai := range handPai.PengPais {
			if pai.GetClientId() == pengPai.GetClientId() {
				gou ++
				break
			}
		}
	}

	log.T("碰牌的勾:%v", gou)


	// 计算手牌中的勾数(未暗杠)
	for _, cnt := range handCounts {
		if cnt == 4 {
			gou ++
		}
	}

	log.T("手牌的勾:%v", gou)

	return gou
}

// 返回胡牌番数
// extraAct:指定HuPaiType.H_GangShangHua(杠上花/炮,海底等)
//
func getHuFan(handPai *MJHandPai, isZimo bool, is19 bool, extraAct HuPaiType, mjDesk *MjDesk) (fan int32, huCardStr[] string) {
	fan = int32(0)
	pais := []*MJPai{}
	pais = append(pais, handPai.Pais...)
	pais = append(pais, handPai.InPai)

	handCounts := GettPaiStats(pais) //计算手牌的每张牌数量

	pais = append(pais, handPai.PengPais...)
	pais = append(pais, handPai.GangPais...)

	isQingYiSe := IsQingYiSe(pais) //清一色
	log.T("判断是否是清一色: %v", isQingYiSe)

	/*
	isDaDuiZi := IsDaDuiZi(pais) //大对子
	log.T("判断是否是大对子: %v", isDaDuiZi)

	//isQiDui := IsQiDui(handPai, handCounts) //七对
	//log.T("判断是否是七对: %v", isQiDui)

	//isLongQiDui := IsLongQiDui(handPai, handCounts) //龙七对
	//log.T("判断是否是龙七对: %v", isLongQiDui)
	//*/
	//log.T("判断是否是将对: %v", IsJiangDui(pais))


	isCountGou := true //是否计算勾 七对 龙七对 清龙七对 将七对 不算勾

	switch  {
	case IsLongQiDui(handPai, handCounts) : //case 清龙七对 龙七对
		log.T("是龙七对")
		isCountGou = false
		if isQingYiSe { //清龙七对
			log.T("是清龙七对")
			fan = FAN_QINGLONGQIDUI
			huCardStr = append(huCardStr, "清龙七对")
		}else { //龙七对
			log.T("是龙七对")
			fan = FAN_LONGQIDUI
			huCardStr = append(huCardStr, "龙七对")
		}
	case IsQiDui(handPai, handCounts): //case 清七对 将七对 七对
		log.T("是七对")
		isCountGou = false
		if isQingYiSe { //清七对
			log.T("是清七对")
			fan = FAN_QINGQIDUI
			huCardStr = append(huCardStr, "清七对")
		}else { //七对
			log.T("是七对")
			fan = FAN_QIDUI
			huCardStr = append(huCardStr, "七对")
		}
	case IsDaDuiZi(pais): //case 清对 将对 大对子
		log.T("是大对子")
		if isQingYiSe { //清对
			log.T("是清对")
			fan = FAN_QINGDUI
			huCardStr = append(huCardStr, "清对")
		}else if mjDesk.IsNeedYaojiuJiangdui() { //将对选项开启
			log.T("是将对")
			if IsJiangDui(handPai) {
				fan = FAN_DADUIZI
				huCardStr = append(huCardStr, "将对")
			}
		}else { //大对子
			log.T("是大对子")
			fan = FAN_DADUIZI
			huCardStr = append(huCardStr, "大对子")
		}
	default: //default 清一色 平胡
		if isQingYiSe { //平胡清一色
			log.T("是清一色")
			fan = FAN_QINGYISE
			huCardStr = append(huCardStr, "清一色")
		}else { //平胡
			log.T("是平胡")
			fan = FAN_PINGHU
			huType := "平胡"
			huCardStr = append(huCardStr, huType)

		}
	}

	//TODO MJOption_JIANGOUDIAO
	//if IsOpenRoomOption(roomInfo.PlayOptions.OthersCheckBox, MJOption_JINGOUDIAO) { //金钩钓
	//	if IsJingGouDiao(handCounts) {
	//		fan += FAN_JINGOUDIAO
	//		huCardStr = append(huCardStr, "金钩钓")
	//	}
	//}

	//附加选项
	if mjDesk.IsNeedYaojiuJiangdui() { //带幺九选项开启
		if is19 && IsPengGang19(handPai) { //手牌带幺九 且 碰杠牌带幺九
			fan += FAN_DAIYAOJIU
			huCardStr = append(huCardStr, "带幺九")
		}
	}

	if mjDesk.IsNeedMenqingZhongzhang() { //门清中张选项开启
		if IsMenqing(handPai) {
			fan += FAN_MENQ_ZHONGZ
			huCardStr = append(huCardStr, "门清")
		}
		if IsZhongzhang(handPai, handCounts) {
			fan += FAN_MENQ_ZHONGZ
			huCardStr = append(huCardStr, "中张")
		}
	}
	isTianDiHuFlag := false //天地胡选项 避免多次搜索
	if mjDesk.IsNeedMenqingZhongzhang() { //天地胡选项开启
		isTianDiHuFlag = true
	}

	switch HuPaiType(extraAct) {

	//天地胡为牌型番数，非加番
	case HuPaiType_H_TianHu :
		if isTianDiHuFlag { //天地胡选项开启
			fan = FAN_TIAN_DI_HU
			huCardStr = append(huCardStr, "天胡")
		}
	case HuPaiType_H_DiHu :
		if isTianDiHuFlag { //天地胡选项开启
			fan = FAN_TIAN_DI_HU
			huCardStr = append(huCardStr, "地胡")
		}

	case HuPaiType_H_GangShangHua:
		fan += FAN_GANGSHANGHUA
		huCardStr = append(huCardStr, "杠上花")

	case HuPaiType_H_GangShangPao:
		fan += FAN_GANGSHANGPAO
		huCardStr = append(huCardStr, "杠上炮")

	case HuPaiType_H_HaiDiLao:
		fan += FAN_HD_HUA
		huCardStr = append(huCardStr, "海底花")

	case HuPaiType_H_HaiDiPao:
		fan += FAN_HD_PAO
		huCardStr = append(huCardStr, "海底炮")

	case HuPaiType_H_QiangGang:
		fan += FAN_QIANGGANG
		huCardStr = append(huCardStr, "抢杠")

	case HuPaiType_H_HaidiGangShangHua:
		fan += FAN_HD_GANGSHANGHUA
		huCardStr = append(huCardStr, "海底杠上花")

	case HuPaiType_H_HaidiGangShangPao:
		fan += FAN_HD_GANGSHANGPAO
		huCardStr = append(huCardStr, "海底杠上炮")
	default:
	}

	//自摸
	if isZimo {
		if mjDesk.IsNeedZiMoJiaFan() {
			fan += FAN_ZIMO
		} else if mjDesk.IsNeedZiMoJiaDi() {
			//result += di
		}
		huCardStr = append(huCardStr, "自摸")
	}

	// 计算有几个"勾"
	if isCountGou {
		log.T("加勾")
		gou := getGou(handPai, handCounts)

		fan += gou
		if gou > 0 {
			str, _ := numUtils.Int2String(gou)
			huCardStr = append(huCardStr, "勾X" + str)
		}
	}

	//顶番
	if fan > FAN_TOP {
		fan = FAN_TOP
	}

	return fan, huCardStr
}


//这张pai是否可碰
// add 增加缺的花色是不能碰的
func CanPengPai(pai *MJPai, handPai *MJHandPai) bool {

	existCount := 0
	for i := 0; i < len(handPai.Pais); i++ {

		//add 判断是否是定缺的花色
		if pai.GetFlower() == handPai.GetQueFlower() {
			continue
		}

		if *pai.Flower == *handPai.Pais[i].Flower && *pai.Value == *handPai.Pais[i].Value {
			existCount ++
		}
	}

	return ( existCount == 2 || existCount == 3 )
}

//这张pai是否可杠( 当pai为nil时, 检测handPai中是否有杠)
func CanGangPai(pai *MJPai, handPai *MJHandPai) (canGang bool, gangPais []*MJPai) {
	if ( pai != nil ) {
		//判断别人打入的牌是否可杠,别人打得牌，不能判断碰牌
		existCount := 0
		for _, p := range handPai.Pais {
			//判断是不是定缺的花色
			if p.GetFlower() == handPai.GetQueFlower() {
				continue
			}

			if p.GetClientId() == pai.GetClientId() {
				existCount ++
			}
		}

		canGang = ( existCount == 3 )
		if ( canGang ) {
			gangPais = append(gangPais, pai)
		}

	} else {
		//检测手牌中是否有杠
		tempPais := make([]*MJPai, len(handPai.Pais) + 1 + len(handPai.PengPais))
		copy(tempPais[0:len(handPai.Pais)], handPai.Pais)
		tempPais[len(handPai.Pais)] = handPai.InPai
		copy(tempPais[len(handPai.Pais) + 1:], handPai.PengPais)

		//counts := GettPaiStats(handPai.Pais)
		counts := GettPaiStats(tempPais)
		for _, p := range tempPais {
			if p.GetFlower() == handPai.GetQueFlower() {
				continue
			}
			//log.T("判断杠牌 p.getValue(%v)+p.GetFlower[%v]*9 = %v", p.GetValue(), p.GetFlower(), p.GetValue() + (p.GetFlower() - 1) * 9)
			if ( 4 == counts[ p.GetValue() - 1 + (p.GetFlower() - 1) * 9 ] ) {
				canGang = true
				gangPais = append(gangPais, p)
			}
		}
	}

	return canGang, gangPais
}


func IsPengGang19(handPai *MJHandPai) bool {
	pengPais := handPai.PengPais
	gangPais := handPai.GangPais
	if pengPais != nil {
		for i := 0; i < len(pengPais); i++ {
			if *pengPais[i].Value != 1 || *pengPais[i].Value != 9 { //
				return false
			}
		}
	}

	if gangPais != nil {
		for i := 0; i < len(gangPais); i++ {
			if *gangPais[i].Value != 1 || *gangPais[i].Value != 9 { //
				return false
			}
		}
	}
	return true
}

//金钩钓牌型判断 手牌只有一对
func IsJingGouDiao(handCounts []int) bool {
	var count int = 0
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] != 2 {
			return false
		}else {
			count++
		}
	}
	if count != 1 {
		return false
	}
	return true
}

//清一色
func IsQingYiSe(pais []*MJPai) bool {
	flower := pais[0].Flower
	for i := 1; i < len(pais); i++ {
		if *flower != *pais[i].Flower {
			return false //不是清一色
		}
	}

	return true
}

//大对子
func IsDaDuiZi(pais []*MJPai) bool {
	counts := GettPaiStats(pais)
	log.T("判断是否是大对子的统计数据:%v", counts)

	jiangDui := 0

	for i := 0; i < len(counts); i++ {
		//log.T("判断是否是大对子count[%v] = %v", i, counts[i])
		//count := counts [ pais[i].GetValue() - 1 + (pais[i].GetFlower() - 1) * 9]
		//if count == 2 {
		//	jiangDui ++
		//	if jiangDui > 1 {
		//		return false
		//	}
		//} else if count < 2 {
		//	return false
		//}

		if counts[i] == 2 {
			jiangDui ++
			if jiangDui > 1 {
				log.T("不是大对子")
				return false
			}
		} else if counts[i] == 1 {
			log.T("不是大对子")
			return false
		}
	}
	log.T("是大对子")
	return true
}

//七对
func IsQiDui(handPai *MJHandPai, handCounts[] int) bool {
	pais := handPai.Pais

	if handPai.GangPais != nil || handPai.PengPais != nil { //不能有碰杠
		return false
	}
	if len(pais) != 13 {
		//手牌需为13张
		return false
	}

	//for i := 0; i < len(pais); i++ {
	//	if handCounts [ pais[i].GetValue() - 1 + (pais[i].GetFlower() - 1) * 9 ] != 2 {
	//		//每张牌都是2张
	//		return false
	//	}
	//}

	for i := 0; i < len(handCounts); i++ {
		if (handCounts [i] != 2) && (handCounts[i] != 0) { //每张牌都是2张
			return false
		}
	}

	return true
}

//龙七对
func IsLongQiDui(handPai *MJHandPai, handCounts[] int) bool {
	//pais := handPai.Pais

	//if !IsQiDui(handPai, handCounts) {
	//	//首先是七对
	//	return false
	//}
	//
	//for i := 0; i < len(pais); i++ {
	//	if handCounts [ *pais[i].Value - 1 + (pais[i].GetFlower() - 1) * 9 ] == 4 {
	//		//有一杠
	//		return true
	//	}
	//}
	longCount := 0 //杠数
	duiCount := 0 //对数

	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 4 { //杠
			longCount++
		}
		if (handCounts[i] == 2) {
			duiCount++
		}
		if (handCounts[i] != 0) && (handCounts[i] != 4) && (handCounts[i] != 2) { //牌数不符合
			//log.T("isLongQiDui: 牌数不符合 0、2、4")
			return false
		}
	}
	if (longCount < 1) || (duiCount < 5) { //杠数小于一，对数小于5
		//log.T("isLongQiDui: 杠对数不符合")
		return false
	}
	return true
}

//将对(全是2,5,8的大对子)
func IsJiangDui(handPai *MJHandPai) bool {
	pais := handPai.Pais
	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 2 && *pais[i].Value != 5 && *pais[i].Value != 8 {
			return false
		}
	}
	return IsDaDuiZi(pais) //是大对子
}

//将七对(全是2,5,8的七对)
//func IsJiangQiDui(handPai *MJHandPai, handCounts[] int) bool {
//	pais := handPai.Pais
//
//	for i := 0; i < len(pais); i++ {
//		if *pais[i].Value != 2 && *pais[i].Value != 5 && *pais[i].Value != 8 {
//			return false
//		}
//	}
//
//	return IsQiDui(handPai, handCounts) //是七对
//}

//门清 没有明杠 碰牌
func IsMenqing(handPai *MJHandPai) bool {
	//TODO 明杠判断
	if handPai.PengPais != nil { //含碰牌
		return false
	}
	return true
}

//中张 没有1、9
func IsZhongzhang(handPai *MJHandPai, handCounts []int) bool {
	//
	for i := 0; i < len(handCounts); i++ {
		switch (i + 1) % 9 {
		case 1, 0 : //牌值为1、9
			if handCounts[i] > 0 { //牌中包含1、9
				return false
			}
		default:
		}
	}
	return true
}

//全带幺
func IsAllDaiYao(handPai *MJHandPai) bool {
	//pais := handPai.Pais
	//counts := GettPaiStats(handPai)
	//for i := 0; i < len(pais); i++ {
	//	if *pais[i].Value != 1 && *pais[i].Value != 2 && *pais[i].Value != 3 && *pais[i].Value != 7 && *pais[i].Value != 8 && *pais[i].Value != 9 {
	//		return false
	//	}
	//}
	//tryHu后的牌型里不含4、5、6即带幺九
	//count := 0
	//for i := 3; i < len(counts); i++ { //
	//	count++
	//	if count == 3 {
	//		count = 0
	//		continue
	//	}
	//	if counts[i] != 0 {
	//		return false
	//	}
	//}

	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//通过描述初始化牌的同条玩和大小
func (p *MJPai) InitByDes() error {

	//拆分描述
	sarry := strings.Split(p.GetDes(), "_")
	flowerStr := sarry[0]

	//初始化花色
	if flowerStr == "T" {
		*p.Flower = T
	} else if flowerStr == "S" {
		*p.Flower = S
	} else {
		*p.Flower = W
	}

	//初始化大小
	*p.Value = int32(numUtils.String2Int(sarry[1]))
	return nil
}

//洗牌的算法,这里可以得到一副洗好的麻将
func XiPai() []*MJPai {

	//初始化一个顺序的牌的集合
	pmap := make([]int, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		pmap[i] = i
	}

	//打乱牌的集合
	pResult := make([]int, MJPAI_COUNT)

	for i := 0; i < MJPAI_COUNT; i++ {
		rand := utils.Rand(0, (int32(MJPAI_COUNT - i)))
		//log.T("得到的rand[%v]", rand)
		pResult[i] = pmap[rand]
		pmap = append(pmap[:int(rand)], pmap[int(rand) + 1:]...)
	}

	log.T("洗牌之后,得到的随机的index数组[%v]", pResult)
	//TestCheckRanIndex(pResult)        //todo  测试代码，之后需要删除
	//pResult = []int{94, 55, 13, 40, 106, 28, 100, 31, 57, 44, 83, 58, 101, 104, 9, 92, 62, 67, 25, 38, 41, 86, 6, 48, 65, 61, 71, 4, 36, 49, 63, 7, 26, 75, 56, 43, 35, 103, 91, 27, 33, 50, 3, 11, 10, 8, 93, 76, 45, 12, 22, 68, 29, 17, 105, 19, 23, 95, 87, 34, 53, 37, 102, 88, 0, 69, 79, 80, 60, 81, 47, 85, 98, 2, 15, 16, 96, 59, 42, 30, 99, 77, 64, 46, 84, 52, 51, 70, 32, 21, 66, 1, 90, 72, 97, 78, 20, 74, 14, 24, 89, 5, 18, 73, 54, 82, 39, 107}

	//开始得到牌的信息
	result := make([]*MJPai, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		result[i] = InitMjPaiByIndex(pResult[i])
	}

	log.T("洗牌之后,得到的牌的数组[%v]", result)
	return result
}

func XiPaiTestHu() []*MJPai {
	//初始化一个顺序的牌的集合
	pmap := make([]int, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		pmap[i] = i
	}

	//打乱牌的集合
	pResult := make([]int, MJPAI_COUNT)
	for i := 0; i < 108; i++ {
		pResult[i] = i;
	}

	//pResult[0] = 15
	//pResult[1] = 16
	//pResult[4] = 19
	//pResult[5] = 20
	//
	//pResult[15] = 0
	//pResult[16] = 1
	//pResult[19] = 4
	//pResult[20] = 5

	//开始得到牌的信息
	result := make([]*MJPai, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		result[i] = InitMjPaiByIndex(pResult[i])
	}

	log.T("洗牌之后,得到的牌的数组[%v]", result)
	return result
}

func TestCheckRanIndex(r []int) {
	copyCheck := make([]int, MJPAI_COUNT)
	copy(copyCheck, r)
	for _, a := range copyCheck {
		count := 0
		for _, b := range copyCheck {
			if a == b {
				count ++
				if count > 1 {
					log.Fatal("随机牌的index[%v] 出错", a)
					panic(errors.New("随便牌的index 出错"))
				}
			}
		}
	}
}


//通过一个index索引来得到一张牌
func InitMjPaiByIndex(index int) *MJPai {
	result := NewMjpai()
	*result.Index = int32(index)
	*result.Des = mjpaiMap[index]
	result.InitByDes()
	return result
}

//通过一个Des描述和现持有的牌来得到一张空闲牌
func InitMjPaiByDes(des string, hand *MJHandPai) *MJPai {

	result := NewMjpai()
	handMJPais := []*MJPai{}

	//加入杠牌
	handMJPais = append(handMJPais, hand.GangPais...)
	//加入手牌
	handMJPais = append(handMJPais, hand.Pais...)
	//加入碰牌
	handMJPais = append(handMJPais, hand.PengPais...)
	//加入摸牌
	handMJPais = append(handMJPais, hand.InPai)

	for mjpaiIndex, mjpaiDes := range mjpaiMap {
		for _, handMJPai := range handMJPais {
			if des == mjpaiDes {
				if (handMJPai != nil) && (int32(mjpaiIndex) == *handMJPai.Index) {
					continue
				} else {
					*result.Index = int32(mjpaiIndex)
					*result.Des = mjpaiDes
					result.InitByDes()
					return result
				}
			}
		}
	}

	return nil
}

func GetFlow(f int32) string {
	if f == 1 {
		return "万"
	} else if f == 2 {
		return "条"
	} else if f == 3 {
		return "筒"
	} else {
		return "白"
	}

}

//判断是否开启房间的某个选项
func IsOpenRoomOption(othersCheckBox []int32, option MJOption) bool {
	for _, opt := range othersCheckBox {
		//判断是否开启房间的某个选
		if opt == int32(option) {
			return true
		}
	}
	return false
}

//
