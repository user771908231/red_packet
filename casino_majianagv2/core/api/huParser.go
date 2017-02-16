package api

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
)

//胡牌解析器
type HuPaerApi interface {
	GetCanHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai, iszimo bool, huType mjproto.HuType) (bool, int32, int64, [] string, []mjproto.PaiType, bool) //是否能胡牌,返回是否能胡，翻数，分数,huCardStr,paiType
	HuScore() (fan int32, score int64)                                                                                                                      //得到胡牌的翻数和分数
	GetJiaoPais(pais []*majiang.MJPai) []*majiang.MJPai
}
