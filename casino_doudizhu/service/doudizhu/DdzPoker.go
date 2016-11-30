package doudizhu

import (
	"errors"
	"sort"
	"casino_doudizhu/msg/protogo"
	"casino_doudizhu/msg/funcsInit"
	"strings"
	"casino_common/utils/chessUtils"
	"casino_common/common/log"
	"casino_common/utils/numUtils"
	"casino_common/utils/pokerUtils"
)

var PokerMap map[int32]string

func init() {
	PokerMap = make(map[int32]string, 54)
	PokerMap[0] = ""
	PokerMap[1] = "POKER_spade_3_3"
	PokerMap[2] = "POKER_spade_4_4"
	PokerMap[3] = "POKER_spade_5_5"
	PokerMap[4] = "POKER_spade_6_6"
	PokerMap[5] = "POKER_spade_7_7"
	PokerMap[6] = "POKER_spade_8_8"
	PokerMap[7] = "POKER_spade_9_9"
	PokerMap[8] = "POKER_spade_10_10"
	PokerMap[9] = "POKER_spade_11_J"
	PokerMap[10] = "POKER_spade_12_Q"
	PokerMap[11] = "POKER_spade_13_K"
	PokerMap[12] = "POKER_spade_14_A"
	PokerMap[13] = "POKER_spade_15_2"

	PokerMap [14] = "POKER_heart_3_3"
	PokerMap[15] = "POKER_heart_4_4"
	PokerMap[16] = "POKER_heart_5_5"
	PokerMap[17] = "POKER_heart_6_6"
	PokerMap[18] = "POKER_heart_7_7"
	PokerMap[19] = "POKER_heart_8_8"
	PokerMap[20] = "POKER_heart_9_9"
	PokerMap[21] = "POKER_heart_10_10"
	PokerMap[22] = "POKER_heart_11_J"
	PokerMap[23] = "POKER_heart_12_Q"
	PokerMap[24] = "POKER_heart_13_K"
	PokerMap[25] = "POKER_heart_14_A"
	PokerMap[26] = "POKER_heart_15_2"

	PokerMap[27] = "POKER_club_3_3"
	PokerMap[28] = "POKER_club_4_4"
	PokerMap[29] = "POKER_club_5_5"
	PokerMap[30] = "POKER_club_6_6"
	PokerMap[31] = "POKER_club_7_7"
	PokerMap[32] = "POKER_club_8_8"
	PokerMap[33] = "POKER_club_9_9"
	PokerMap[34] = "POKER_club_10_10"
	PokerMap[35] = "POKER_club_11_J	"
	PokerMap[36] = "POKER_club_12_Q"
	PokerMap[37] = "POKER_club_13_K"
	PokerMap[38] = "POKER_club_14_A"
	PokerMap[39] = "POKER_club_15_2"

	PokerMap[40] = "POKER_diamond_3_3"
	PokerMap[41] = "POKER_diamond_4_4"
	PokerMap[42] = "POKER_diamond_5_5"
	PokerMap[43] = "POKER_diamond_6_6"
	PokerMap[44] = "POKER_diamond_7_7"
	PokerMap[45] = "POKER_diamond_8_8"
	PokerMap[46] = "POKER_diamond_9_9"
	PokerMap[47] = "POKER_diamond_10_10"
	PokerMap[48] = "POKER_diamond_11_J"
	PokerMap[49] = "POKER_diamond_12_Q"
	PokerMap[50] = "POKER_diamond_13_K"
	PokerMap[51] = "POKER_diamond_14_A"
	PokerMap[52] = "POKER_diamond_15_2"

	PokerMap[53] = "POKER_blackjoker_16_JOKER"
	PokerMap[54] = "POKER_redjoker_17_JOKER"
}

//返回值
func parseByIndex(index int32) (int32, string, int32, string, string) {
	var rmapdes string = PokerMap[index]
	sarry := strings.Split(rmapdes, "_")
	log.T("通过index[%v]解析出来的sarry[%v]", index, sarry)
	var pvalue int32 = int32(numUtils.String2Int(sarry[2]))
	var pname string = sarry[3]
	var pflower string = sarry[1]
	return index, rmapdes, pvalue, pflower, pname
}

