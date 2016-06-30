package porkService

import (
	"casino_server/utils"
	"strings"
	"casino_server/common/log"
	"fmt"
	"sort"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/numUtils"
)

/**
	制作一副扑克牌
	1,一副扑克牌,完整的54张牌
	2,能更具程序随机的产生指定的牌
	3,适应各种游戏模式

 */

/**
	定义一副牌
 */

var PLAYER_COUNT int32 = 5				//总的玩家个数
const  PLAY_PORK_COUNT int32 = 15		//每局需要的总牌数


var ZJH_TYPE_SANPAI 	int32	= 1;
var ZJH_TYPE_DUIZI  	int32	= 2;
var ZJH_TYPE_SHUNZI 	int32	= 3;
var ZJH_TYPE_TONGHUA  	int32	= 4;
var ZJH_TYPE_TONGHUASHUN int32	= 5;
var ZJH_TYPE_BAOZI  	int32	= 6;


var porkMap map[int32]string

func init(){
	porkMap = make(map[int32]string,54)
	porkMap[1]  = "POKER_HEART_2"
	porkMap[2]  = "POKER_DIAMOND_2"
	porkMap[3]  = "POKER_CLUB_2"
	porkMap[4]  = "POKER_SPADE_2"
	porkMap[5]  = "POKER_HEART_3"
	porkMap[6]  = "POKER_DIAMOND_3"
	porkMap[7]  = "POKER_CLUB_3"
	porkMap[8]  = "POKER_SPADE_3"
	porkMap[9]  = "POKER_HEART_4"
	porkMap[10] = "POKER_DIAMOND_4"
	porkMap[11] = "POKER_CLUB_4"
	porkMap[12] = "POKER_SPADE_4"
	porkMap[13] = "POKER_HEART_5"
	porkMap[14] = "POKER_DIAMOND_5"
	porkMap[15] = "POKER_CLUB_5"
	porkMap[16] = "POKER_SPADE_5"
	porkMap[17] = "POKER_HEART_6"
	porkMap[18] = "POKER_DIAMOND_6"
	porkMap[19] = "POKER_CLUB_6"
	porkMap[20] = "POKER_SPADE_6"
	porkMap[21] = "POKER_HEART_7"
	porkMap[22] = "POKER_DIAMOND_7"
	porkMap[23] = "POKER_CLUB_7"
	porkMap[24] = "POKER_SPADE_7"
	porkMap[25] = "POKER_HEART_8"
	porkMap[26] = "POKER_DIAMOND_8"
	porkMap[27] = "POKER_CLUB_8"
	porkMap[28] = "POKER_SPADE_8"
	porkMap[29] = "POKER_HEART_9"
	porkMap[30] = "POKER_DIAMOND_9"
	porkMap[31] = "POKER_CLUB_9"
	porkMap[32] = "POKER_SPADE_9"
	porkMap[33] = "POKER_HEART_10"
	porkMap[34] = "POKER_DIAMOND_10"
	porkMap[35] = "POKER_CLUB_10"
	porkMap[36] = "POKER_SPADE_10"
	porkMap[37] = "POKER_HEART_11_J"
	porkMap[38] = "POKER_DIAMOND_11_J"
	porkMap[39] = "POKER_CLUB_11_J	"
	porkMap[40] = "POKER_SPADE_11_J"
	porkMap[41] = "POKER_HEART_12_Q"
	porkMap[42] = "POKER_DIAMOND_12_Q"
	porkMap[43] = "POKER_CLUB_12_Q"
	porkMap[44] = "POKER_SPADE_12_Q"
	porkMap[45] = "POKER_HEART_13_K"
	porkMap[46] = "POKER_DIAMOND_13_K"
	porkMap[47] = "POKER_CLUB_13_K"
	porkMap[48] = "POKER_SPADE_13_K"
	porkMap[49] = "POKER_HEART_14_A"
	porkMap[50] = "POKER_DIAMOND_14_A"
	porkMap[51] = "POKER_CLUB_14_A"
	porkMap[52] = "POKER_SPADE_14_A"
	porkMap[53] = "POKER_RED_JOKER"
	porkMap[54] = "POKER_BLACK_JOKER"
}

/**
	5个人玩的一组牌,需要更具大小进行排序
 */
type ZjhPorkList []*ZjhPork		//表示扎金花的一组牌

/**
使用sort包需要实现的方法
 */
func ( list ZjhPorkList) Len() int{
	return len(list)
}

