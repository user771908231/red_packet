package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"errors"
	"casino_server/common/log"
	"casino_server/utils/numUtils"
)

//1=明杠、2=巴杠、3=暗杠
var GANG_TYPE_MING int32 = 1//明杠
var GANG_TYPE_BA int32 = 2//明杠
var GANG_TYPE_AN int32 = 3//明杠

//得到一张牌的信息
func (p *MJPai) GetCardInfo() *mjproto.CardInfo {
	cardInfo := newProto.NewCardInfo()
	*cardInfo.Id = p.GetIndex()
	*cardInfo.Type = p.GetFlower()
	*cardInfo.Value = p.GetClientId()
	return cardInfo
}

func (p *MJPai) GetBackPai() *mjproto.CardInfo {
	return NewBackPai()
}

//返回前端需要的id号
func (p *MJPai) GetClientId() int32 {
	return clienMap[int(p.GetIndex())]
}

func ( p *MJPai) LogDes() string {
	valueStr, _ := numUtils.Int2String(p.GetValue())
	return valueStr + GetFlow(p.GetFlower())
}


//是否可以胡牌
func (p *MJHandPai) GetCanHu() bool {
	return CanHuPai(p)
}

func (p *MJHandPai) GetCanPeng(pai *MJPai) bool {
	return CanPengPai(pai, p)
}

func (p *MJHandPai) GetCanGang(pai *MJPai) (bool, []*MJPai) {
	return CanGangPai(pai, p)
}

//增加一张牌
func (hand *MJHandPai) AddPai(pai *MJPai) error {
	hand.Pais = append(hand.Pais, pai)
	return nil
}

func (hand *MJHandPai) DelHandlPai(key int32) error {
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
		log.E("服务器错误：删除手牌的时候出错，没有找到对应的手牌[%v]", key)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}
}

func (hand *MJHandPai) DelPengPai(key int32) error {
	index := -1
	for i, pai := range hand.PengPais {
		if pai != nil && pai.GetIndex() == key {
			index = i
			break
		}
	}
	if index > -1 {
		hand.PengPais = append(hand.PengPais[:index], hand.PengPais[index + 1:]...)
		return nil

	} else {
		log.E("服务器错误：删除碰牌的时候出错，没有找到对应的碰牌[%v]", key)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}
}

//判断碰牌中是否有指定的牌
func (hand *MJHandPai) IsExistPengPai(pai *MJPai) bool {
	for _, p := range hand.PengPais {
		if p != nil {
			if p.GetValue() == pai.GetValue() && p.GetFlower() == pai.GetFlower() {
				//表示花色相同，有碰牌
				return true
			}
		}
	}
	return false;
}

//删除杠牌的信息
func (d *PlayerGameData) DelGangInfo(pai *MJPai) error {
	index := -1
	for i, info := range d.GangInfo {
		if info != nil && info.GetPai().GetClientId() == pai.GetClientId() {
			index = i
			break
		}
	}

	if index > -1 {
		d.GangInfo = append(d.GangInfo[:index], d.GangInfo[index + 1:]...)
		return nil
	} else {
		log.E("服务器错误：删除gangInfo的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}

	return nil
}