//返回一张牌
func InitPaiByIndex(index int32) *PPokerPai {
	if index < 0 || index > 54 {
		log.E("非法的牌id，解析出错", index)
		return nil
	}

	_, rmapdes, pvalue, pflower, pname := parseByIndex(index)
	//返回一张需要的牌
	pokerPai := NewPPokerPai()
	*pokerPai.Id = index
	*pokerPai.Name = pname
	*pokerPai.Value = pvalue
	*pokerPai.Flower = Flower2int(pflower)
	*pokerPai.Des = rmapdes
	return pokerPai
}

func Flower2int(f string) int32 {
	if f == "diamond" {
		return pokerUtil.FLOWER_DIAMOND
	} else if f == "club" {
		return pokerUtil.FLOWER_CLUB
	} else if f == "heart" {
		return pokerUtil.FLOWER_HEART
	} else if f == "spade" {
		return pokerUtil.FLOWER_SPADE
	} else if f == "redjoker" {
		return pokerUtil.FLOWER_REDJOKER
	} else if f == "blackjoker" {
		return pokerUtil.FLOWER_BLACKJOKER
	}
	return 0
}

//喜好衣服扑克牌
func XiPai() []*PPokerPai {
	//得到随机的牌的index...
	randIndex := chessUtils.Xipai(1, 54)

	//通过index 得到牌的值
	var pokerPais []*PPokerPai
	for _, i := range randIndex {
		pai := InitPaiByIndex(i)
		if pai == nil {
			//中断...初始化牌的时候出错..
			return nil
		}
		pokerPais = append(pokerPais, pai)
	}

	//返回得到的牌...
	return pokerPais
}

//为POutPokerPais 增加方法

//通过牌来初始化其他的参数
func (out *POutPokerPais) init() error {
	//1,先判断是否有牌
	if out.PokerPais == nil || len(out.PokerPais) <= 0 {
		return errors.New("初始化出的牌失败...没有找到牌。")
	}

	//2，通过牌来进行初始化
	out.sortPais()        //首先对牌进行排序
	out.initTypeAndKeyValue()        //初始化类型和比较值


	return nil
}

type DdzPokerOutList []*PPokerPai;

func (list DdzPokerOutList) Less(i, j  int) bool {
	if list[i].GetValue() < list[j].GetValue() {
		return true
	} else {
		return false
	}
}

// Len 为集合内元素的总数
func (list DdzPokerOutList)Len() int {
	return len(list)
}

// Swap 交换索引为 i 和 j 的元素
func (list DdzPokerOutList) Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
}

//得到牌的张数
func (out *POutPokerPais) getPaiCount() int32 {
	return int32(len(out.PokerPais))
}


//对牌进行排序,左边小，右边大的值进行排序
func (out *POutPokerPais) sortPais() error {
	var list DdzPokerOutList = out.PokerPais
	log.T("befor sort : %v", list)
	sort.Sort(list)        //进行排序
	log.T("after sort : %v", out.PokerPais)

	return nil
}


