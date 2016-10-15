package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"errors"
	"casino_server/common/log"
)

//得到一张牌的信息
func (p *MJPai) GetCardInfo() *mjproto.CardInfo {
	cardInfo := newProto.NewCardInfo()
	*cardInfo.Id = p.GetIndex()
	*cardInfo.Type = p.GetFlower()
	*cardInfo.Value = p.GetClientId()
	return cardInfo
}

func (p *MJPai) GetBackPai() *mjproto.CardInfo {
	cardInfo := newProto.NewCardInfo()
	*cardInfo.Id = 0
	*cardInfo.Type = 0
	*cardInfo.Value = 0
	return cardInfo
}

//返回前端需要的id号
func (p *MJPai) GetClientId() int32 {
	return clienMap[int(p.GetIndex())]
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

//增加一张牌
func (hand *MJHandPai) AddPai(pai *MJPai) error {
	hand.Pais = append(hand.Pais, pai)
	return nil
}

func (hand *MJHandPai) DelPai(key int32) error {
	index := -1
	for i, pai := range hand.Pais {
		if pai != nil && pai.GetIndex() == key {
			index = i
			break
		}
	}
	if index > -1 {
		hand.Pais = append(hand.Pais[:index], hand.Pais[index + 1:]...)
		return nil

	} else {
		log.E("服务器错误：删除手牌的时候出错，没有找到对应的手牌")
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}
}
