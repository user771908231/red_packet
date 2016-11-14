package doudizhu

import (
	"casino_server/utils/pokerUtils"
	"errors"
	"casino_server/common/log"
	"sort"
	"github.com/derekparker/delve/dwarf/line"
)

//返回一张牌
func InitPaiByIndex(index int32) *PPokerPai {
	_, rmapdes, pvalue, pflower, pname := pokerUtils.ParseByIndex(index)
	//返回一张需要的牌
	pokerPai := NewPPokerPai()
	*pokerPai.Id = index
	*pokerPai.Name = pname
	*pokerPai.Value = pvalue
	*pokerPai.Flower = pflower
	*pokerPai.Des = rmapdes
	return pokerPai
}

//喜好衣服扑克牌
func XiPai() []*PPokerPai {
	//得到随机的牌的index...
	randIndex := pokerUtils.Xipai(54)

	//通过index 得到牌

	var pokerPais []*PPokerPai
	for _, i := range randIndex {
		pai := InitPaiByIndex(i)
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
func (out *POutPokerPais) getPaiCount() int {
	return len(out.PokerPais)
}


//对牌进行排序,左边小，右边大的值进行排序
func (out *POutPokerPais) sortPais() error {
	list := out.PokerPais
	log.T("befor sort : %v", list)
	sort.Sort(list)        //进行排序
	log.T("after sort : %v", out.PokerPais)

	return nil
}


//初始化类型
func (out *POutPokerPais) initTypeAndKeyValue() error {
	//统计数据
	counts := make([]int32, 15)
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
		for k := 0; k < out.getPaiCount() - 1; k++ {
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

	isFeijiChibang := false
	if out.GetCountSanzhang() * 4 == out.getPaiCount() {
		boolFlag := true
		for i := 0; i < out.GetCountSanzhang() - 1; i++ {
			if countsSanzhang[i] + 1 != countsSanzhang[i + 1] {
				boolFlag = false
				break
			}
		}
		isFeijiChibang = boolFlag
	}

	isLianDui := false
	if out.GetCountDuizi() * 2 == out.getPaiCount() {
		boolFlag := true
		for i := 0; i < out.GetCountDuizi() - 1; i++ {
			if countsLiagnzhang[i] + 1 != countsLiagnzhang[i + 1] {
				boolFlag = false
				break
			}
		}
		isLianDui = boolFlag
	}

	//飞机不带翅膀

	//四个带两个单牌

	//四个带两个对子


	//判断是否是单张
	if out.getPaiCount() == 1 {
		//out.Type = 单张
		*out.KeyValue = out.GetPokerPais()[0].GetValue()
	} else if out.getPaiCount() == 2 {
		if out.GetCountDuizi() == 1 {
			//这里需要判断是否是王炸
			if out.GetPokerPais()[0].GetValue() == -1 {
				//判断是否是王炸
				//out.Type = 炸弹
				*out.KeyValue = out.GetPokerPais()[0].GetValue()
			} else {
				//out.Type = 对子
				*out.KeyValue = out.GetPokerPais()[0].GetValue()
			}
		} else {
			//error
		}
	} else if out.getPaiCount() == 3 {
		if out.GetCountSanzhang() == 1 {
			//out.Type == 三张
			*out.KeyValue = countsSanzhang[0]
		} else {
			//error
		}
	} else if out.getPaiCount() == 4 {
		if out.GetCountSizhang() == 1 {
			//out.Type == 炸弹
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSanzhang() == 1 && out.GetCountYizhang() == 1 {
			//out.Type == 三带一
			*out.KeyValue = countsSanzhang[0]

		} else {
			//	error
		}
	} else if out.getPaiCount() == 5 {
		if isShunZi {
			//out.Type == 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSanzhang() == 1 && out.GetCountDuizi() == 1 {
			//out.Type == 三带二
			*out.KeyValue = countsSanzhang[0]
		} else {
			//error
		}

	} else if out.getPaiCount() == 6 {
		if isShunZi {
			//out.Type == 顺子
		} else if out.GetCountSizhang() == 1 && out.GetCountYizhang() == 2 {
			//out.Type == 四带二
		} else {
			//error
		}

	} else if out.getPaiCount() == 8 {
		if isShunZi {
			//out.Type == 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if out.GetCountSizhang() == 1 && out.GetCountDuizi() == 2 {
			//out.Type == 四带两队
			*out.KeyValue = countsSizhang[0]
		} else {
			//error
		}

	} else {
		if isShunZi {
			//out.Type == 顺子
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if isFeiji {
			//out.Type == 飞机
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		} else if isFeijiChibang {
			//out.Type = 飞机带翅膀
			*out.KeyValue = countsSanzhang[0]
		} else if isLianDui {
			//out.Type = 连队
			*out.KeyValue = out.GetPokerPais()[0].GetValue()
		}
	}

	return nil;
}
