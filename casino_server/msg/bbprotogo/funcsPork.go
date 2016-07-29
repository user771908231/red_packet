package bbproto

import (
	"strings"
	"casino_server/utils/numUtils"
	"sort"
	"casino_server/utils"
	"casino_server/common/log"
)

var PLAYER_COUNT int32 = 5				//总的玩家个数
const  PLAY_PORK_COUNT int32 = 15			//每局需要的总牌数


var porkMap map[int32]string

func init(){
	porkMap = make(map[int32]string,54)
	//porkMap[1]  = "POKER_heart_2_2"
	//porkMap[2]  = "POKER_heart_3_3"
	//porkMap[3]  = "POKER_heart_4_4"
	//porkMap[4]  = "POKER_heart_5_5"
	//porkMap[5]  = "POKER_heart_6_6"
	//porkMap[6]  = "POKER_heart_7_7"
	//porkMap[7]  = "POKER_heart_8_8"
	//porkMap[8]  = "POKER_heart_9_9"
	//porkMap[9]  = "POKER_heart_10_10"
	//porkMap[10] = "POKER_heart_11_J"
	//porkMap[11] = "POKER_heart_12_Q"
	//porkMap[12] = "POKER_heart_13_K"
	//porkMap[13] = "POKER_heart_14_A"
	//
	//porkMap[14] = "POKER_diamond_2_2"
	//porkMap[15] = "POKER_diamond_3_3"
	//porkMap[16] = "POKER_diamond_4_4"
	//porkMap[17] = "POKER_diamond_5_5"
	//porkMap[18] = "POKER_diamond_6_6"
	//porkMap[19] = "POKER_diamond_7_7"
	//porkMap[20] = "POKER_diamond_8_8"
	//porkMap[21] = "POKER_diamond_9_9"
	//porkMap[22] = "POKER_diamond_10_10"
	//porkMap[23] = "POKER_diamond_11_J"
	//porkMap[24] = "POKER_diamond_12_Q"
	//porkMap[25] = "POKER_diamond_13_K"
	//porkMap[26] = "POKER_diamond_14_A"
	//
	//porkMap[27] = "POKER_club_2_2"
	//porkMap[28] = "POKER_club_3_3"
	//porkMap[29] = "POKER_club_4_4"
	//porkMap[30] = "POKER_club_5_5"
	//porkMap[31] = "POKER_club_6_6"
	//porkMap[32] = "POKER_club_7_7"
	//porkMap[33] = "POKER_club_8_8"
	//porkMap[34] = "POKER_club_9_9"
	//porkMap[35] = "POKER_club_10_10"
	//porkMap[36] = "POKER_club_11_J	"
	//porkMap[37] = "POKER_club_12_Q"
	//porkMap[38] = "POKER_club_13_K"
	//porkMap[39] = "POKER_club_14_A"
	//
	//porkMap[40] = "POKER_spade_2_2"
	//porkMap[41] = "POKER_spade_3_3"
	//porkMap[42] = "POKER_spade_4_4"
	//porkMap[43] = "POKER_spade_5_5"
	//porkMap[44] = "POKER_spade_6_6"
	//porkMap[45] = "POKER_spade_7_7"
	//porkMap[46] = "POKER_spade_8_8"
	//porkMap[47] = "POKER_spade_9_9"
	//porkMap[48] = "POKER_spade_10_10"
	//porkMap[49] = "POKER_spade_11_J"
	//porkMap[50] = "POKER_spade_12_Q"
	//porkMap[51] = "POKER_spade_13_K"
	//porkMap[52] = "POKER_spade_14_A"
	//
	//porkMap[53] = "POKER_RED_JOKER"
	//porkMap[54] = "POKER_BLACK_JOKER"

	porkMap[1]  = "POKER_diamond_2_2"
	porkMap[2]  = "POKER_diamond_3_3"
	porkMap[3]  = "POKER_diamond_4_4"
	porkMap[4]  = "POKER_diamond_5_5"
	porkMap[5]  = "POKER_diamond_6_6"
	porkMap[6]  = "POKER_diamond_7_7"
	porkMap[7]  = "POKER_diamond_8_8"
	porkMap[8]  = "POKER_diamond_9_9"
	porkMap[9]  = "POKER_diamond_10_10"
	porkMap[10] = "POKER_diamond_11_J"
	porkMap[11] = "POKER_diamond_12_Q"
	porkMap[12] = "POKER_diamond_13_K"
	porkMap[13] = "POKER_diamond_14_A"

	porkMap[14] = "POKER_club_2_2"
	porkMap[15] = "POKER_club_3_3"
	porkMap[16] = "POKER_club_4_4"
	porkMap[17] = "POKER_club_5_5"
	porkMap[18] = "POKER_club_6_6"
	porkMap[19] = "POKER_club_7_7"
	porkMap[20] = "POKER_club_8_8"
	porkMap[21] = "POKER_club_9_9"
	porkMap[22] = "POKER_club_10_10"
	porkMap[23] = "POKER_club_11_J	"
	porkMap[24] = "POKER_club_12_Q"
	porkMap[25] = "POKER_club_13_K"
	porkMap[26] = "POKER_club_14_A"


	porkMap[27] = "POKER_heart_2_2"
	porkMap[28] = "POKER_heart_3_3"
	porkMap[29] = "POKER_heart_4_4"
	porkMap[30] = "POKER_heart_5_5"
	porkMap[31] = "POKER_heart_6_6"
	porkMap[32] = "POKER_heart_7_7"
	porkMap[33] = "POKER_heart_8_8"
	porkMap[34] = "POKER_heart_9_9"
	porkMap[35] = "POKER_heart_10_10"
	porkMap[36] = "POKER_heart_11_J"
	porkMap[37] = "POKER_heart_12_Q"
	porkMap[38] = "POKER_heart_13_K"
	porkMap[39] = "POKER_heart_14_A"


	porkMap[40] = "POKER_spade_2_2"
	porkMap[41] = "POKER_spade_3_3"
	porkMap[42] = "POKER_spade_4_4"
	porkMap[43] = "POKER_spade_5_5"
	porkMap[44] = "POKER_spade_6_6"
	porkMap[45] = "POKER_spade_7_7"
	porkMap[46] = "POKER_spade_8_8"
	porkMap[47] = "POKER_spade_9_9"
	porkMap[48] = "POKER_spade_10_10"
	porkMap[49] = "POKER_spade_11_J"
	porkMap[50] = "POKER_spade_12_Q"
	porkMap[51] = "POKER_spade_13_K"
	porkMap[52] = "POKER_spade_14_A"

	porkMap[53] = "POKER_RED_JOKER"
	porkMap[54] = "POKER_BLACK_JOKER"
}

