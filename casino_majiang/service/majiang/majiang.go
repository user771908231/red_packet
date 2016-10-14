package majiang

import (
	"strings"
	"github.com/name5566/leaf/log"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	. "casino_majiang/msg/protogo"
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
	clienMap[52] = 13
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
}

//一张麻将牌的结构


//开始的时候,一手麻将牌,关于这幅麻将牌的统计都在这里
type PairMjPai struct {
	pais          []*MJPai
	paiStatistics []int
	fan           int32   //这局牌 有多少番
	que           int32   //缺什么牌,不要什么花色

}


//统计麻将牌
func (p *PairMjPai)  InitHandPaisStats(handPai *MJHandPai) []int {
	//统计牌的张数

	//p.InitPaiStats( handPai.PengPais )  //已碰的牌
	//p.InitPaiStats( handPai.GangPais )  //已杠的牌

	return p.GettPaiStats( handPai.Pais ) //手里剩余牌
}

func (p *PairMjPai)  GettPaiStats(pais []*MJPai) []int {
	//统计每张牌的重复次数
	counts := make([]int, 27) //0~27
	for i := 0; i < len(pais); i++ {
		pai := pais[i]
		value := pai.GetValue() - 1
		flower := pai.GetFlower()    //flower=1,2,3

		value += (flower-1)*9

		counts[ value ] ++
	}

	return counts
}

func is19(val int) bool {
	return (val % 9 == 0) || (val % 9 == 8)
}

