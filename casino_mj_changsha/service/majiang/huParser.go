package majiang

import (
	mjproto "casino_mj_changsha/msg/protogo"
)

//胡牌解析起的接口
type HuParser interface {
	GetCanHu(handPai *MJHandPai, hupai *MJPai, iszimo bool, huType mjproto.HuType, isBanker bool) (bool, int32, int64, [] string, []mjproto.PaiType, bool) //是否能胡牌,返回是否能胡，翻数，分数,huCardStr,paiType
	HuScore() (fan int32, score int64)                                                                                                      //得到胡牌的翻数和分数
	GetJiaoPais(handPai *MJHandPai,) []*MJPai                                                                                                     //得到叫牌
}

//胡牌的骨架，常用方法，公用方法的集合
type HuParserSkeleton struct {
}

//返回胡牌可以使用的通用方法
func NewHuParserSkeleton() *HuParserSkeleton {
	return &HuParserSkeleton{}
}

//清一色
func (p *HuParserSkeleton) IsQingYiSe(pais []*MJPai) bool {
	flower := pais[0].Flower
	for i := 1; i < len(pais); i++ {
		if *flower != *pais[i].Flower {
			return false //不是清一色
		}
	}

	return true
}


//判断是否是七对 返回勾数
func (p *HuParserSkeleton) IsQiDui(handPai *MJHandPai, hupai *MJPai) (isQidui bool, gou int32) {
	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)

	//合并手牌 长度不为14即有碰杠吃 不满足七对基本要求
	if len(joinedPais) != 14 {
		return false, 0
	}

	handCounts := GettPaiStats(joinedPais)
	for i := 0; i < len(handCounts); i++ {
		switch handCounts[i] {
		case 0, 2:
			continue //牌数为0和2都忽略
		case 4:
			gou++ //统计勾数
		default:
			return false, 0 //牌数1、3不为七对
		}
	}
	return true, gou
}

//七对 龙七对牌型胡牌判断
func (p *HuParserSkeleton) tryHU7(handPai *MJHandPai, hupai *MJPai) bool {
	isQidui, _ := p.IsQiDui(handPai, hupai)
	return isQidui
}

func (p *HuParserSkeleton) is19(val int) bool {
	v := GettPaiValueByCountPos(val)
	return v == 1 || v == 9
	//return (val % 9 == 0) || (val % 9 == 8)
}

func (p *HuParserSkeleton) is258(val int) bool {
	v := GettPaiValueByCountPos(val)
	return v == 2 || v == 5 || v == 8
}

func (p *HuParserSkeleton) tryHU(count []int, len int) (result bool, isAll19 bool, jiang int) {
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

				result, isAll19, jiang = p.tryHU(count, len-2)
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
				result, isAll19, jiang = p.tryHU(count, len-3)
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
				result, isAll19, jiang = p.tryHU(count, len-3)
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
				result, isAll19, jiang = p.tryHU(count, len-3)
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
				result, isAll19, jiang = p.tryHU(count, len-3)
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





////////////////////////////////////////////////////////////////////////

//将手牌碰牌刚牌与huPai拼接成数组
func (p *HuParserSkeleton) JoinAllHandPaiAndHuPai(handPai *MJHandPai, hupai *MJPai) []*MJPai {
	pais := []*MJPai{}
	if hupai != nil {
		pais = append(pais, hupai)
	}

	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}

	if handPai.GangPais != nil {
		pais = append(pais, handPai.GangPais...)
	}

	if handPai.PengPais != nil {
		pais = append(pais, handPai.PengPais...)
	}

	if handPai.ChiPais != nil {
		pais = append(pais, handPai.ChiPais...)
	}
	return pais
}

//将手牌与huPai拼接成数组
func (p *HuParserSkeleton) JoinHandPaiPaisAndHuPai(handPai *MJHandPai, hupai *MJPai) []*MJPai {
	pais := []*MJPai{}
	if hupai != nil {
		pais = append(pais, hupai)
	}

	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}

	return pais
}

//判断手牌中是否有杠牌
func (p *HuParserSkeleton) isYouGang(handPai *MJHandPai) bool {
	return handPai.GangPais != nil && len(handPai.GangPais) > 0
}

//判断手牌中是否有碰牌
func (p *HuParserSkeleton) isYouPeng(handPai *MJHandPai) bool {
	return handPai.PengPais != nil && len(handPai.PengPais) > 0
}

//判断手牌中是否有吃牌
func (p *HuParserSkeleton) isYouChi(handPai *MJHandPai) bool {
	return handPai.ChiPais != nil && len(handPai.ChiPais) > 0
}