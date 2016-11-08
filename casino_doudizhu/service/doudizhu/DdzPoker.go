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

