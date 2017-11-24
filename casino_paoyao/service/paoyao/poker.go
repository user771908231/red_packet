package paoyao

import (
	"casino_common/proto/ddproto"
	"strings"
	"casino_common/utils/numUtils"
	"casino_common/utils/pokerUtils"
	"casino_common/common/log"
)

var PokerMap map[int32]string

//一副扑克
var PokerList = make([]*ddproto.CommonSrvPokerPai, 108)

func init() {
	PokerMap = make(map[int32]string, 108)
	//POKER 花色 牌值 牌型
	PokerMap[0] = "POKER_spade_1_4"  //黑桃
	PokerMap[1] = "POKER_spade_2_5"
	PokerMap[2] = "POKER_spade_3_6"
	PokerMap[3] = "POKER_spade_4_7"
	PokerMap[4] = "POKER_spade_5_8"
	PokerMap[5] = "POKER_spade_6_9"
	PokerMap[6] = "POKER_spade_7_10"
	PokerMap[7] = "POKER_spade_8_J"
	PokerMap[8] = "POKER_spade_9_Q"
	PokerMap[9] = "POKER_spade_10_K"
	PokerMap[10] = "POKER_spade_11_A"
	PokerMap[11] = "POKER_spade_12_2"
	PokerMap[12] = "POKER_spade_13_3"

	PokerMap[13] = "POKER_heart_1_4"  //红桃
	PokerMap[14] = "POKER_heart_2_5"
	PokerMap[15] = "POKER_heart_3_6"
	PokerMap[16] = "POKER_heart_4_7"
	PokerMap[17] = "POKER_heart_5_8"
	PokerMap[18] = "POKER_heart_6_9"
	PokerMap[19] = "POKER_heart_7_10"
	PokerMap[20] = "POKER_heart_8_J"
	PokerMap[21] = "POKER_heart_9_Q"
	PokerMap[22] = "POKER_heart_10_K"
	PokerMap[23] = "POKER_heart_11_A"
	PokerMap[24] = "POKER_heart_12_2"
	PokerMap[25] = "POKER_heart_13_3"

	PokerMap[26] = "POKER_club_1_4"  //梅花
	PokerMap[27] = "POKER_club_2_5"
	PokerMap[28] = "POKER_club_3_6"
	PokerMap[39] = "POKER_club_4_7"
	PokerMap[30] = "POKER_club_5_8"
	PokerMap[31] = "POKER_club_6_9"
	PokerMap[32] = "POKER_club_7_10"
	PokerMap[33] = "POKER_club_8_J"
	PokerMap[34] = "POKER_club_9_Q"
	PokerMap[35] = "POKER_club_10_K"
	PokerMap[36] = "POKER_club_11_A"
	PokerMap[37] = "POKER_club_12_2"
	PokerMap[38] = "POKER_club_13_3"

	PokerMap[39] = "POKER_diamond_1_4"  //方块
	PokerMap[40] = "POKER_diamond_2_5"
	PokerMap[41] = "POKER_diamond_3_6"
	PokerMap[42] = "POKER_diamond_4_7"
	PokerMap[43] = "POKER_diamond_5_8"
	PokerMap[44] = "POKER_diamond_6_9"
	PokerMap[45] = "POKER_diamond_7_10"
	PokerMap[46] = "POKER_diamond_8_J"
	PokerMap[47] = "POKER_diamond_9_Q"
	PokerMap[48] = "POKER_diamond_10_K"
	PokerMap[49] = "POKER_diamond_11_A"
	PokerMap[50] = "POKER_diamond_12_2"
	PokerMap[51] = "POKER_diamond_13_3"

	PokerMap[52] = "POKER_blackjoker_14_JOKER"
	PokerMap[53] = "POKER_redjoker_15_JOKER"

	PokerMap[54] = "POKER_spade_1_4"  //黑桃
	PokerMap[55] = "POKER_spade_2_5"
	PokerMap[56] = "POKER_spade_3_6"
	PokerMap[57] = "POKER_spade_4_7"
	PokerMap[58] = "POKER_spade_5_8"
	PokerMap[59] = "POKER_spade_6_9"
	PokerMap[60] = "POKER_spade_7_10"
	PokerMap[61] = "POKER_spade_8_J"
	PokerMap[62] = "POKER_spade_9_Q"
	PokerMap[63] = "POKER_spade_10_K"
	PokerMap[64] = "POKER_spade_11_A"
	PokerMap[65] = "POKER_spade_12_2"
	PokerMap[66] = "POKER_spade_13_3"

	PokerMap[67] = "POKER_heart_1_4"  //红桃
	PokerMap[68] = "POKER_heart_2_5"
	PokerMap[69] = "POKER_heart_3_6"
	PokerMap[70] = "POKER_heart_4_7"
	PokerMap[71] = "POKER_heart_5_8"
	PokerMap[72] = "POKER_heart_6_9"
	PokerMap[73] = "POKER_heart_7_10"
	PokerMap[74] = "POKER_heart_8_J"
	PokerMap[75] = "POKER_heart_9_Q"
	PokerMap[76] = "POKER_heart_10_K"
	PokerMap[77] = "POKER_heart_11_A"
	PokerMap[78] = "POKER_heart_12_2"
	PokerMap[79] = "POKER_heart_13_3"

	PokerMap[80] = "POKER_club_1_4"  //梅花
	PokerMap[81] = "POKER_club_2_5"
	PokerMap[82] = "POKER_club_3_6"
	PokerMap[83] = "POKER_club_4_7"
	PokerMap[84] = "POKER_club_5_8"
	PokerMap[85] = "POKER_club_6_9"
	PokerMap[86] = "POKER_club_7_10"
	PokerMap[87] = "POKER_club_8_J"
	PokerMap[88] = "POKER_club_9_Q"
	PokerMap[89] = "POKER_club_10_K"
	PokerMap[90] = "POKER_club_11_A"
	PokerMap[91] = "POKER_club_12_2"
	PokerMap[92] = "POKER_club_13_3"

	PokerMap[93] = "POKER_diamond_1_4"  //方块
	PokerMap[94] = "POKER_diamond_2_5"
	PokerMap[95] = "POKER_diamond_3_6"
	PokerMap[96] = "POKER_diamond_4_7"
	PokerMap[97] = "POKER_diamond_5_8"
	PokerMap[98] = "POKER_diamond_6_9"
	PokerMap[99] = "POKER_diamond_7_10"
	PokerMap[100] = "POKER_diamond_8_J"
	PokerMap[101] = "POKER_diamond_9_Q"
	PokerMap[102] = "POKER_diamond_10_K"
	PokerMap[103] = "POKER_diamond_11_A"
	PokerMap[104] = "POKER_diamond_12_2"
	PokerMap[105] = "POKER_diamond_13_3"

	PokerMap[106] = "POKER_blackjoker_14_JOKER"
	PokerMap[107] = "POKER_redjoker_15_JOKER"

	//初始化一副扑克
	for k,_ := range PokerMap {
		new_pai := InitPaiByIndex(k)
		PokerList[k]= new_pai
	}

}

