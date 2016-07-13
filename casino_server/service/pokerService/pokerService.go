package pokerService

import (
	"casino_server/utils"
	"strings"
	"casino_server/common/log"
	"sort"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/numUtils"
)

/**
	制作一副扑克牌
	1,一副扑克牌,完整的54张牌
	2,能更具程序随机的产生指定的牌
	3,适应各种游戏模式


	HEART	红桃
	DIAMOND	方块
	CLUB	梅花
	SPADE	黑桃

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
	porkMap[1]  = "POKER_HEART_2_2"
	porkMap[2]  = "POKER_HEART_3_3"
	porkMap[3]  = "POKER_HEART_4_4"
	porkMap[4] = "POKER_HEART_5_5"
	porkMap[5] = "POKER_HEART_6_6"
	porkMap[6] = "POKER_HEART_7_7"
	porkMap[7] = "POKER_HEART_8_8"
	porkMap[8] = "POKER_HEART_9_9"
	porkMap[9] = "POKER_HEART_10_10"
	porkMap[10] = "POKER_HEART_11_J"
	porkMap[11] = "POKER_HEART_12_Q"
	porkMap[12] = "POKER_HEART_13_K"
	porkMap[13] = "POKER_HEART_14_A"

	porkMap[14]  = "POKER_DIAMOND_2_2"
	porkMap[15]  = "POKER_DIAMOND_3_3"
	porkMap[16] = "POKER_DIAMOND_4_4"
	porkMap[17] = "POKER_DIAMOND_5_5"
	porkMap[18] = "POKER_DIAMOND_6_6"
	porkMap[19] = "POKER_DIAMOND_7_7"
	porkMap[20] = "POKER_DIAMOND_8_8"
	porkMap[21] = "POKER_DIAMOND_9_9"
	porkMap[22] = "POKER_DIAMOND_10_10"
	porkMap[23] = "POKER_DIAMOND_11_J"
	porkMap[24] = "POKER_DIAMOND_12_Q"
	porkMap[25] = "POKER_DIAMOND_13_K"
	porkMap[26] = "POKER_DIAMOND_14_A"

	porkMap[27]  = "POKER_CLUB_2_2"
	porkMap[28]  = "POKER_CLUB_3_3"
	porkMap[29] = "POKER_CLUB_4_4"
	porkMap[30] = "POKER_CLUB_5_5"
	porkMap[31] = "POKER_CLUB_6_6"
	porkMap[32] = "POKER_CLUB_7_7"
	porkMap[33] = "POKER_CLUB_8_8"
	porkMap[34] = "POKER_CLUB_9_9"
	porkMap[35] = "POKER_CLUB_10_10"
	porkMap[36] = "POKER_CLUB_11_J	"
	porkMap[37] = "POKER_CLUB_12_Q"
	porkMap[38] = "POKER_CLUB_13_K"
	porkMap[39] = "POKER_CLUB_14_A"

	porkMap[40]  = "POKER_SPADE_2_2"
	porkMap[41]  = "POKER_SPADE_3_3"
	porkMap[42] = "POKER_SPADE_4_4"
	porkMap[43] = "POKER_SPADE_5_5"
	porkMap[44] = "POKER_SPADE_6_6"
	porkMap[45] = "POKER_SPADE_7_7"
	porkMap[46] = "POKER_SPADE_8_8"
	porkMap[47] = "POKER_SPADE_9_9"
	porkMap[48] = "POKER_SPADE_10_10"
	porkMap[49] = "POKER_SPADE_11_J"
	porkMap[50] = "POKER_SPADE_12_Q"
	porkMap[51] = "POKER_SPADE_13_K"
	porkMap[52] = "POKER_SPADE_14_A"

	porkMap[53] = "POKER_RED_JOKER"
	porkMap[54] = "POKER_BLACK_JOKER"
}

/**
	5个人玩的一组牌,需要更具大小进行排序
 */
type ZjhPorkList []*ZjhPoker                //表示扎金花的一组牌

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
	var temp *ZjhPoker = list[i]
	list[i] = list[j]
	list[j] = temp
}



func (z ZjhPorkList) String() string{
	result := ""
	for i := 0; i < len(z); i++ {
		t,_ := numUtils.Int2String(z[i].zjhType)
		result = strings.Join([]string{result,"type:",t ,",",z[i].String()}, "-------\n")
	}
	return result
}

/**
	表示一张牌
 */
type Poker struct {
	mapKey	int32
	mapDes  string		//从map出来时候的描述
	value int		//值,A对应14,K对应13
	flower string		//花色 四中
	name string		//牌值,2345678910JQKA
}

/**
通过描述来初始化一张牌
 */
func (p *Poker) initPoker() error{
	sarry :=  strings.Split(p.mapDes,"_")
	p.value = numUtils.String2Int(sarry[2])
	p.name = sarry[3]
	p.flower = sarry[1]
	return nil
}


/**
	扎金花牌的结构
 */
type ZjhPoker struct {
	zjhType		int32 //哪种牌,散牌,对子,顺子,同花,同花顺
	pork	[]*Poker         //第一张牌
}

