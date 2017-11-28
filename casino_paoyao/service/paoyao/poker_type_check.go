package paoyao

import (
	"casino_common/proto/ddproto"
	"casino_common/common/log"
)

//是否是单牌
func IsDanpai(pais []*ddproto.CommonSrvPokerPai) bool {
	if len(pais) == 1 {
		return true
	}
	return false
}

//是否是对子
func IsDuipai(pais []*ddproto.CommonSrvPokerPai) bool {
	if len(pais) == 2 && pais[0].GetValue() == pais[1].GetValue() {
		return true
	}
	return false
}

//是否是单顺
func IsDanshun(pais []*ddproto.CommonSrvPokerPai) bool {
	if len(pais) < 3 {
		return false
	}

	for i,p := range pais {
		//单牌必须为4-A
		if p.GetValue() > 10 {
			return false
		}
		if i == 0 {
			continue
		}
		if p.GetValue() - pais[i-1].GetValue() != 1 {
			return false
		}
	}
	return true
}

//是否为双顺
func IsShuangshun(pais []*ddproto.CommonSrvPokerPai) bool {
	//必须为长度大于4的偶数牌
	if len(pais) < 6 || len(pais)%2 == 1 {
		return false
	}

	for i,p := range pais {
		//单牌必须为4-A
		if p.GetValue() > 10 {
			return false
		}

		if i == 0 {
			continue
		}

		if p.GetValue() - pais[i-1].GetValue() != int32((i+1)%2) {
			return false
		}
	}
	return true
}

//是否为路，三路、四路、五路、六路、七路、八路
func IsLu(pais []*ddproto.CommonSrvPokerPai) bool {
	//长度必须大于等于3
	if len(pais) < 3 {
		return false
	}

	for i,p := range pais {
		if i == 0 {
			continue
		}
		//必须每张牌相同
		if p.GetValue() != pais[i-1].GetValue() {
			return false
		}
	}


	return true
}

//是否为幺，小妖、中夭、老幺、老老幺
func IsYao(pais []*ddproto.CommonSrvPokerPai) bool {
	map_num := map[string]int32 {}
	for _,p := range pais {
		if p.GetName() != "A" && p.GetName() != "4" {
			return false
		}
		if _,ok := map_num[p.GetName()];ok {
			map_num[p.GetName()] += 1
		} else {
			map_num[p.GetName()] = 1
		}
	}
	if len(map_num) != 2 || map_num["A"] > 1 || map_num["4"] < 2 || map_num["4"] > 8 {
		return false
	}
	return true
}


//解析牌型
func ParsePokerType(pais []*ddproto.CommonSrvPokerPai) ddproto.PaoyaoEnumPokerType {
	switch {
	case IsDanpai(pais):
		return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DAN_PAI
	case IsDuipai(pais):
		return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DUI_PAI
	case IsDanshun(pais):
		return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DAN_SHUN
	case IsShuangshun(pais):
		return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_SHUANG_SHUN
	case IsLu(pais):
		switch len(pais) {
		case 3:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_3_LU
		case 4:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_4_LU
		case 5:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_5_LU
		case 6:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_6_LU
		case 7:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_7_LU
		case 8:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_8_LU
		default:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER
		}
	case IsYao(pais):
		switch len(pais) {
		case 3:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_XIAO_YAO
		case 4:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_ZHONG_YAO
		case 5:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_LAO_YAO
		default:
			return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_LAOLAO_YAO
		}
	default:
		return ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER
	}
}

//将客户端牌转为服务器牌
func ParseOutPai(client_pai []int32) *ddproto.PaoyaoSrvPoker  {
	out_pai := &ddproto.PaoyaoSrvPoker{
		Pais: []*ddproto.CommonSrvPokerPai{},
		Type: ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER.Enum(),
	}

	for _,pai_id := range client_pai {
		if pai_id < 0 || pai_id > 107 {
			log.E("ParseOutPai %v id=%v is nil.", client_pai, pai_id)
			continue
		}

		pai_item := PokerList[int32(pai_id)]
		out_pai.Pais = append(out_pai.Pais, pai_item)
	}

	//牌值排序
	out_pai.Pais = SortPokers(out_pai.Pais)

	if len(client_pai) != 27 {
		out_pai.Type = ParsePokerType(out_pai.Pais).Enum()
	}

	return out_pai
}

//牌型比较大小
func IsBigThanPoker(paisA, paisB *ddproto.PaoyaoSrvPoker) bool {

	AisDanpai := paisA.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DAN_PAI
	AisDanshun := paisA.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DAN_SHUN

	BisDanshun := paisB.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_DAN_SHUN
	BisShuangshun := paisB.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_SHUANG_SHUN

	//过滤非法牌型
	if paisA.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER || paisB.GetType() == ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER {
		log.E("出牌错误：牌型解析出错 %v %v", paisA, paisB)
		return false
	}

	//牌型一样的情况
	if paisA.GetType() == paisB.GetType() {
		//如果两牌型相同，则先比较牌长度，再比较第一张牌大小
		//适用于：单牌 对牌 单顺 双顺 路 幺
		if len(paisB.Pais) > len(paisA.Pais) {
			return true
		}else if len(paisB.Pais) == len(paisA.Pais) {
			//相同长度，则比较第一张牌大小
			if paisB.Pais[0].GetValue() > paisA.Pais[0].GetValue() {
				return true
			}else {
				log.E("出牌错误：相同长度比第一张牌大小 %v %v", paisA, paisB)
				return false
			}
		}else {
			log.E("出牌错误：出的牌比上家牌长度小 %v %v", paisA, paisB)
			return false
		}
	}

	//判断 单牌 和 对子 的特殊情况
	if AisDanpai && (BisDanshun || BisShuangshun) {
		log.E("出牌错误：单牌不能与对子比牌 %v %v", paisA, paisB)
		return false
	}

	//判断 单顺和双顺 的特殊情况
	if AisDanshun && BisShuangshun {
		if len(paisB.Pais) < len(paisA.Pais)*2 {
			log.E("出牌错误：双顺 顺子数比单顺顺子数小 %v %v", paisA, paisB)
			return false
		}
	}

	//最后比较牌型大小
	if paisB.GetType() > paisA.GetType() {
		return true
	}

	log.E("出牌错误：出牌牌型小于上家牌型 %v %v", paisA, paisB)
	return false
}