//胡牌的算法
func tryHU(count []int, len int) (result bool, isAll19 bool) {
	isAll19 = true //全带幺
	result = false

	//递归完所有的牌表示 胡了
	if (len == 0) {
		return true, isAll19
	}

	if (len % 3 == 2) {
		// 说明对牌没出现
		for i := 0; i < 27; i++ {
			if (count[i] >= 2) {
				count[i] -= 2
				result, isAll19 = tryHU(count, len - 2)
				if (result) {
					if ! is19(i) { //不是幺九
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 2
			}
		}
	} else {
		// 三个一样的
		for i := 0; i < 27; i++ {
			if (count[i] >= 3) {
				count[i] -= 3
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					if !is19(i) { //不是幺九
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 3
			}
		}
		// 是否是顺子，这里应该分开判断
		for i := 0; i < 7; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				result, isAll19 = tryHU(count, len - 3)
				if (result) {
					if !is19(i) && !is19(i+1) && !is19(i+2) { //不是幺九
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
					if !is19(i) && !is19(i+1) && !is19(i+2) { //不是幺九
						isAll19 = false
					}
					return true,isAll19
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
					if !is19(i) && !is19(i+1) && !is19(i+2) { //不是幺九
						isAll19 = false
					}
					return true, isAll19
				}
				count[i] += 1;
				count[i + 1] += 1;
				count[i + 2] += 1;
			}
		}

	}

	return false, isAll19
}

//确定要胡牌的时候,做出的处理
func (p *PairMjPai) HuPai(handPai *MJHandPai) error {
	//排序
	//统计
	//p.InitHandPaisStats()
	//p.IsDaDui        //判断是否是大队子
	//p.IsQingYiSe        //判断是否是清一色
	//p.IsJiangDui        //判断是否是将对

	return nil

}

//这张pai能不能胡
func (p *PairMjPai) CanHuPai(pai *MJPai, handPai *MJHandPai) bool {
	//在所有的牌中增加 pai,判断此牌是否能和
	handPai.Pais = append(handPai.Pais, pai)
	counts := p.InitHandPaisStats( handPai )

	canHu, isAll19 := tryHU(counts, len(handPai.Pais))
	if canHu {
		log.Debug("牌= %d  可以胡! isAll19=%V", *pai.Value, isAll19)
	} else {
		log.Debug("牌= %d  不能胡! isAll19=%V", *pai.Value, isAll19)
	}

	//最后需要删除最后一张牌
	handPai.Pais = handPai.Pais[0:len(handPai.Pais) - 1]

	return canHu
}


func (p *PairMjPai) getHuScore(handPai *MJHandPai, isZimo bool, extraAct HuPaiType, roomInfo RoomTypeInfo ) (fan int32, score int64, huCardStr[] string) {
	//底分
	score = int64(*roomInfo.BaseValue)

	//取得番数
	huFan, huCardStr := p.getHuFan(handPai, isZimo, extraAct, roomInfo)

	for i:=int32(0); i< huFan; i++ {
		score *= 2
	}

	log.Debug("胡[%d]番 (%v)", huFan, huCardStr)

	if isZimo {
		if MJOption(*roomInfo.PlayOptions.ZiMoRadio) == MJOption_ZIMO_JIA_DI { //自摸加底
			score += *roomInfo.BaseValue
		}
	}

	return huFan, score, huCardStr
}

//计算带几个"勾"
func (p *PairMjPai) getGou(handPai *MJHandPai, counts[] int ) (gou int32) {
	// 已杠的牌
	gou = int32(len(handPai.GangPais))

	// 计算 碰牌+手牌 的勾数
	for _,pengPai := range handPai.PengPais {
		for _, pai := range handPai.Pais {
			if (*pengPai.Flower == *pai.Flower && *pengPai.Value == *pai.Value ) {
				gou ++
			}
		}
	}

	// 计算手牌中的勾数(未暗杠)
	for _, pai := range handPai.Pais {
		count := counts [ pai.GetValue()-1 + (pai.GetFlower()-1)*9]
		if count == 4 {
			gou ++
		}
	}

	return gou
}

// 返回胡牌番数
// extraAct:指定HuPaiType.H_GangShangHua(杠上花/炮,海底等)
//
func (p *PairMjPai) getHuFan(handPai *MJHandPai, isZimo bool, extraAct HuPaiType, roomInfo RoomTypeInfo ) (fan int32,  huCardStr[] string) {
	fan = int32(0)
	pais := []*MJPai{}
	pais = append(pais, handPai.Pais...)
	pais = append(pais, handPai.InPai)


	handCounts := p.GettPaiStats( handPai.Pais ) //计算手牌的每张牌数量


	isDaDuiZi := p.IsDaDuiZi( pais ) //大对子

	isQingYiSe := p.IsQingYiSe( pais ) //清一色

	isQiDui := p.IsQiDui( handPai ) //七对

	isLongQiDui := p.IsLongQiDui( handPai ) //龙七对

	if isQiDui {
		if isLongQiDui && isQingYiSe {
			fan = 5
			huCardStr = append(huCardStr, "清龙七对")
		} else if isLongQiDui {
			fan = 3
			huCardStr = append(huCardStr, "龙七对")
		} else if isQingYiSe {
			fan = 4
			huCardStr = append(huCardStr, "清七对")
		} else if p.IsJiangQiDui( handPai ) {
			fan = 4
			huCardStr = append(huCardStr, "将七对")
		}
	} else if isDaDuiZi && isQingYiSe {
		fan = 3
		huCardStr = append(huCardStr, "清对")
	} else if isDaDuiZi {
		fan = 1
		huCardStr = append(huCardStr, "大对子")
	} else {
		fan = 1
		huType := "平胡"

		//TODO: if 附加选项开启时
		for _, opt := range roomInfo.PlayOptions.OthersCheckBox {
			switch MJOption(opt) {
			case MJOption_YAOJIU_JIANGDUI: { //幺九将对
				if p.IsJiangDui( handPai ) {
					fan = 2
					huType = "将对"
				} else if p.IsAllDaiYao( handPai ) {
					fan = 2
					huType = "带幺九"
				}
			}

			case MJOption_TIAN_DI_HU: { //天地胡

			}
			case MJOption_KA_ER_TIAO: { //卡2条

			}
			case MJOption_TIAN_DI_HU: { //天地胡

			}

			case MJOption_MENQING_MID_CARD: { //门清中张

			}

			default:

			}
		}
		huCardStr = append(huCardStr, huType)

	}

	switch HuPaiType(extraAct) {
		case HuPaiType_H_GangShangHua:
			fan += 1
			huCardStr = append(huCardStr, "杠上花")

		case HuPaiType_H_GangShangPao:
			fan += 1
			huCardStr = append(huCardStr, "杠上炮")

		case HuPaiType_H_HaiDiLao:
			fan += 1
			huCardStr = append(huCardStr, "海底花")

		case HuPaiType_H_HaiDiPao:
			fan += 1
			huCardStr = append(huCardStr, "海底炮")

		case HuPaiType_H_QiangGang:
			fan += 1
			huCardStr = append(huCardStr, "抢杠")

		case HuPaiType_H_HaidiGangShangHua:
			fan += 1
			huCardStr = append(huCardStr, "海底杠上花")

		case HuPaiType_H_HaidiGangShangPao:
			fan += 2
			huCardStr = append(huCardStr, "海底杠上炮")
		default:
	}

	//自摸
	if isZimo {
		if MJOption(*roomInfo.PlayOptions.ZiMoRadio) == MJOption_ZIMO_JIA_FAN {
			fan += 1
		} else if MJOption(*roomInfo.PlayOptions.ZiMoRadio) == MJOption_ZIMO_JIA_DI {
			//result += di
		}
		huCardStr = append(huCardStr, "自摸")
	}

	// 计算有几个"勾"
	gou := p.getGou(handPai, handCounts)

	fan += gou
	if gou > 0 {
		str, _ := numUtils.Int2String(gou)
		huCardStr = append(huCardStr, "勾X"+str)
	}

	return fan, huCardStr
}


//这张pai是否可碰
func (p *PairMjPai) CanPengPai(pai *MJPai, handPai *MJHandPai) bool {

	existCount := 0
	for i:=0; i < len(handPai.Pais); i++ {
		if *pai.Value == *handPai.Pais[i].Value {
			existCount ++
		}
	}

	return ( existCount == 2 || existCount == 3 )
}


//这张pai是否可杠
func (p *PairMjPai) CanGangPai(pai *MJPai, handPai *MJHandPai) bool {

	existCount := 0
	for i:=0; i < len(handPai.Pais); i++ {
		if *pai.Value == *handPai.Pais[i].Value {
			existCount ++
		}
	}

	return ( existCount == 3 )
}

//清一色
func (p *PairMjPai) IsQingYiSe(pais []*MJPai) bool {
	flower := pais[0].Flower
	for i := 1; i < len(pais); i++ {
		if *flower != *pais[i].Flower {
			return false //不是清一色
		}
	}

	return true
}

//大对子
func (p *PairMjPai) IsDaDuiZi(pais []*MJPai) bool {
	counts := p.GettPaiStats( pais )

	jiangDui := 0
	for i := 0; i < len(pais); i++ {
		count := counts [ pais[i].GetValue() - 1 + (pais[i].GetFlower()-1)*9]
		if count == 2 {
			jiangDui ++
			if jiangDui > 1 {
				return false
			}
		} else if count < 2 {
			return false
		}
	}

	return true
}

//七对
func (p *PairMjPai) IsQiDui(handPai *MJHandPai) bool {

	if len(handPai.Pais) != 13 { //手牌需为13张
		return false
	}

	pais := handPai.Pais
	counts := p.GettPaiStats( pais )
	for i := 0; i < len(pais); i++ {
		if counts [ pais[i].GetValue() - 1 + (pais[i].GetFlower()-1)*9 ] != 2 {
			return false
		}
	}

	return true
}

//龙七对
func (p *PairMjPai) IsLongQiDui(handPai *MJHandPai) bool {
	pais := handPai.Pais

	if !p.IsQiDui( handPai ) { //首先是七对
		return false
	}

	counts := p.GettPaiStats( pais )
	for i := 0; i < len(pais); i++ {
		if counts [ *pais[i].Value - 1 + (pais[i].GetFlower()-1)*9 ] == 4 { //有一杠
			return true
		}
	}

	return false
}

//将对(全是2,5,8的大对子)
func (p *PairMjPai) IsJiangDui(handPai *MJHandPai) bool {
	pais := handPai.Pais

	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 2 && *pais[i].Value != 5 && *pais[i].Value != 8 {
			return false
		}
	}

	return p.IsDaDuiZi(pais) //是大对子
}

//将七对(全是2,5,8的七对)
func (p *PairMjPai) IsJiangQiDui(handPai *MJHandPai) bool {
	pais := handPai.Pais

	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 2 && *pais[i].Value != 5 && *pais[i].Value != 8 {
			return false
		}
	}

	return p.IsQiDui(handPai) //是七对
}

//全带幺
func (p *PairMjPai) IsAllDaiYao(handPai *MJHandPai) bool {
	pais := handPai.Pais

	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 1 && *pais[i].Value != 2 && *pais[i].Value != 3 && *pais[i].Value != 7 && *pais[i].Value != 8 && *pais[i].Value != 9 {
			return false
		}
	}

	return true
}

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
		//log.Debug("得到的rand[%v]", rand)
		pResult[i] = pmap[rand]
		pmap = append(pmap[:int(rand)], pmap[int(rand) + 1:]...)
	}

	log.Debug("洗牌之后,得到的随机的index数组[%v]", pResult)

	//开始得到牌的信息
	result := make([]*MJPai, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		result[i] = InitMjPaiByIndex(pResult[i])
	}

	log.Debug("洗牌之后,得到的牌的数组[%v]", result)
	return result
}


//通过一个index索引来得到一张牌
func InitMjPaiByIndex(index int) *MJPai {
	result := NewMjpai()
	*result.Index = int32(index)
	*result.Des = mjpaiMap[index]
	result.InitByDes()
	return result
}
