package paosangong

import (
	"casino_common/proto/ddproto"
	"strings"
	"casino_common/utils/numUtils"
	"casino_common/utils/pokerUtils"
	//"log"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
	"fmt"
)

var PokerMap map[int32]string

//一副扑克
var PokerList = make([]*ddproto.CommonSrvPokerPai, 52)

func init() {
	PokerMap = make(map[int32]string, 52)
	PokerMap[0] = "POKER_spade_3_3"  //黑桃
	PokerMap[1] = "POKER_spade_4_4"
	PokerMap[2] = "POKER_spade_5_5"
	PokerMap[3] = "POKER_spade_6_6"
	PokerMap[4] = "POKER_spade_7_7"
	PokerMap[5] = "POKER_spade_8_8"
	PokerMap[6] = "POKER_spade_9_9"
	PokerMap[7] = "POKER_spade_10_10"
	PokerMap[8] = "POKER_spade_10_J"
	PokerMap[9] = "POKER_spade_10_Q"
	PokerMap[10] = "POKER_spade_10_K"
	PokerMap[11] = "POKER_spade_1_A"
	PokerMap[12] = "POKER_spade_2_2"

	PokerMap[13] = "POKER_heart_3_3"  //红桃
	PokerMap[14] = "POKER_heart_4_4"
	PokerMap[15] = "POKER_heart_5_5"
	PokerMap[16] = "POKER_heart_6_6"
	PokerMap[17] = "POKER_heart_7_7"
	PokerMap[18] = "POKER_heart_8_8"
	PokerMap[19] = "POKER_heart_9_9"
	PokerMap[20] = "POKER_heart_10_10"
	PokerMap[21] = "POKER_heart_10_J"
	PokerMap[22] = "POKER_heart_10_Q"
	PokerMap[23] = "POKER_heart_10_K"
	PokerMap[24] = "POKER_heart_1_A"
	PokerMap[25] = "POKER_heart_2_2"

	PokerMap[26] = "POKER_club_3_3"  //梅花
	PokerMap[27] = "POKER_club_4_4"
	PokerMap[28] = "POKER_club_5_5"
	PokerMap[29] = "POKER_club_6_6"
	PokerMap[30] = "POKER_club_7_7"
	PokerMap[31] = "POKER_club_8_8"
	PokerMap[32] = "POKER_club_9_9"
	PokerMap[33] = "POKER_club_10_10"
	PokerMap[34] = "POKER_club_10_J"
	PokerMap[35] = "POKER_club_10_Q"
	PokerMap[36] = "POKER_club_10_K"
	PokerMap[37] = "POKER_club_1_A"
	PokerMap[38] = "POKER_club_2_2"

	PokerMap[39] = "POKER_diamond_3_3"  //方块
	PokerMap[40] = "POKER_diamond_4_4"
	PokerMap[41] = "POKER_diamond_5_5"
	PokerMap[42] = "POKER_diamond_6_6"
	PokerMap[43] = "POKER_diamond_7_7"
	PokerMap[44] = "POKER_diamond_8_8"
	PokerMap[45] = "POKER_diamond_9_9"
	PokerMap[46] = "POKER_diamond_10_10"
	PokerMap[47] = "POKER_diamond_10_J"
	PokerMap[48] = "POKER_diamond_10_Q"
	PokerMap[49] = "POKER_diamond_10_K"
	PokerMap[50] = "POKER_diamond_1_A"
	PokerMap[51] = "POKER_diamond_2_2"

	//PokerMap[52] = "POKER_blackjoker_16_JOKER"
	//PokerMap[53] = "POKER_redjoker_17_JOKER"

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
	if index < 0 || index > 51 {
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

//解析牛牛牌型 是否花样玩法
func ParseNiuPoker(pais []*ddproto.CommonSrvPokerPai, option *ddproto.NiuniuDeskOption) *ddproto.NiuniuSrvPoker {
	niu_poker := &ddproto.NiuniuSrvPoker{
		Pais: pais,
		Type: ddproto.NiuniuEnum_PokerType_NO_NIU.Enum(),
		SelectedId: []int32{},
	}

	// 为顺子牛（一条龙）
	if option.GetHasShunziNiu() {
		if isYiTiaoLong(pais) {
			//同花顺-小帅牛
			if isTongHua(pais) {
				niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_XIAOSHUAI_NIU.Enum()
			}else {
				//顺子牛
				niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_YI_TIAO_LONG.Enum()
			}
			return niu_poker
		}
	}

	//是否为炸弹牛
	if option.GetHasZhadanNiu(){
		name_map := map[string]int{}
		for _,pai := range pais {
			name := pai.GetName()
			if _,ok := name_map[name]; ok {
				name_map[name] += 1
			}else {
				name_map[name] = 1
			}
			if name_map[name] == 4 {
				// 为炸弹
				niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_ZHA_DAN.Enum()
				return niu_poker
			}
		}
	}

	//是否为同花牛
	if option.GetHasTonghuaNiu() {
		if isTongHua(pais) {
			niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_TONGHUA_NIU.Enum()
			return niu_poker
		}
	}

	//葫芦牛
	if option.GetHasHuluNiu() {
		if isHulu(pais) {
			niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_TONGHUA_NIU.Enum()
			return niu_poker
		}
	}

	index_arr := [][]int32{
		[]int32{0, 1, 2},
		[]int32{2, 3, 4},
		[]int32{0, 3, 4},
		[]int32{0, 1, 4},
		[]int32{0, 2, 4},
		[]int32{1, 2, 3},
		[]int32{0, 1, 3},
		[]int32{1, 2, 4},
		[]int32{0, 2, 3},
		[]int32{1, 3, 4},
	}

	niu_num := map[int]int{}
	for niu_index,arr := range index_arr {
		//是否有牛
		if (pais[arr[0]].GetValue() + pais[arr[1]].GetValue() + pais[arr[2]].GetValue())%10 == 0 {
			niu_num[niu_index] = int(( pais[0].GetValue() + pais[1].GetValue() + pais[2].GetValue() + pais[3].GetValue() + pais[4].GetValue() )%10)
		}
	}

	max_index := 0
	max_niu := 0
	for tmp_index,tmp_niu := range niu_num {
		if tmp_niu > max_niu {
			max_niu = tmp_niu
			max_index = tmp_index
		}
	}
	niu_poker.SelectedId = index_arr[max_index]
	for tmp_index,_ := range niu_num {
		if tmp_index != max_index {
			niu_poker.SelectedId = append(niu_poker.SelectedId, index_arr[tmp_index]...)
		}
	}

	if len(niu_num) > 1 {
		fmt.Println(niu_num)
	}

	//银牛
	num_yin := 0
	//金牛
	num_jin := 0
	//五小牛
	num_xiao_niu := 0
	var val_count int32 = 0
	for _,pai := range pais {
		name := pai.GetName()
		val := pai.GetValue()
		if name == "J" || name == "Q" || name == "K" {
			num_jin++
		}
		if val >= 10 {
			num_yin++
		}
		if val <= 5 {
			num_xiao_niu++
		}
		val_count += val
	}

	//有牛
	if len(niu_num) > 0 {

		if max_niu > 0 {
			//牛1-9
			niu_poker.Type = ddproto.NiuniuEnum_PokerType(1+max_niu).Enum()
			return niu_poker
		}else {
			//牛牛
			//是否为花样玩法
			//if num_jin == 5 {
			//	// 金牛
			//	niu_poker.Type = ddproto.NiuniuEnum_PokerType_JIN_NIU.Enum()
			//	return niu_poker
			//}

			// 银牛（五花牛）
			if option.GetHasWuhuaNiu() {
				if num_yin == 5 {
					niu_poker.Type = ddproto.NiuniuEnum_PokerType_YIN_NIU.Enum()
					return niu_poker
				}
			}

			// 牛牛
			niu_poker.Type = ddproto.NiuniuEnum_PokerType_NIU_NIU.Enum()
			return niu_poker
		}

	}else {
		//if num_xiao_niu == 5 && val_count <= 10 {
		//	// 五小牛
		//	niu_poker.Type = ddproto.NiuniuEnum_PokerType_WU_XIAO_NIU.Enum()
		//	return niu_poker
		//}
		//没牛
		niu_poker.Type = ddproto.NiuniuEnum_PokerType_NO_NIU.Enum()
		return niu_poker
	}

	return niu_poker
}

//是否一条龙
func isYiTiaoLong(pais []*ddproto.CommonSrvPokerPai) bool {
	arr_src := []int{}
	for _,pai := range pais {
		val := pai.GetValue()
		name := pai.GetName()
		switch name {
		case "J":
			val = 11
		case "Q":
			val = 12
		case "K":
			val = 13
		}
		arr_src = append(arr_src, int(val))
	}
	//log.Println(1, arr_src)
	//排序
	for i:=0;i<4;i++ {
		for j:=i+1;j<5;j++ {
			if arr_src[i] > arr_src[j] {
				tmp := arr_src[i]
				arr_src[i] = arr_src[j]
				arr_src[j] = tmp
			}
		}
	}
	//log.Println(2, arr_src)

	sum := 0
	for i, val := range arr_src {
		if i > 0 {
			if val - arr_src[i-1] == 1 {
				sum++
			}
		}
	}

	if sum == 4 {
		return true
	}
	return false
}

//是否为同花牛
func isTongHua(pais []*ddproto.CommonSrvPokerPai) bool {
	if len(pais) == 0 {
		return false
	}

	last_flower := pais[0].GetFlower()
	for _,p := range pais {
		if p.GetFlower() != last_flower {
			return false
		}
	}

	return true
}

//是否葫芦牛
func isHulu(pais []*ddproto.CommonSrvPokerPai) bool {
	pais_count := map[string]int{}
	for _,p := range pais {
		if _,ok := pais_count[p.GetName()]; ok {
			pais_count[p.GetName()]++
		}else {
			pais_count[p.GetName()] = 1
		}
	}

	if len(pais_count) == 2 {
		return true
	}
	return false
}

//=========================跑三公牌型解析===========================
//是否为爆久
func isBaoJiu(pais []*ddproto.CommonSrvPokerPai) bool {
	pais_count := map[string]int{}
	for _,p := range pais {
		if _,ok := pais_count[p.GetName()]; ok {
			pais_count[p.GetName()]++
		}else {
			pais_count[p.GetName()] = 1
		}
	}

	if len(pais_count) == 1 && pais[0].GetName() == "3" {
		return true
	}

	return false
}

//是否为三条
func isSanTiao(pais []*ddproto.CommonSrvPokerPai) bool {
	pais_count := map[string]int{}
	for _,p := range pais {
		if _,ok := pais_count[p.GetName()]; ok {
			pais_count[p.GetName()]++
		}else {
			pais_count[p.GetName()] = 1
		}
	}

	if len(pais_count) == 1 && pais[0].GetName() != "3" {
		return true
	}

	return false
}

//是否为三公
func isSanGong(pais []*ddproto.CommonSrvPokerPai) bool {
	pais_count := map[string]int{}
	for _,p := range pais {
		if _,ok := pais_count[p.GetName()]; ok {
			pais_count[p.GetName()]++
		}else {
			pais_count[p.GetName()] = 1
		}
		if p.GetName() != "J" && p.GetName() != "Q" && p.GetName() != "K" {
			return false
		}
	}

	if len(pais_count) == 3 {
		return true
	}

	return false
}

//解析三公牌型
func ParsePSGPoker(pais []*ddproto.CommonSrvPokerPai) *ddproto.PsgSrvPoker {
	psg_poker := &ddproto.PsgSrvPoker{
		Pais: pais,
		Type: ddproto.PsgEnum_PokerType_PSG_POKER_0_DIAN.Enum(),
	}

	if len(pais) != 3 {
		log.E("跑三公牌型错误！pais:%v", pais)
		return psg_poker
	}

	for _,p := range pais {
		if p == nil {
			log.E("跑三公牌型错误！存在nil值！pais:%v", pais)
			return psg_poker
		}
	}

	//是否为爆久
	if isBaoJiu(pais) {
		psg_poker.Type = ddproto.PsgEnum_PokerType_PSG_POKER_BAO_JIU.Enum()
		return psg_poker
	}

	//是否为三条
	if isSanTiao(pais) {
		psg_poker.Type = ddproto.PsgEnum_PokerType_PSG_POKER_SAN_TIAO.Enum()
		return psg_poker
	}

	//是否为三公
	if isSanGong(pais) {
		psg_poker.Type = ddproto.PsgEnum_PokerType_PSG_POKER_SAN_GONG.Enum()
		return psg_poker
	}

	//0-9点
	point := (pais[0].GetValue() + pais[1].GetValue() + pais[2].GetValue())%10
	psg_poker.Type = ddproto.PsgEnum_PokerType(point + int32(ddproto.PsgEnum_PokerType_PSG_POKER_0_DIAN)).Enum()
	return psg_poker
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
func GetClientPoker(srv *ddproto.NiuniuSrvPoker) *ddproto.NiuniuClientPoker {
	if srv == nil {
		log.E("poker nil")
		return nil
	}
	server_pais := srv.Pais
	client_pais := []*ddproto.ClientBasePoker{}

	for _,s_pai := range server_pais {
		var p_num int32 = 0
		switch s_pai.GetName() {
		case "J":
			p_num = 11
		case "Q":
			p_num = 12
		case "K":
			p_num = 13
		case "A":
			p_num = 14
		default:
			p_num = s_pai.GetValue()
		}
		c_pai := &ddproto.ClientBasePoker{
			Suit: flower2clientFlower(s_pai.GetDes()).Enum(),
			Num: proto.Int32(p_num),
			Id: proto.Int32(s_pai.GetId()+1),
		}
		client_pais = append(client_pais, c_pai)
	}

	clent_poker := &ddproto.NiuniuClientPoker{
		Pais: client_pais,
		Type: ddproto.NiuniuEnum_PokerType(srv.GetType()).Enum(),
		SelectedId: srv.SelectedId,
	}

	return clent_poker
}

//牌值map
var poker_score_map map[ddproto.NiuniuEnum_PokerType]int64 = map[ddproto.NiuniuEnum_PokerType]int64{
	ddproto.NiuniuEnum_PokerType_NO_NIU: 1,
	ddproto.NiuniuEnum_PokerType_NIU_1: 1,
	ddproto.NiuniuEnum_PokerType_NIU_2: 1,
	ddproto.NiuniuEnum_PokerType_NIU_3: 1,
	ddproto.NiuniuEnum_PokerType_NIU_4: 1,
	ddproto.NiuniuEnum_PokerType_NIU_5: 1,
	ddproto.NiuniuEnum_PokerType_NIU_6: 1,
	ddproto.NiuniuEnum_PokerType_NIU_7: 2,
	ddproto.NiuniuEnum_PokerType_NIU_8: 2,
	ddproto.NiuniuEnum_PokerType_NIU_9: 2,
	ddproto.NiuniuEnum_PokerType_NIU_NIU: 3,
	ddproto.NiuniuEnum_PokerType_YIN_NIU: 5,  //银牛、五花牛
	ddproto.NiuniuEnum_PokerType_NIU_ZHA_DAN: 6,  //炸弹牛
	ddproto.NiuniuEnum_PokerType_NIU_YI_TIAO_LONG: 5,  //一条龙、顺子牛
	ddproto.NiuniuEnum_PokerType_NIU_TONGHUA_NIU: 5,  //同花牛
	ddproto.NiuniuEnum_PokerType_NIU_HULU_NIU: 6,  //葫芦牛
	ddproto.NiuniuEnum_PokerType_NIU_XIAOSHUAI_NIU: 8,  //同花顺牛、小帅牛、五小牛
}

//疯狂加倍-牌值
var super_poker_score_map map[ddproto.NiuniuEnum_PokerType]int64 = map[ddproto.NiuniuEnum_PokerType]int64{
	ddproto.NiuniuEnum_PokerType_NO_NIU: 1,
	ddproto.NiuniuEnum_PokerType_NIU_1: 1,
	ddproto.NiuniuEnum_PokerType_NIU_2: 2,
	ddproto.NiuniuEnum_PokerType_NIU_3: 3,
	ddproto.NiuniuEnum_PokerType_NIU_4: 4,
	ddproto.NiuniuEnum_PokerType_NIU_5: 5,
	ddproto.NiuniuEnum_PokerType_NIU_6: 6,
	ddproto.NiuniuEnum_PokerType_NIU_7: 7,
	ddproto.NiuniuEnum_PokerType_NIU_8: 8,
	ddproto.NiuniuEnum_PokerType_NIU_9: 9,
	ddproto.NiuniuEnum_PokerType_NIU_NIU: 10,
	ddproto.NiuniuEnum_PokerType_YIN_NIU: 11,
	ddproto.NiuniuEnum_PokerType_NIU_ZHA_DAN: 15,
	ddproto.NiuniuEnum_PokerType_NIU_YI_TIAO_LONG: 13,
	ddproto.NiuniuEnum_PokerType_NIU_TONGHUA_NIU: 14,
	ddproto.NiuniuEnum_PokerType_NIU_HULU_NIU: 12,
	ddproto.NiuniuEnum_PokerType_NIU_XIAOSHUAI_NIU: 16,
}


//牌型换牌概率100
var Exchange_poker_score_map map[ddproto.NiuniuEnum_PokerType]int = map[ddproto.NiuniuEnum_PokerType]int{
	ddproto.NiuniuEnum_PokerType_NO_NIU: 0,
	ddproto.NiuniuEnum_PokerType_NIU_1: 0,
	ddproto.NiuniuEnum_PokerType_NIU_2: 0,
	ddproto.NiuniuEnum_PokerType_NIU_3: 0,
	ddproto.NiuniuEnum_PokerType_NIU_4: 0,
	ddproto.NiuniuEnum_PokerType_NIU_5: 0,
	ddproto.NiuniuEnum_PokerType_NIU_6: 0,
	ddproto.NiuniuEnum_PokerType_NIU_7: 30,
	ddproto.NiuniuEnum_PokerType_NIU_8: 40,
	ddproto.NiuniuEnum_PokerType_NIU_9: 50,
	ddproto.NiuniuEnum_PokerType_NIU_NIU: 60,
	ddproto.NiuniuEnum_PokerType_YIN_NIU: 0,
	ddproto.NiuniuEnum_PokerType_NIU_ZHA_DAN: 0,
	ddproto.NiuniuEnum_PokerType_NIU_YI_TIAO_LONG: 0,
	ddproto.NiuniuEnum_PokerType_NIU_TONGHUA_NIU: 0,
	ddproto.NiuniuEnum_PokerType_NIU_HULU_NIU: 0,
	ddproto.NiuniuEnum_PokerType_NIU_XIAOSHUAI_NIU: 0,
}

//获取牌型对应的分值
func GetPokerScore(poker *ddproto.NiuniuSrvPoker, option *ddproto.NiuniuDeskOption) int64 {
	//疯狂加倍
	if option.GetIsSuperJiabei() {
		if score,ok := super_poker_score_map[poker.GetType()];ok {
			return score
		}
	}
	switch option.GetFanbeiRule() {
	case ddproto.NiuniuEnumFanbeiRule_FANBEI_RULE_1:
		switch poker.GetType() {
		case ddproto.NiuniuEnum_PokerType_NIU_NIU:
			return 4
		case ddproto.NiuniuEnum_PokerType_NIU_9:
			return 3
		case ddproto.NiuniuEnum_PokerType_NIU_8:
			return 2
		case ddproto.NiuniuEnum_PokerType_NIU_7:
			return 2
		default:
			if score,ok := poker_score_map[poker.GetType()];ok {
				return score
			}
		}
	case ddproto.NiuniuEnumFanbeiRule_FANBEI_RULE_2:
		switch poker.GetType() {
		case ddproto.NiuniuEnum_PokerType_NIU_NIU:
			return 3
		case ddproto.NiuniuEnum_PokerType_NIU_9:
			return 2
		case ddproto.NiuniuEnum_PokerType_NIU_8:
			return 2
		default:
			if score,ok := poker_score_map[poker.GetType()];ok {
				return score
			}
		}
	default:
		if score,ok := poker_score_map[poker.GetType()];ok {
			return score
		}
	}
	return 1
}

//比较牌值大小
func GetPokerPaiValue(pai *ddproto.CommonSrvPokerPai) int32 {
	switch pai.GetName() {
	case "A":
		return 1
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	default:
		return pai.GetValue()
	}
}

//比牌
func IsBigThanBanker(banker_poker *ddproto.NiuniuSrvPoker, user_poker *ddproto.NiuniuSrvPoker) bool {
	if banker_poker == nil || user_poker == nil || len(banker_poker.GetPais()) != 5 || len(user_poker.GetPais()) != 5 {
		log.E("牌型异常 banker_poker:%v user_poker:%v", banker_poker, user_poker)
		if len(banker_poker.GetPais()) > len(banker_poker.GetPais()) {
			return false
		}else {
			return true
		}
	}
	if banker_poker.GetType() == user_poker.GetType() {
		var max_1, max_2, max_1_id, max_2_id int32 = 0, 0, 0, 0
		for i:=0;i<5;i++ {
			if val_1 := GetPokerPaiValue(banker_poker.GetPais()[i]);val_1 > max_1 {
				//先比较牌值
				max_1 = val_1
				max_1_id = banker_poker.GetPais()[i].GetId()
			}else if val_1 == max_1 {
				//再比较花色
				if max_1_id > banker_poker.GetPais()[i].GetId() {
					max_1_id = banker_poker.GetPais()[i].GetId()
				}
			}
			if val_2 := GetPokerPaiValue(user_poker.GetPais()[i]);val_2 > max_2 {
				max_2 = val_2
				max_2_id = user_poker.GetPais()[i].GetId()
			}else if val_2 == max_2 {
				//再比较花色
				if max_2_id > user_poker.GetPais()[i].GetId() {
					max_2_id = user_poker.GetPais()[i].GetId()
				}
			}
		}
		if max_1 == max_2 {
			//比花色，黑桃>红桃>梅花>方块
			if max_2_id < max_1_id {
				return true
			}else {
				return false
			}
		}else if max_2 > max_1{
			return true
		}else {
			return false
		}
	}else if user_poker.GetType() > banker_poker.GetType() {
		return true
	}else {
		return false
	}

	return false
}
