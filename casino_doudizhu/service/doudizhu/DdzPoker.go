package doudizhu

import "casino_server/utils/pokerUtils"

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