//由大到小的排序
func ( list ZjhPorkList) Less(i,j int) bool{

	if list[i].zjhType > list[j].zjhType{	//比较类型
		return true
	}else if list[i].zjhType < list[j].zjhType{
		return false
	}else {
		//当类型相同时候的比较各种不同的牌型
		switch list[i].zjhType {
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
func ( list ZjhPorkList) Swap(i,j int){
	var temp *ZjhPork = list[i]
	list[i] = list[j]
	list[j] = temp
}



func (z ZjhPorkList) String() string{
	result := ""
	for i := 0; i < len(z); i++ {
		result = strings.Join([]string{result,z[i].String()}, "-------\n")
	}
	return result
}

/**
	表示一张牌
 */
type Pork struct {
	mapKey	int32
	mapDes  string		//从map出来时候的描述
	value int32		//值,A对应14,K对应13
	flower string		//花色 四中
	name string		//牌值,2345678910JQKA
}

/**
通过描述来初始化一张牌
 */
func (p *Pork) initPork() error{
	sarry :=  strings.Split(p.mapDes,"_")
	p.value = int32(numUtils.String2Int(sarry[2]))
	p.name = sarry[len(sarry)-1]
	p.flower = sarry[1]
	return nil
}


/**
	扎金花牌的结构
 */
type ZjhPork struct {
	zjhType		int32 	//哪种牌,散牌,对子,顺子,同花,同花顺
	pork1	*Pork		//第一张牌
	pork2	*Pork		//第二张牌
	pork3 	*Pork		//第三张牌
}

func (z *ZjhPork) String() string{
	return strings.Join([]string{z.pork1.mapDes, z.pork2.mapDes,z.pork3.mapDes}, "-")
}


/**
比较两个散牌的大小
 */
func compareSanPai(a,b *ZjhPork) bool{
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
func compreZuiDa(a,b *ZjhPork) bool{
	if a.getSanpai1() > b.getSanpai1() {
		return true
	}else {
		return false
	}
}

func compareDuizi(a,b *ZjhPork) bool{
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

/**
	取散牌最大的
 */
func (z *ZjhPork) getSanpai1() int32{
	if z.pork1.value > z.pork2.value {
		if z.pork1.value > z.pork3.value {
			return z.pork1.value
		}else{
			return z.pork3.value
		}
	}else {
		if z.pork2.value > z.pork3.value {
			return z.pork2.value
		}else{
			return z.pork3.value
		}
	}

}
/**
	取散牌第二大的
 */
func (z *ZjhPork) getSanpai2() int32{
	if z.pork1.value > z.pork2.value {
		if z.pork1.value > z.pork3.value {
			if z.pork2.value > z.pork3.value {
				return z.pork2.value
			}else{
				return z.pork3.value
			}
		}else{
			return z.pork1.value
		}
	}else {
		if z.pork1.value < z.pork3.value {
			if z.pork2.value < z.pork3.value {
				return z.pork2.value
			}else{
				return z.pork3.value
			}
		}else {
			return z.pork1.value
		}
	}
}

/**
	取散牌第三大的
 */
func (z *ZjhPork) getSanpai3() int32{
	if z.pork1.value > z.pork2.value {
		if z.pork2.value > z.pork3.value {
			return z.pork3.value
		}else{
			return z.pork2.value
		}
	}else {
		if z.pork3.value < z.pork1.value {
			return z.pork3.value
		}else{
			return z.pork1.value
		}
	}

}

/**
	取对子的值
 */
func (z *ZjhPork) getDuizi() int32{
	//是否需要是对子才能取对子
	if z.pork1.value == z.pork2.value {
		return z.pork1.value
	}else if z.pork1.value == z.pork3.value {
		return z.pork1.value
	}else {
		return z.pork3.value
	}
}


/**
	取对子中的散牌
 */
func (z *ZjhPork) getDuiziSanPai() int32 {
	//是否需要是对子才能取对子
	if z.pork1.value == z.pork2.value {
		return z.pork3.value
	}else if z.pork1.value == z.pork3.value {
		return z.pork2.value
	}else {
		return z.pork3.value
	}
}

/**
	得到扎金花牌的类型
 */
func (z *ZjhPork) getZjhType(){

}


func CreateZjhList() ZjhPorkList{
	result := ZjhPorkList{}
	indexs := RandomPorkIndex(1,53)	//5组扎金花牌 总共15张牌
	log.T("找到的索引%v",indexs)
	fmt.Println("找到的索引%v",indexs)
	for i := int32(0);i < PLAYER_COUNT;i++ {
		z := &ZjhPork{}
		z.pork1 = CreatePorkByIndex(indexs[i*3])
		z.pork2 = CreatePorkByIndex(indexs[i*3+1])
		z.pork3 = CreatePorkByIndex(indexs[i*3+2])
		result = append(result,z)
	}

	sort.Sort(result)		//对扎金牌数组进行从大到小的排序
	return result
}

/**
更具index生成一张纸牌
 */
func CreatePorkByIndex(i int32) *Pork{
	result :=&Pork{}
	result.mapKey = i
	result.mapDes = porkMap[i]
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



//把这里的拍转换成扎金花proto的牌类型

func ConstructZjhPai(p *ZjhPork) *bbproto.ZjhPai{
	result := &bbproto.ZjhPai{}
	result.Pai1 = constructPai(p.pork1)
	result.Pai2 = constructPai(p.pork2)
	result.Pai3 = constructPai(p.pork3)
	result.PaiType = &p.zjhType
	return result
}

func constructPai(p *Pork) *bbproto.Pai{
	result := &bbproto.Pai{}
	result.Flower = &p.flower
	result.Mapdes = &p.mapDes
	result.Value = &p.value
	result.Name = &p.name
	return result
}