//初始化类型
func (out *POutPokerPais) initTypeAndKeyValue() error {
	//统计数据
	counts := make([]int32, 17)
	for _, pai := range out.PokerPais {
		counts[pai.GetValue()]++
	}

	var countsLiagnzhang []int32
	var countsSanzhang []int32
	var countsSizhang []int32

	//统计对子,三张
	for _, v := range counts {

		if v == 1 {
			*out.CountYizhang ++
		}

		if v == 2 {
			*out.CountDuizi ++
			countsLiagnzhang = append(countsLiagnzhang, v)
		}

		if v == 3 {
			*out.CountSanzhang ++
			countsSanzhang = append(countsSanzhang, v)
		}

		if v == 4 {
			*out.CountSizhang ++
			countsSizhang = append(countsSizhang, v)
		}
	}

	//判断是否是顺子
	isShunZi := false
	if out.GetCountYizhang() == out.getPaiCount() {
		boolFlag := true
		for k := 0; k < int(out.getPaiCount() - 1); k++ {
			if out.PokerPais[k].GetValue() + 1 != out.PokerPais[k + 1].GetValue() {
				boolFlag = false
				break
			}
		}
		isShunZi = boolFlag
	}

	//飞机带翅膀
	isFeiji := false
	if out.GetCountSanzhang() * 3 == out.getPaiCount() {
		isFeiji = true
	}

	//飞机带单张
	isFeijiChibang := false
	if out.GetCountSanzhang() * 4 == out.getPaiCount() {
		boolFlag := true
		for i := 0; i < int(out.GetCountSanzhang()) - 1; i++ {
			if countsSanzhang[i] + 1 != countsSanzhang[i + 1] {
				boolFlag = false
				break
			}
		}
		isFeijiChibang = boolFlag
	}

	//飞机带对子
	isFeijiDuizi := false
	if out.GetCountSanzhang() * 5 == out.getPaiCount() && out.GetCountSanzhang() == out.GetCountDuizi() {
		boolFlag := true
		for i := 0; i < int(out.GetCountSanzhang() - 1); i++ {
			if countsSanzhang[i] + 1 != countsSanzhang[i + 1] {
				boolFlag = false
				break
			}
		}
		isFeijiDuizi = boolFlag
	}

	isSiDaiLiangDui := false
	if out.getPaiCount() == 8 && out.GetCountSizhang() == 1 && out.GetCountDuizi() == 2 {
		isSiDaiLiangDui = true        //四个带两队
	}

	isLianDui := false
	if out.GetCountDuizi() * 2 == out.getPaiCount() {
		boolFlag := true
		for i := 0; i < int(out.GetCountDuizi() - 1); i++ {
			if countsLiagnzhang[i] + 1 != countsLiagnzhang[i + 1] {
				boolFlag = false
				break
			}
		}
		isLianDui = boolFlag
	}


	//判断是否是单张
	if out.getPaiCount() == 1 {
		*out.Type = int32(ddzproto.DdzPaiType_SINGLECARD) //单张
		*out.KeyValue = out.GetPokerPais()[0].GetValue()
	} else if out.getPaiCount() == 2 {
		if out.GetCountDuizi() == 1 {
			//这里需要判断是否是王炸
			if out.GetPokerPais()[0].GetValue() == -1 {
				//判断是否是王炸
				*out.Type = int32(ddzproto.DdzPaiType_SUPERBOMB) //王炸弹
				*out.KeyValue = out.GetPokerPais()[0].GetValue()
			} else {
				*out.Type = int32(ddzproto.DdzPaiType_DOUBLECARD) // 对子
				*out.KeyValue = out.GetPokerPais()[0].GetValue()
			}
		} else {
			//error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)   //错误牌型
		}
	} else if out.getPaiCount() == 3 {
		if out.GetCountSanzhang() == 1 {
			*out.Type = int32(ddzproto.DdzPaiType_THREECARD)                //三张
			*out.KeyValue = countsSanzhang[0]
		} else {
			//error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)   //错误牌型
		}
	} else if out.getPaiCount() == 4 {
		if out.GetCountSizhang() == 1 {
			*out.Type = int32(ddzproto.DdzPaiType_BOMBCARD)        // 炸弹
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSanzhang() == 1 && out.GetCountYizhang() == 1 {
			*out.Type = int32(ddzproto.DdzPaiType_THREEONECARD)        // 三带一
			*out.KeyValue = countsSanzhang[0]
		} else {
			//	error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)   //错误牌型

		}
	} else if out.getPaiCount() == 5 {
		if isShunZi {
			*out.Type = int32(ddzproto.DdzPaiType_CONNECTCARD)        // 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSanzhang() == 1 && out.GetCountDuizi() == 1 {
			*out.Type = int32(ddzproto.DdzPaiType_THREETWOCARD)        //三带二
			*out.KeyValue = countsSanzhang[0]
		} else {
			//error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)   //错误牌型
		}

	} else if out.getPaiCount() == 6 {
		if isShunZi {
			*out.Type = int32(ddzproto.DdzPaiType_CONNECTCARD)        // 顺子
		} else if out.GetCountSizhang() == 1 && out.GetCountYizhang() == 2 {
			*out.Type = int32(ddzproto.DdzPaiType_BOMBTWOCARD)      // 四带二个单张
		} else {
			//error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)  //错误牌型
		}

	} else if out.getPaiCount() == 8 {
		if isShunZi {
			*out.Type = int32(ddzproto.DdzPaiType_CONNECTCARD)        // 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSizhang() == 1 && out.GetCountDuizi() == 2 {
			*out.Type = int32(ddzproto.DdzPaiType_BOMBTWOOOCARD)       // 四带两队
			*out.KeyValue = countsSizhang[0]
		} else {
			//error
			*out.Type = int32(ddzproto.DdzPaiType_ERRORCARD)   //错误牌型
		}

	} else {
		if isShunZi {
			*out.Type = int32(ddzproto.DdzPaiType_CONNECTCARD)        // 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if isFeiji {
			*out.Type = int32(ddzproto.DdzPaiType_AIRCRAFTCARD)      //飞机
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if isFeijiChibang {
			*out.Type = int32(ddzproto.DdzPaiType_AIRCRAFTSINGLECARD) // 飞机带翅膀
			*out.KeyValue = countsSanzhang[0]
		} else if isFeijiDuizi {
			*out.Type = int32(ddzproto.DdzPaiType_AIRCRAFTDOUBLECARD) // 飞机带翅膀
			*out.KeyValue = countsSanzhang[0]
		} else if isLianDui {
			*out.Type = int32(ddzproto.DdzPaiType_COMPANYCARD) //连队
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if isSiDaiLiangDui {
			*out.Type = int32(ddzproto.DdzPaiType_FOURWITHONEDOUBLE) //连队
			*out.KeyValue = countsSizhang[0]   //比较值
		}
	}

	return nil;
}

//比较两幅牌的大小
func (out *POutPokerPais)  GT(outb *POutPokerPais) (bool, error) {
	if out.GetType() == outb.GetType() {
		return out.GetKeyValue() > outb.GetKeyValue(), nil
	} else {
		//比较类型不同的情况
		if out.GetIsBomb() || outb.GetIsBomb() {
			return out.GetIsBomb(), nil
		} else {
			log.E("牌型[%v]和牌型[%v] 无法比较...", out.GetType(), outb.GetType())
			return false, errors.New("无法比较，类型有错误..")
		}
	}

}

func (out *POutPokerPais) GetClientPokers() []*ddzproto.Poker {
	return ServerPoker2ClienPoker(out.PokerPais)
}

//得到牌
func GetOutPais(outcards []*ddzproto.Poker, userId uint32) *POutPokerPais {
	ret := NewPOutPokerPais()
	for _, c := range outcards {
		ret.PokerPais = append(ret.PokerPais, InitPaiByIndex(c.GetId()))
	}
	*ret.UserId = userId
	ret.init()
	return ret
}

func ( p *PPokerPai) GetClientPoker() *ddzproto.Poker {
	ret := newProto.NewPoker()
	*ret.Id = p.GetId()
	*ret.Suit = p.GetSuit()
	*ret.Num = p.GetValue()
	return ret
}

func (p *PPokerPai) GetSuit() ddzproto.PokerColor {
	if p.GetFlower() == pokerUtil.FLOWER_DIAMOND {
		return ddzproto.PokerColor_FANGKUAI
	} else if p.GetFlower() == pokerUtil.FLOWER_CLUB {
		return ddzproto.PokerColor_MEIHUA
	} else if p.GetFlower() == pokerUtil.FLOWER_HEART {
		return ddzproto.PokerColor_HONGTAO
	} else if p.GetFlower() == pokerUtil.FLOWER_SPADE {
		return ddzproto.PokerColor_HEITAO
	} else if p.GetFlower() == pokerUtil.FLOWER_BLACKJOKER {
		return ddzproto.PokerColor_BLACKBIGJOKER
	} else {
		return ddzproto.PokerColor_REDJOKER
	}
}

func (p *PPokerPai) GetLogDes() string {
	suit, _ := numUtils.Int2String(p.GetId())
	switch p.GetSuit() {
	case ddzproto.PokerColor_FANGKUAI:
		suit = "方块"
	case ddzproto.PokerColor_BLACKBIGJOKER:
		suit = "小王"
	case ddzproto.PokerColor_HEITAO:
		suit = "黑桃"
	case ddzproto.PokerColor_HONGTAO:
		suit = "红桃"
	case ddzproto.PokerColor_REDJOKER:
		suit = "大王"
	case ddzproto.PokerColor_MEIHUA:
		suit = "梅花"
	default:
	}
	suit += p.GetName()
	return suit
}

//把服务端的牌转化成客户端需要的牌
func ServerPoker2ClienPoker(sps []*PPokerPai) []*ddzproto.Poker {
	if sps == nil || len(sps) == 0 {
		return nil
	}

	var ret []*ddzproto.Poker
	for _, p := range sps {
		ret = append(ret, p.GetClientPoker())
	}
	return ret
}