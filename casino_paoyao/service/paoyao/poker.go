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
var PokerList = make([]*ddproto.CommonSrvPokerPai, 54)

func init() {
	PokerMap = make(map[int32]string, 54)
	//POKER 花色 牌值 牌型
	PokerMap[0] = "POKER_spade_13_3"  //黑桃
	PokerMap[1] = "POKER_spade_1_4"
	PokerMap[2] = "POKER_spade_2_5"
	PokerMap[3] = "POKER_spade_3_6"
	PokerMap[4] = "POKER_spade_4_7"
	PokerMap[5] = "POKER_spade_5_8"
	PokerMap[6] = "POKER_spade_6_9"
	PokerMap[7] = "POKER_spade_7_10"
	PokerMap[8] = "POKER_spade_8_J"
	PokerMap[9] = "POKER_spade_9_Q"
	PokerMap[10] = "POKER_spade_10_K"
	PokerMap[11] = "POKER_spade_11_A"
	PokerMap[12] = "POKER_spade_12_2"

	PokerMap[13] = "POKER_heart_13_3"  //红桃
	PokerMap[14] = "POKER_heart_1_4"
	PokerMap[15] = "POKER_heart_2_5"
	PokerMap[16] = "POKER_heart_3_6"
	PokerMap[17] = "POKER_heart_4_7"
	PokerMap[18] = "POKER_heart_5_8"
	PokerMap[19] = "POKER_heart_6_9"
	PokerMap[20] = "POKER_heart_7_10"
	PokerMap[21] = "POKER_heart_8_J"
	PokerMap[22] = "POKER_heart_9_Q"
	PokerMap[23] = "POKER_heart_10_K"
	PokerMap[24] = "POKER_heart_11_A"
	PokerMap[25] = "POKER_heart_12_2"

	PokerMap[26] = "POKER_club_13_3"  //梅花
	PokerMap[27] = "POKER_club_1_4"
	PokerMap[28] = "POKER_club_2_5"
	PokerMap[29] = "POKER_club_3_6"
	PokerMap[30] = "POKER_club_4_7"
	PokerMap[31] = "POKER_club_5_8"
	PokerMap[32] = "POKER_club_6_9"
	PokerMap[33] = "POKER_club_7_10"
	PokerMap[34] = "POKER_club_8_J"
	PokerMap[35] = "POKER_club_9_Q"
	PokerMap[36] = "POKER_club_10_K"
	PokerMap[37] = "POKER_club_11_A"
	PokerMap[38] = "POKER_club_12_2"

	PokerMap[39] = "POKER_diamond_13_3"  //方块
	PokerMap[40] = "POKER_diamond_1_4"
	PokerMap[41] = "POKER_diamond_2_5"
	PokerMap[42] = "POKER_diamond_3_6"
	PokerMap[43] = "POKER_diamond_4_7"
	PokerMap[44] = "POKER_diamond_5_8"
	PokerMap[45] = "POKER_diamond_6_9"
	PokerMap[46] = "POKER_diamond_7_10"
	PokerMap[47] = "POKER_diamond_8_J"
	PokerMap[48] = "POKER_diamond_9_Q"
	PokerMap[49] = "POKER_diamond_10_K"
	PokerMap[50] = "POKER_diamond_11_A"
	PokerMap[51] = "POKER_diamond_12_2"

	PokerMap[52] = "POKER_blackjoker_14_JOKER"
	PokerMap[53] = "POKER_redjoker_15_JOKER"

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
	if index < 0 || index > 53 {
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
	if srv == nil {
		log.E("poker nil")
		return nil
	}
	server_pais := srv.Pais
	clent_poker := &ddproto.PaoyaoClientPoker{
		Pais: []int32{},
	}

	for _,s_pai := range server_pais {
		clent_poker.Pais = append(clent_poker.Pais, s_pai.GetId())
	}

	return clent_poker
}
