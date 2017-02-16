package majiangv2

import (
	"casino_common/common/log"
	"casino_common/utils/rand"
	"sort"
	"casino_majiang/service/majiang"
)

//得到一副牌...

//**处理忙将牌相关的service

var MjpaiMap map[int]string
var clienMap map[int]int32

//const
var W int32 = 1           //万
var S int32 = 2           //条
var T int32 = 3           //筒
var MJPAI_COUNT int = 108 //牌的张数

//基本牌型番数
var FAN_PINGHU int32 = 0        //平胡 0番
var FAN_DADUIZI int32 = 1       //大对子 1番
var FAN_QINGYISE int32 = 2      //平胡清一色 2番
var FAN_DAIYAOJIU int32 = 2     //带幺九 2番
var FAN_QIDUI int32 = 2         //七对 2番
var FAN_QINGDUI int32 = 3       //清对 3番
var FAN_JIANGDUI int32 = 3      //将对 3番
var FAN_LONGQIDUI int32 = 4     //龙七对 4番
var FAN_QINGQIDUI int32 = 4     //清七对 4番
var FAN_JIANGQIDUI int32 = 4    //将七对 4番
var FAN_QINGYAOJIU int32 = 4    //清幺九 4番
var FAN_QINGLONGQIDUI int32 = 5 //清龙七对
var FAN_TIAN_DI_HU int32 = 5    //天地胡

//加番
var FAN_ZIMO int32 = 1            //自摸
var FAN_JINGOUDIAO int32 = 1      //金钩钓
var FAN_MENQ_ZHONGZ int32 = 1     //门清中张
var FAN_GANGSHANGHUA int32 = 1    //杠上花
var FAN_GANGSHANGPAO int32 = 1    //杠上炮
var FAN_HD_LAO int32 = 1          //海底捞
var FAN_HD_PAO int32 = 1          //海底炮
var FAN_QIANGGANG int32 = 1       //抢杠
var FAN_HD_GANGSHANGHUA int32 = 1 //海底杠上花
var FAN_HD_GANGSHANGPAO int32 = 2 //海底杠上炮

