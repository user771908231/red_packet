package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
)

//得到一张牌的信息
func (p *MJPai) GetCardInfo() *mjproto.CardInfo {
	cardInfo := newProto.NewCardInfo()
	*cardInfo.Id = p.GetIndex()
	*cardInfo.Type = p.GetFlower()
	*cardInfo.Value = p.GetValue()
	return cardInfo
}

//返回前端需要的id号
func (p *MJPai) GetClientId() int32 {
	return 0
}

//是否可以胡牌
func (p *MJHandPai) GetCanHu() bool {
	return false
}

func (p *MJHandPai) GetCanPeng() bool {
	return false
}

func (p *MJHandPai) GetCanGang() bool {
	return false
}
