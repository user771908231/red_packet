package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"errors"
	"casino_server/common/log"
	"casino_server/utils/numUtils"
	"sort"
)

//1=明杠、2=巴杠、3=暗杠
var GANG_TYPE_DIAN int32 = 1//明杠
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
func (p *MJHandPai) GetCanHu() (bool,bool) {
	return CanHuPai(p)
}

func (p *MJHandPai) GetCanPeng(pai *MJPai) bool {
	return CanPengPai(pai, p)
}

func (p *MJHandPai) GetCanGang(pai *MJPai) (bool, []*MJPai) {
	//自己手中的杠牌，只返回一张就可以了...
	var result []*MJPai
	boolReuls, pais := CanGangPai(pai, p)
	if pais != nil && len(pais) >= 4 {
		for _, pai := range pais {
			//杠牌中没有，并且手牌中有的才能放在list中...
			if !IsListExisGangPais(result, pai) && (p.IsExistHandPais(pai) || pai.GetIndex() == p.InPai.GetIndex()) {
				result = append(result, pai)
			}
		}
	}
	return boolReuls, result
}

func (p *MJHandPai) IsExistHandPais(pai *MJPai) bool {
	for _, p1 := range p.Pais {
		if p1 != nil && p1.GetIndex() == pai.GetIndex() {
			return true
		}
	}
	return false;
}

func IsListExisGangPais(ps []*MJPai, p *MJPai) bool {
	for _, pss := range ps {
		if pss.GetClientId() == p.GetClientId() {
			return true
		}
	}

	return false
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
			if p.GetClientId() == pai.GetClientId() {
				//表示花色相同，有碰牌
				return true
			}
		}
	}
	return false;
}

//得到手牌的描述
func (hand *MJHandPai) GetDes() string {
	var tempPais MjPaiList = make([]*MJPai, len(hand.Pais))
	copy(tempPais, hand.Pais)
	sort.Sort(tempPais)

	s := ""
	for _, p := range tempPais {
		if p != nil {
			s += (" " + p.LogDes() + " ")
		}
	}
	return s
}

//手牌排序

type MjPaiList []*MJPai

func (list MjPaiList)Len() int {
	return len(list)

}

func ( list MjPaiList)Less(i, j int) bool {
	if list[i].GetFlower() < list[j].GetFlower() {
		return true
	} else if list[i].GetFlower() == list[j].GetFlower() {
		if list[i].GetValue() < list[j].GetValue() {
			return true
		} else {
			return false
		}

	} else {
		return false
	}

}
func (list MjPaiList)Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
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