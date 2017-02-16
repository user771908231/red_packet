package skeleton

import (
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/majiangv2"
)

//胡牌的骨架，常用方法，公用方法的集合
type HuParserSkeleton struct {
}

//返回胡牌可以使用的通用方法
func NewHuParserSkeleton() *HuParserSkeleton {
	return &HuParserSkeleton{}
}

//清一色
func (p *HuParserSkeleton) IsQingYiSe(pais []*majiang.MJPai) bool {
	flower := pais[0].Flower
	for i := 1; i < len(pais); i++ {
		if *flower != *pais[i].Flower {
			return false //不是清一色
		}
	}

	return true
}

//七对
func (p *HuParserSkeleton) IsQiDui(handCounts [] int) bool {
	countDuizi := 0
	for i := 0; i < len(handCounts); i++ {
		if (handCounts [i] != 2) && (handCounts[i] != 0) {
			//每张牌都是2张
			return false
		} else {
			countDuizi ++
		}
	}
	return countDuizi == 7
}

//龙七对
func (p *HuParserSkeleton) IsLongQiDui(handCounts [] int) bool {
	longCount := 0 //杠数
	duiCount := 0  //对数

	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 4 {
			//杠
			longCount++
		}
		if (handCounts[i] == 2) {
			duiCount++
		}
		if (handCounts[i] != 0) && (handCounts[i] != 4) && (handCounts[i] != 2) {
			//牌数不符合
			//log.T("isLongQiDui: 牌数不符合 0、2、4")
			return false
		}
	}
	if (longCount < 1) || (duiCount < 5) {
		//杠数小于一，对数小于5
		//log.T("isLongQiDui: 杠对数不符合")
		return false
	}
	return true
}

//七对 龙七对牌型胡牌判断
func (p *HuParserSkeleton) TryHU7(handCounts [] int) bool {
	if p.IsQiDui(handCounts) || p.IsLongQiDui(handCounts) {
		return true
	} else {
		return false
	}
}

func (p *HuParserSkeleton) is19(val int) bool {
	v := majiangv2.GettPaiValueByCountPos(val)
	return v == 1 || v == 9
	//return (val % 9 == 0) || (val % 9 == 8)
}

func (p *HuParserSkeleton) Is258(val int) bool {
	v := majiangv2.GettPaiValueByCountPos(val)
	return v == 2 || v == 5 || v == 8
}

func (p *HuParserSkeleton) TryHU(count []int, len int) (result bool, isAll19 bool, jiang int) {
	//log.T("开始判断tryHu(%v,%v)", count, len)
	isAll19 = true //全带幺
	result = false
	jiang = -1
	//递归完所有的牌表示 胡了
	if (len == 0) {
		//log.T("len == 0")
		return true, isAll19, jiang
	}

	if (len%3 == 2) {
		//log.T("if %v 取模 3 == 2", len)
		// 说明对牌出现
		for i := 0; i < 27; i++ {
			if (count[i] >= 2) {
				count[i] -= 2

				result, isAll19, jiang = p.TryHU(count, len-2)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if ! p.is19(i) {
						//不是幺九
						isAll19 = false
					}
					jiang = i
					return true, isAll19, jiang
				}
				count[i] += 2
			}
		}
	} else {
		//log.T("else %v", len)
		// 是否是顺子，这里应该分开判断
		for i := 0; i < 7; i++ {
			if (count[i] > 0 && count[i+1] > 0 && count[i+2] > 0) {
				count[i] -= 1;
				count[i+1] -= 1;
				count[i+2] -= 1;
				result, isAll19, jiang = p.TryHU(count, len-3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !p.is19(i) && !p.is19(i + 1) && !p.is19(i + 2) {
						//不是幺九
						//log.T("branch 2 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19, jiang
				}
				count[i] += 1
				count[i+1] += 1
				count[i+2] += 1
			}
		}

		for i := 9; i < 16; i++ {
			if (count[i] > 0 && count[i+1] > 0 && count[i+2] > 0) {
				count[i] -= 1
				count[i+1] -= 1
				count[i+2] -= 1
				result, isAll19, jiang = p.TryHU(count, len-3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !p.is19(i) && !p.is19(i + 1) && !p.is19(i + 2) {
						//不是幺九
						//log.T("branch 3 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19, jiang
				}
				count[i] += 1
				count[i+1] += 1
				count[i+2] += 1
			}
		}

		for i := 18; i < 25; i++ {
			if (count[i] > 0 && count[i+1] > 0 && count[i+2] > 0) {
				count[i] -= 1;
				count[i+1] -= 1;
				count[i+2] -= 1;
				result, isAll19, jiang = p.TryHU(count, len-3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !p.is19(i) && !p.is19(i + 1) && !p.is19(i + 2) {
						//不是幺九
						//log.T("branch 4 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19, jiang
				}
				count[i] += 1;
				count[i+1] += 1;
				count[i+2] += 1;
			}
		}

		// 三个一样的
		for i := 0; i < 27; i++ {
			if (count[i] >= 3) {
				count[i] -= 3
				result, isAll19, jiang = p.TryHU(count, len-3)
				if (result) {
					//log.T("i: %v, value: %v", i, count[i])
					if !p.is19(i) {
						//不是幺九
						//log.T("branch 5 pos%v不是幺九", i)
						isAll19 = false
					}
					return true, isAll19, jiang
				}
				count[i] += 3
			}
		}
	}
	return false, isAll19, jiang
}