//初始化麻将牌
func init() {
	MjpaiMap = make(map[int]string, 108)
	MjpaiMap[0] = "T_1"
	MjpaiMap[1] = "T_1"
	MjpaiMap[2] = "T_1"
	MjpaiMap[3] = "T_1"
	MjpaiMap[4] = "T_2"
	MjpaiMap[5] = "T_2"
	MjpaiMap[6] = "T_2"
	MjpaiMap[7] = "T_2"
	MjpaiMap[8] = "T_3"
	MjpaiMap[9] = "T_3"
	MjpaiMap[10] = "T_3"
	MjpaiMap[11] = "T_3"
	MjpaiMap[12] = "T_4"
	MjpaiMap[13] = "T_4"
	MjpaiMap[14] = "T_4"
	MjpaiMap[15] = "T_4"
	MjpaiMap[16] = "T_5"
	MjpaiMap[17] = "T_5"
	MjpaiMap[18] = "T_5"
	MjpaiMap[19] = "T_5"
	MjpaiMap[20] = "T_6"
	MjpaiMap[21] = "T_6"
	MjpaiMap[22] = "T_6"
	MjpaiMap[23] = "T_6"
	MjpaiMap[24] = "T_7"
	MjpaiMap[25] = "T_7"
	MjpaiMap[26] = "T_7"
	MjpaiMap[27] = "T_7"
	MjpaiMap[28] = "T_8"
	MjpaiMap[29] = "T_8"
	MjpaiMap[30] = "T_8"
	MjpaiMap[31] = "T_8"
	MjpaiMap[32] = "T_9"
	MjpaiMap[33] = "T_9"
	MjpaiMap[34] = "T_9"
	MjpaiMap[35] = "T_9"
	MjpaiMap[36] = "S_1"
	MjpaiMap[37] = "S_1"
	MjpaiMap[38] = "S_1"
	MjpaiMap[39] = "S_1"
	MjpaiMap[40] = "S_2"
	MjpaiMap[41] = "S_2"
	MjpaiMap[42] = "S_2"
	MjpaiMap[43] = "S_2"
	MjpaiMap[44] = "S_3"
	MjpaiMap[45] = "S_3"
	MjpaiMap[46] = "S_3"
	MjpaiMap[47] = "S_3"
	MjpaiMap[48] = "S_4"
	MjpaiMap[49] = "S_4"
	MjpaiMap[50] = "S_4"
	MjpaiMap[51] = "S_4"
	MjpaiMap[52] = "S_5"
	MjpaiMap[53] = "S_5"
	MjpaiMap[54] = "S_5"
	MjpaiMap[55] = "S_5"
	MjpaiMap[56] = "S_6"
	MjpaiMap[57] = "S_6"
	MjpaiMap[58] = "S_6"
	MjpaiMap[59] = "S_6"
	MjpaiMap[60] = "S_7"
	MjpaiMap[61] = "S_7"
	MjpaiMap[62] = "S_7"
	MjpaiMap[63] = "S_7"
	MjpaiMap[64] = "S_8"
	MjpaiMap[65] = "S_8"
	MjpaiMap[66] = "S_8"
	MjpaiMap[67] = "S_8"
	MjpaiMap[68] = "S_9"
	MjpaiMap[69] = "S_9"
	MjpaiMap[70] = "S_9"
	MjpaiMap[71] = "S_9"
	MjpaiMap[72] = "W_1"
	MjpaiMap[73] = "W_1"
	MjpaiMap[74] = "W_1"
	MjpaiMap[75] = "W_1"
	MjpaiMap[76] = "W_2"
	MjpaiMap[77] = "W_2"
	MjpaiMap[78] = "W_2"
	MjpaiMap[79] = "W_2"
	MjpaiMap[80] = "W_3"
	MjpaiMap[81] = "W_3"
	MjpaiMap[82] = "W_3"
	MjpaiMap[83] = "W_3"
	MjpaiMap[84] = "W_4"
	MjpaiMap[85] = "W_4"
	MjpaiMap[86] = "W_4"
	MjpaiMap[87] = "W_4"
	MjpaiMap[88] = "W_5"
	MjpaiMap[89] = "W_5"
	MjpaiMap[90] = "W_5"
	MjpaiMap[91] = "W_5"
	MjpaiMap[92] = "W_6"
	MjpaiMap[93] = "W_6"
	MjpaiMap[94] = "W_6"
	MjpaiMap[95] = "W_6"
	MjpaiMap[96] = "W_7"
	MjpaiMap[97] = "W_7"
	MjpaiMap[98] = "W_7"
	MjpaiMap[99] = "W_7"
	MjpaiMap[100] = "W_8"
	MjpaiMap[101] = "W_8"
	MjpaiMap[102] = "W_8"
	MjpaiMap[103] = "W_8"
	MjpaiMap[104] = "W_9"
	MjpaiMap[105] = "W_9"
	MjpaiMap[106] = "W_9"
	MjpaiMap[107] = "W_9"

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
/**
	
 */
func GettPaiStats(pais []*majiang.MJPai) []int {
	//统计每张牌的重复次数
	counts := make([]int, 27) //0~27
	for i := 0; i < len(pais); i++ {
		pai := pais[i]
		value := pai.GetValue() - 1
		flower := pai.GetFlower() //flower=1,2,3
		value += (flower - 1) * 9
		counts[ value ] ++
	}
	return counts
}

func GettPaiValueByCountPos(countPos int) int32 {
	return int32(countPos%9 + 1)
}

//从pais数组里删除一张pos位置的pai 注 pos是索引值 使用覆盖的方式
func removeFromPais(pais []*majiang.MJPai, pos int) []*majiang.MJPai {
	pais[pos] = pais[len(pais)-1]
	return pais[:len(pais)-1]
}

//将一张pai插入到指定pos的pais数组里去
func addPaiIntoPais(pai *majiang.MJPai, pais []*majiang.MJPai, pos int) []*majiang.MJPai {
	tempPais := make([]*majiang.MJPai, len(pais)+1)

	copy(tempPais[:pos], pais[:pos])
	tempPais[pos] = pai
	copy(tempPais[pos+1:], pais[pos:])
	return tempPais
}

//这张pai是否可碰
// add 增加缺的花色是不能碰的
func CanPengPai(pai *majiang.MJPai, handPai *majiang.MJHandPai) bool {

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
func CanGangPai(pai *majiang.MJPai, handPai *majiang.MJHandPai) (canGang bool, gangPais []*majiang.MJPai) {
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
		log.T("判断手中的牌是否可以杠...")
		//检测手牌中是否有杠
		tempPais := make([]*majiang.MJPai, len(handPai.Pais)+1+len(handPai.PengPais))
		copy(tempPais[0:len(handPai.Pais)], handPai.Pais)
		tempPais[len(handPai.Pais)] = handPai.InPai
		copy(tempPais[len(handPai.Pais)+1:], handPai.PengPais)

		//log.T("tempPais[len(handPai.Pais)] is %v", tempPais[len(handPai.Pais)])
		//log.T("tempPais is %v", tempPais)
		//counts := GettPaiStats(handPai.Pais)
		counts := GettPaiStats(tempPais)
		for _, p := range tempPais {
			if p.GetFlower() == handPai.GetQueFlower() {
				continue
			}
			//log.T("判断杠牌 p.getValue(%v)+p.GetFlower[%v]*9 = %v", p.GetValue(), p.GetFlower(), p.GetValue() + (p.GetFlower() - 1) * 9)
			if ( 4 == counts[ p.GetValue()-1+(p.GetFlower()-1)*9 ] ) {
				canGang = true
				gangPais = append(gangPais, p)
			}
		}
	}

	return canGang, gangPais
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//过滤一副牌中的某一个花色
func IgnoreFlower(pais []*majiang.MJPai, flower int32) []*majiang.MJPai {
	newPais := []*majiang.MJPai{}
	for i := 0; i < len(pais); i++ {
		//log.T("IgnoreFlower: pais[%v] is %v, flower is [%v]", i, pais[i].GetDes(), pais[i].GetFlower())
		if pais[i].GetFlower() != flower {
			//不是需要过滤的花色 append
			newPais = append(newPais, pais[i])
		}
	}
	return newPais
}

//洗牌的算法,这里可以得到一副洗好的麻将
func XiPai() []*majiang.MJPai {

	//初始化一个顺序的牌的集合
	pmap := make([]int, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		pmap[i] = i
	}

	//打乱牌的集合
	pResult := make([]int, MJPAI_COUNT)

	for i := 0; i < MJPAI_COUNT; i++ {
		rand := rand.Rand(0, (int32(MJPAI_COUNT - i)))
		//log.T("得到的rand[%v]", rand)
		pResult[i] = pmap[rand]
		pmap = append(pmap[:int(rand)], pmap[int(rand)+1:]...)
	}

	log.T("洗牌之后,得到的随机的index数组[%v]", pResult)
	//TestCheckRanIndex(pResult)        //todo  测试代码，之后需要删除
	//pResult = []int{78,60,93,14,12,45,28,34,51,61,4,40,36,11,48,37,100,35,1,68,32,66,67,85,101,20,7,69,3,15,95,55,16,56,41,54,33,59,79,38,86,99,87,83,8,73,57,102,90,104,77,63,105,82,62,91,5,39,13,44,46,25,65,27,49,19,80,26,74,71,29,88,53,97,10,92,23,9,81,50,98,22,75,31,6,106,47,17,72,43,94,89,2,42,0,64,24,107,76,18,30,96,58,52,21,70,84,103}

	//开始得到牌的信息
	result := make([]*majiang.MJPai, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		result[i] = InitMjPaiByIndex(pResult[i])
	}

	//log.T("洗牌之后,得到的牌的数组[%v]", result)
	return result
}

func XiPaiTestHu() []*majiang.MJPai {
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

	//开始得到牌的信息
	result := make([]*majiang.MJPai, MJPAI_COUNT)
	for i := 0; i < MJPAI_COUNT; i++ {
		result[i] = InitMjPaiByIndex(pResult[i])
	}

	log.T("洗牌之后,得到的牌的数组[%v]", result)
	return result
}

//通过一个index索引来得到一张牌
func InitMjPaiByIndex(index int) *majiang.MJPai {
	result := majiang.NewMjpai()
	*result.Index = int32(index)
	*result.Des = MjpaiMap[index]
	result.InitByDes()
	return result
}

//通过一个Des描述和现持有的牌来得到一张空闲牌
func InitMjPaiByDes(des string, hand *majiang.MJHandPai) *majiang.MJPai {

	result := majiang.NewMjpai()
	handMJPais := []*majiang.MJPai{}

	//加入杠牌
	handMJPais = append(handMJPais, hand.GangPais...)
	//加入手牌
	handMJPais = append(handMJPais, hand.Pais...)
	//加入碰牌
	handMJPais = append(handMJPais, hand.PengPais...)
	//加入摸牌
	handMJPais = append(handMJPais, hand.InPai)

	for mjpaiIndex, mjpaiDes := range MjpaiMap {
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

func ServerPais2string(pais []*majiang.MJPai) string {
	//首先进行排序
	var tempPais majiang.MjPaiList = make([]*majiang.MJPai, len(pais))
	copy(tempPais, pais)
	sort.Sort(tempPais)

	//得到描述字符串
	s := ""
	for _, p := range pais {
		s = s + p.LogDes() + "\t"
	}

	//返回值
	return s
}