var ZJH_TYPE_SANPAI 	int32	= 1;
var ZJH_TYPE_DUIZI  	int32	= 2;
var ZJH_TYPE_SHUNZI 	int32	= 3;
var ZJH_TYPE_TONGHUA  	int32	= 4;
var ZJH_TYPE_TONGHUASHUN int32	= 5;
var ZJH_TYPE_BAOZI  	int32	= 6;


type ZjhPaiList []*ZjhPai

/**
使用sort包需要实现的方法
 */
func ( list ZjhPaiList) Len() int{
	return len(list)
}

//由大到小的排序
func ( list ZjhPaiList) Less(i,j int) bool{

	if list[i].GetPaiType() > list[j].GetPaiType(){	//比较类型
		return true
	}else if list[i].GetPaiType() < list[j].GetPaiType(){
		return false
	}else {
		//当类型相同时候的比较各种不同的牌型
		switch list[i].GetPaiType() {
		case ZJH_TYPE_SANPAI:	//散牌
			return compareSanPai(list[i],list[j])
		case ZJH_TYPE_DUIZI:
			return compareDuizi(list[i],list[j])
		case ZJH_TYPE_SHUNZI:
			return compreZuiDa(list[i],list[j])
		case ZJH_TYPE_TONGHUA:
			return compareSanPai(list[i],list[j])
		case ZJH_TYPE_TONGHUASHUN:
			return compreZuiDa(list[i],list[j])
		case ZJH_TYPE_BAOZI:
			return compreZuiDa(list[i],list[j])
		}

		return true
	}
}

//交换函数
func ( list ZjhPaiList) Swap(i,j int){
	var temp *ZjhPai = list[i]
	list[i] = list[j]
	list[j] = temp
}


/**
通过描述来初始化一张牌
 */
func (p *Pai) initPork() error{
	sarry :=  strings.Split(p.GetMapdes(),"_")

	var pvalue int32 = int32(numUtils.String2Int(sarry[2]))
	var pname string = sarry[3]
	var pflower string = sarry[1]

	p.Value  = &pvalue
	p.Name   = &pname
	p.Flower = &pflower
	return nil
}

func (z *ZjhPai) ToString() string{
	return strings.Join([]string{z.Pai[0].GetMapdes(), z.Pai[1].GetMapdes(),z.Pai[2].GetMapdes()}, "-")
}


/**
比较两个散牌的大小
 */
func compareSanPai(a,b *ZjhPai) bool{
	if a.getSanpai1() > b.getSanpai1() {
		return true
	}else if a.getSanpai1() == b.getSanpai1() {
		if a.getSanpai2() > b.getSanpai2() {
			return true
		}else if a.getSanpai2() == b.getSanpai2() {
			if a.getSanpai3() > b.getSanpai3() {
				return true
			}else{
				return false
			}
		}else{
			return false
		}
	}else{
		return  false
	}
}

/**
只比较一组牌中最大的
 */
func compreZuiDa(a,b *ZjhPai) bool{
	if a.getSanpai1() > b.getSanpai1() {
		return true
	}else {
		return false
	}
}

func compareDuizi(a,b *ZjhPai) bool{
	if a.getDuizi() > b.getDuizi() {
		return true
	}else if a.getDuizi() == b.getDuizi() {
		if a.getDuiziSanPai() > b.getDuiziSanPai() {
			return true
		}else{
			return false
		}
	}else{
		return false
	}

}

