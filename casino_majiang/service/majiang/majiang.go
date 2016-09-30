package majiang

import (
	"strings"
	"github.com/name5566/leaf/log"
	"casino_server/utils/numUtils"
	"casino_server/utils"
)

//得到一副牌...


//**处理忙将牌相关的service

var mjpaiMap map[int]string
var clienMap map[int]int32

//const
var T int32 = 1        //筒
var S int32 = 2        //条
var W int32 = 3        //万
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


	IsQingYiSe    bool    //是否是清一色
	IsQiDui       bool    //是否是7对
	IsJiangDui    bool    //是否是将对
	IsDaiYao      bool    //是否是带幺
	IsDaDui       bool    //是否是大对子

	PengPaied     []int32 //碰牌
	GangPaied     []int32 //杠牌
	JiangPai      int32   //将牌
}


//统计麻将牌
func (p *PairMjPai)  InitPaiStatistics() error {
	//统计牌的张数
	for i := 0; i < len(p.pais); i++ {
		pai := p.pais[i]
		if pai.GetFlower() == T {
			p.paiStatistics[pai.GetValue() - 1]++
		} else if pai.GetFlower() == S {
			p.paiStatistics[pai.GetValue() - 1 + 9]++
		} else if pai.GetFlower() == W {
			p.paiStatistics[pai.GetValue() - 1 + 9 + 9 ]++
		}
	}
	return nil
}

//判断打这张牌能不能胡牌
func (p *PairMjPai) IsHuPai(pai *MJPai) bool {
	//在所有的牌中增加 pai,判断此牌是否能和

	p.pais = append(p.pais, pai)
	p.InitPaiStatistics()

	if tryHU(p.paiStatistics, len(p.pais)) {
		log.Debug("胡牌了")
	} else {
		log.Debug("没有胡牌")
	}

	//最后需要删除最后一张牌
	p.pais = p.pais[0:len(p.pais) - 1]
	return false;
}


//胡牌的算法
func tryHU(count []int, len int) bool {

	//递归完所有的牌表示 胡了
	if (len == 0) {
		return true;
	}

	if (len % 3 == 2) {
		// 说明对牌没出现
		for i := 0; i < 27; i++ {
			if (count[i] >= 2) {
				count[i] -= 2;
				if (tryHU(count, len - 2)) {
					return true;
				}
				count[i] += 2;
			}
		}
	} else {
		// 三个一样的
		for i := 0; i < 27; i++ {
			if (count[i] >= 3) {
				count[i] -= 3;
				if (tryHU(count, len - 3)) {
					return true;
				}
				count[i] += 3;
			}
		}
		// 是否是顺子，这里应该分开判断
		for i := 0; i < 7; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				if (tryHU(count, len - 3)) {
					return true;
				}
				count[i] += 1;
				count[i + 1] += 1;
				count[i + 2] += 1;
			}
		}

		for i := 9; i < 16; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				if (tryHU(count, len - 3)) {
					return true;
				}
				count[i] += 1;
				count[i + 1] += 1;
				count[i + 2] += 1;
			}
		}

		for i := 18; i < 25; i++ {
			if (count[i] > 0 && count[i + 1] > 0 && count[i + 2] > 0) {
				count[i] -= 1;
				count[i + 1] -= 1;
				count[i + 2] -= 1;
				if (tryHU(count, len - 3)) {
					return true;

				}
				count[i] += 1;
				count[i + 1] += 1;
				count[i + 2] += 1;
			}
		}

	}
	return false;
}


//是否能碰牌
func (p *PairMjPai) IsPengPai(pai *MJPai) bool {
	return true
}

//是否能杠牌
func (p *PairMjPai) IsGangPai(Pai *MJPai) bool {
	return true
}

//确定要胡牌的时候,做出的处理
func (p *PairMjPai) HuPai(pai *MJPai) error {
	//排序
	//统计
	p.InitPaiStatistics()
	//p.IsDaDui        //判断是否是大队子
	//p.IsQingYiSe        //判断是否是清一色
	//p.IsJiangDui        //判断是否是将对
	//
	return nil

}

//碰牌
func (p *PairMjPai) CanPengPai(pai *MJPai) error {
	return nil
}


//杠牌
func (p *PairMjPai) CanGangPai(pai *MJPai) error {
	return nil
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