//初始化--单张扑克牌
func newPPokerPai() *ddproto.CommonSrvPokerPai {
	pai := new(ddproto.CommonSrvPokerPai)
	pai.Des = new(string)
	pai.Flower = new(int32)
	pai.Id = new(int32)
	pai.Name = new(string)
	pai.Value = new(int32)
	return pai
}

//返回一张牌
func InitPaiByIndex(index int32) *ddproto.CommonSrvPokerPai {
	if index < 0 || index > 107 {
		return nil
	}

	_, rmapdes, pvalue, pflower, pname := parseByIndex(index)
	//返回一张需要的牌
	pokerPai := newPPokerPai()
	*pokerPai.Id = index
	*pokerPai.Name = pname
	*pokerPai.Value = pvalue
	*pokerPai.Flower = flower2int(pflower)
	*pokerPai.Des = rmapdes
	return pokerPai
}

func flower2int(f string) int32 {
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

//解析单张
func parseByIndex(index int32) (int32, string, int32, string, string) {
	var rmapdes string = PokerMap[index]
	sarry := strings.Split(rmapdes, "_")
	//log.T("通过index[%v]解析出来的sarry[%v]", index, sarry)
	var pvalue int32 = int32(numUtils.String2Int(sarry[2]))
	var pname string = sarry[3]
	var pflower string = sarry[1]
	return index, sarry[1], pvalue, pflower, pname
}

//转换成客户端花色
func flower2clientFlower(f string) ddproto.CommonEnumPokerColor {
	if f == "diamond" {
		return ddproto.CommonEnumPokerColor_FANGKUAI
	} else if f == "club" {
		return ddproto.CommonEnumPokerColor_MEIHUA
	} else if f == "heart" {
		return ddproto.CommonEnumPokerColor_HONGTAO
	} else if f == "spade" {
		return ddproto.CommonEnumPokerColor_HEITAO
	} else if f == "redjoker" {
		return ddproto.CommonEnumPokerColor_REDJOKER
	} else if f == "blackjoker" {
		return ddproto.CommonEnumPokerColor_BLACKBIGJOKER
	}
	return ddproto.CommonEnumPokerColor_FANGKUAI
}

//将服务器的牌型转换成客户端的
func GetClientPoker(srv *ddproto.PaoyaoSrvPoker) *ddproto.PaoyaoClientPoker {
	clent_poker := &ddproto.PaoyaoClientPoker{
		Pais: []int32{},
		PokerType: ddproto.PaoyaoEnumPokerType(srv.GetType()).Enum(),
	}
	if srv == nil {
		log.E("poker nil")
		return clent_poker
	}

	for _,s_pai := range srv.Pais {
		clent_poker.Pais = append(clent_poker.Pais, s_pai.GetId())
	}

	return clent_poker
}

//牌型合法性验证
func GetLostPokersByOutpai(out_pai *ddproto.PaoyaoSrvPoker, all_poker *ddproto.PaoyaoSrvPoker) (*ddproto.PaoyaoSrvPoker) {
	lost_poker := &ddproto.PaoyaoSrvPoker{
		Pais: []*ddproto.CommonSrvPokerPai{},
		Type: ddproto.PaoyaoEnumPokerType_PAOYAO_POKER_TYPE_OTHER.Enum(),
	}

	for _, pa := range all_poker.Pais {
		has_poker := false
		for _, pi := range out_pai.Pais {
			if pa.GetId() == pi.GetId() {
				has_poker = true
				break
			}
		}
		if has_poker == false {
			lost_poker.Pais = append(lost_poker.Pais, pa)
		}
	}

	return lost_poker
}

//计算出牌的分数
func GetOutpaiScore(out_pai *ddproto.PaoyaoSrvPoker) (score int32) {
	for _,p := range out_pai.Pais {
		switch p.GetName() {
		case "5":
			score += 5
		case "10", "K":
			score += 10
		}
	}
	return
}

//牌面排序
func SortPokers(src_pokers []*ddproto.CommonSrvPokerPai) []*ddproto.CommonSrvPokerPai {
	if len(src_pokers) < 2 {
		return src_pokers
	}
	for i:=0;i<len(src_pokers)-1;i++ {
		for j:=i+1;j<len(src_pokers);j++ {
			if src_pokers[i].GetValue() > src_pokers[j].GetValue() {
				tmp := src_pokers[i]
				src_pokers[i] = src_pokers[j]
				src_pokers[j] = tmp
			}
		}
	}
	return src_pokers
}