func (z *ZjhPoker) String() string{
	return strings.Join([]string{z.pork[0].mapDes, z.pork[1].mapDes,z.pork[2].mapDes}, "-")
}

/**
比较两个散牌的大小
 */
func compareSanPai(a,b *ZjhPoker) bool{
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
func compreZuiDa(a,b *ZjhPoker) bool{
	if a.getSanpai1() > b.getSanpai1() {
		return true
	}else {
		return false
	}
}

func compareDuizi(a,b *ZjhPoker) bool{
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
func (z *ZjhPoker) getSanpai1() int{
	data := []int{z.pork[0].value,z.pork[1].value,z.pork[2].value}
	sort.Ints(data)
	return data[2]

}
/**
	取散牌第二大的
 */
func (z *ZjhPoker) getSanpai2() int{
	data := []int{z.pork[0].value,z.pork[1].value,z.pork[2].value}
	sort.Ints(data)
	return data[1]
}

/**
	取散牌第三大的
 */
func (z *ZjhPoker) getSanpai3() int{
	data := []int{z.pork[0].value,z.pork[1].value,z.pork[2].value}
	sort.Ints(data)
	return data[0]

}

/**
	取对子的值
 */
func (z *ZjhPoker) getDuizi() int{
	data := []int{z.pork[0].value,z.pork[1].value,z.pork[2].value}
	sort.Ints(data)
	return data[1]
}

/**
	取对子中的散牌
 */
func (z *ZjhPoker) getDuiziSanPai() int {
	//是否需要是对子才能取对子
	if z.pork[0].value == z.pork[1].value {
		return z.pork[2].value
	}else if z.pork[0].value == z.pork[2].value {
		return z.pork[1].value
	}else {
		return z.pork[2].value
	}
}

/**
	得到扎金花牌的类型
 */
func (z *ZjhPoker) initZjhType(){
	if z.zjhType == 0{
		//开始初始化扎金花的类型,从最大的开始算
		if z.pork[0].value == z.pork[1].value && z.pork[0].value == z.pork[2].value {
			z.zjhType = ZJH_TYPE_BAOZI
		}else if strings.EqualFold(z.pork[0].flower, z.pork[1].flower) && strings.EqualFold(z.pork[0].flower,z.pork[2].flower){
			//判断是否是同花顺
			if  checkShunzi(z.pork[0].value,z.pork[1].value,z.pork[2].value){
				z.zjhType = ZJH_TYPE_TONGHUASHUN

			}else{
				z.zjhType = ZJH_TYPE_TONGHUA
			}
		}else if checkShunzi(z.pork[0].value,z.pork[1].value,z.pork[2].value) {
			//判断是否是顺子
			z.zjhType = ZJH_TYPE_SHUNZI
		}else if z.pork[0].value != z.pork[1].value && z.pork[0].value != z.pork[2].value && z.pork[1].value != z.pork[2].value{
			z.zjhType = ZJH_TYPE_SANPAI
		}else{
			z.zjhType = ZJH_TYPE_DUIZI
		}
	}
}

/**
判断三个数字是否是顺子
 */
func checkShunzi(a,b,c int) bool{
	data := []int{a,b,c}
	sort.Ints(data)
	if data[2] - data[1] == 1 && data[1] - data[0] == 1 {
		return true
	}else{
		return false
	}
}


func CreateZjhList() ZjhPorkList{
	result := ZjhPorkList{}
	indexs := RandomPorkIndex(1,53)	//5组扎金花牌 总共15张牌
	log.T("找到的索引%v",indexs)
	//fmt.Println("找到的索引%v",indexs)
	for i := int32(0);i < PLAYER_COUNT;i++ {
		z := &ZjhPoker{}
		z.pork = make([]*Poker,3)
		z.pork[0] = CreatePorkByIndex(indexs[i*3])
		z.pork[1] = CreatePorkByIndex(indexs[i*3+1])
		z.pork[2] = CreatePorkByIndex(indexs[i*3+2])
		z.initZjhType()		//初始化扎金花牌的类型
		//log.T("%v副扎金花的牌:%v",i,z)
		result = append(result,z)
	}

	sort.Sort(result)		//对扎金牌数组进行从大到小的排序
	log.T("排序之后的牌:",result)
	return result
}

/**
更具index生成一张纸牌
 */
func CreatePorkByIndex(i int32) *Poker {
	result :=&Poker{}
	result.mapKey = i
	result.mapDes = porkMap[i]
	result.initPoker()
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

func ConstructZjhPai(p *ZjhPoker) *bbproto.ZjhPai{
	result := &bbproto.ZjhPai{}
	ps := make([]*bbproto.Pai,3)
	ps[0] = constructPai(p.pork[0])
	ps[1] = constructPai(p.pork[1])
	ps[2] = constructPai(p.pork[2])
	result.PaiType = &p.zjhType
	return result
}

func constructPai(p *Poker) *bbproto.Pai{
	result := &bbproto.Pai{}
	result.Flower = &p.flower
	result.Mapdes = &p.mapDes
	var v int32  = int32(p.value)
	result.Value = &v
	result.Name = &p.name
	return result
}