func (z *ZjhPai) getIntArray() []int{
	data := make([]int,3)
	data[0] = int(z.Pai[0].GetValue())
	data[1] = int(z.Pai[1].GetValue())
	data[2] = int(z.Pai[2].GetValue())
	return data
}

func (z *ZjhPai) getSortIntArray() []int{
	data := z.getIntArray()
	sort.Ints(data)
	return data
}



/**
	取散牌最大的
 */
func (z *ZjhPai) getSanpai1() int{
	data := z.getSortIntArray()
	return data[2]

}
/**
	取散牌第二大的
 */
func (z *ZjhPai) getSanpai2() int{
	data := z.getSortIntArray()
	sort.Ints(data)
	return data[1]
}

/**
	取散牌第三大的
 */
func (z *ZjhPai) getSanpai3() int{
	data := z.getSortIntArray()
	sort.Ints(data)
	return data[0]

}

/**
	取对子的值
 */
func (z *ZjhPai) getDuizi() int{
	data := z.getSortIntArray()
	sort.Ints(data)
	return data[1]
}

/**
	取对子中的散牌
 */
func (z *ZjhPai) getDuiziSanPai() int32 {
	//是否需要是对子才能取对子
	if z.Pai[0].GetValue() == z.Pai[1].GetValue() {
		return z.Pai[2].GetValue()
	}else if z.Pai[0].GetValue() == z.Pai[2].GetValue() {
		return z.Pai[1].GetValue()
	}else {
		return z.Pai[2].GetValue()
	}
}

/**
	得到扎金花牌的类型
 */
func (z *ZjhPai) initZjhType(){
	if z.GetPaiType() == 0{
		//开始初始化扎金花的类型,从最大的开始算
		if z.Pai[0].GetValue() == z.Pai[1].GetValue() && z.Pai[0].GetValue() == z.Pai[2].GetValue() {
			z.PaiType = &ZJH_TYPE_BAOZI
		}else if strings.EqualFold(z.Pai[0].GetFlower(), z.Pai[1].GetFlower()) && strings.EqualFold(z.Pai[0].GetFlower(),z.Pai[2].GetFlower()){
			//判断是否是同花顺
			if  checkShunzi(z.Pai[0].GetValue(),z.Pai[1].GetValue(),z.Pai[2].GetValue()){
				z.PaiType = &ZJH_TYPE_TONGHUASHUN
			}else{
				z.PaiType = &ZJH_TYPE_TONGHUA
			}
		}else if checkShunzi(z.Pai[0].GetValue(),z.Pai[1].GetValue(),z.Pai[2].GetValue()) {
			//判断是否是顺子
			z.PaiType = &ZJH_TYPE_SHUNZI
		}else if z.Pai[0].GetValue() != z.Pai[1].GetValue()  && z.Pai[0].GetValue()  != z.Pai[2].GetValue()  && z.Pai[1].GetValue()  != z.Pai[2].GetValue() {
			z.PaiType = &ZJH_TYPE_SANPAI
		}else{
			z.PaiType = &ZJH_TYPE_DUIZI
		}
	}
}

/**
判断三个数字是否是顺子
 */
func checkShunzi(a,b,c int32) bool{
	data := []int{int(a),int(b),int(c)}
	sort.Ints(data)
	if data[2] - data[1] == 1 && data[1] - data[0] == 1 {
		return true
	}else{
		return false
	}
}


func CreateZjhList() ZjhPaiList{
	result := ZjhPaiList{}
	indexs := RandomPorkIndex(1,53)	//5组扎金花牌 总共15张牌
	log.T("找到的索引%v",indexs)
	//fmt.Println("找到的索引%v",indexs)
	for i := int32(0);i < PLAYER_COUNT;i++ {
		z := &ZjhPai{}
		z.Pai = make([]*Pai,3)
		z.Pai[0] = CreatePorkByIndex(indexs[i*3])
		z.Pai[1] = CreatePorkByIndex(indexs[i*3+1])
		z.Pai[2] = CreatePorkByIndex(indexs[i*3+2])
		z.initZjhType()		//初始化扎金花牌的类型
		//log.T("%v副扎金花的牌:%v",i,z)
		result = append(result,z)
	}

	//todo 可以不需要排序
	sort.Sort(result)		//对扎金牌数组进行从大到小的排序
	log.T("排序之后的牌:",result)
	return result
}

/**
更具index生成一张纸牌
 */
func CreatePorkByIndex(i int32) *Pai{
	//var rmapkey int32 = i
	result :=&Pai{}
	result.MapKey = &i
	var rmapdes string = porkMap[i]
	result.Mapdes = &rmapdes
	result.initPork()
	return result
}


/**
	产生随机数:min==1:max==53
 */
func RandomPorkIndex(min, max int32) [PLAY_PORK_COUNT]int32 {
	result := new([PLAY_PORK_COUNT]int32);
	count := int32(0);
	for count < PLAY_PORK_COUNT {
		num := utils.Rand(min, max)
		flag := true;
		for j := int32(0); j < PLAY_PORK_COUNT; j++ {
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
	return *result;
}

func CreateZjhPaiList(){

}