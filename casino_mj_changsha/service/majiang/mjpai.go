package majiang

import (
	mjproto"casino_mj_changsha/msg/protogo"
	"casino_mj_changsha/msg/funcsInit"
	"errors"
	"casino_common/common/log"
	"casino_common/utils/numUtils"
)

//1=明杠、2=巴杠、3=暗杠
var GANG_TYPE_DIAN int32 = 1 //明杠
var GANG_TYPE_BA int32 = 2   //巴杠
var GANG_TYPE_AN int32 = 3   //暗杠
var GANG_TYPE_BU int32 = 4   //补杠	目前只有长沙麻将才有

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

func (p *MJPai) LogDes() string {
	if p == nil {
		return "空牌"
	}
	valueStr, _ := numUtils.Int2String(p.GetValue())
	idStr, _ := numUtils.Int2String(p.GetIndex())
	return idStr + "-" + valueStr + GetFlow(p.GetFlower())
}

//用户牌是否包含缺牌
func (p *MJHandPai) IsContainQue() bool {
	//判断inPai
	if p.InPai != nil {
		if p.GetQueFlower() == p.InPai.GetFlower() {
			return true
		}
	}

	//判断手牌是否缺牌
	for i := 0; i < len(p.Pais); i++ {
		if p.GetQueFlower() == p.Pais[i].GetFlower() {
			return true
		}
	}
	return false
}

func (p *MJHandPai) GetCanPeng(pai *MJPai) bool {
	//判断是否是缺牌
	if pai.GetFlower() == p.GetQueFlower() {
		return false
	}

	return CanPengPai(pai, p)
}

//判断手牌能否吃这张牌
func (p *MJHandPai) GetCanChi(pai *MJPai) (bool, []*MJPai) {
	//log.T("开始判断能否吃[%v]", pai.GetDes())
	var result []*MJPai
	var pai1, pai2 *MJPai

	canChi := false

	pai1 = p.IsExistHandPaisByValueAndFlower(*pai.Value-2, *pai.Flower);
	pai2 = p.IsExistHandPaisByValueAndFlower(*pai.Value-1, *pai.Flower);
	if pai1 != nil && pai2 != nil {
		//log.T("能吃前两位[%v] [%v]", pai1.GetDes(), pai2.GetDes())
		canChi = true
		result = append(result, pai1)
		result = append(result, pai)
		result = append(result, pai2)
	}

	pai1 = p.IsExistHandPaisByValueAndFlower(*pai.Value-1, *pai.Flower);
	pai2 = p.IsExistHandPaisByValueAndFlower(*pai.Value+1, *pai.Flower);
	if pai1 != nil && pai2 != nil {
		//log.T("能吃左右两位[%v] [%v]", pai1.GetDes(), pai2.GetDes())
		canChi = true
		result = append(result, pai1)
		result = append(result, pai)
		result = append(result, pai2)
	}

	pai1 = p.IsExistHandPaisByValueAndFlower(*pai.Value+1, *pai.Flower);
	pai2 = p.IsExistHandPaisByValueAndFlower(*pai.Value+2, *pai.Flower);
	if pai1 != nil && pai2 != nil {
		//log.T("能吃后两位[%v] [%v]", pai1.GetDes(), pai2.GetDes())
		canChi = true
		result = append(result, pai1)
		result = append(result, pai)
		result = append(result, pai2)
	}
	return canChi, result
}

func (p *MJHandPai) GetCanGang(pai *MJPai, remainPaiCount int32) (bool, []*MJPai) {

	//判断有没有剩余的牌
	if remainPaiCount == 0 {
		return false, nil
	}

	//自己手中的杠牌，只返回一张就可以了...
	boolReuls, pais := CanGangPai(pai, p)
	if pais != nil && len(pais) >= 4 {
		var result []*MJPai
		for _, pai := range pais {
			//杠牌中没有，并且手牌中有的才能放在list中...
			if !IsListExisGangPais(result, pai) && (p.IsExistHandPais(pai) || pai.GetIndex() == p.InPai.GetIndex()) {
				result = append(result, pai)
			}
		}
		return boolReuls, result
	} else {
		return boolReuls, pais
	}
}

func (p *MJHandPai) IsExistHandPais(pai *MJPai) bool {
	for _, p1 := range p.Pais {
		if p1 != nil && p1.GetIndex() == pai.GetIndex() {
			return true
		}
	}
	return false;
}

//根据牌值和牌花色从手牌中找到一样的牌
func (p *MJHandPai) IsExistHandPaisByValueAndFlower(value, flower int32) *MJPai {
	//log.T("开始从手牌[%v]中找是否同花色[%v] 同值的牌[%v]", p.Pais, flower, value)
	for i, _ := range p.Pais {
		if *p.Pais[i].Flower != flower {
			continue
		}
		if *p.Pais[i].Value == value {
			//log.T("已找到, 返回牌[%v]", p.Pais[i].GetDes())
			return p.Pais[i]
		}
	}
	//log.T("未找到, 返回nil")
	return nil
}

//从手牌中找到一样的牌
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
		hand.Pais = append(hand.Pais[:index], hand.Pais[index+1:]...)
		return nil

	} else {
		log.E("服务器错误：删除手牌的时候出错，没有找到对应的手牌[%v]", key)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}
}

//删除杠牌的信息
func (d *MJHandPai) DelGangInfo(pai *MJPai) error {
	index := -1
	for i, info := range d.gangInfos {
		if info != nil && info.GetPai().GetClientId() == pai.GetClientId() {
			index = i
			break
		}
	}

	if index > -1 {
		d.gangInfos = append(d.gangInfos[:index], d.gangInfos[index+1:]...)
		return nil
	} else {
		log.E("服务器错误：删除gangInfo的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}

	return nil
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
		hand.PengPais = append(hand.PengPais[:index], hand.PengPais[index+1:]...)
		return nil

	} else {
		log.E("服务器错误：删除碰牌的时候出错，没有找到对应的碰牌[%v]", key)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}
}

func (hand *MJHandPai) DelOutPai(key int32) error {
	index := -1
	for i, pai := range hand.OutPais {
		if pai != nil && pai.GetIndex() == key {
			index = i
			break
		}
	}
	if index > -1 {
		hand.OutPais = append(hand.OutPais[:index], hand.OutPais[index+1:]...)
		return nil

	} else {
		log.E("服务器错误：删除out牌的时候出错，没有找到对应的out牌[%v]", key)
		return errors.New("删除出牌时出错，没有找到对应的出牌...")
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
	return ServerPais2string(hand.Pais)
}

//手牌排序

type MjPaiList []*MJPai

func (list MjPaiList) Len() int {
	return len(list)

}

func (list MjPaiList) Less(i, j int) bool {
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
func (list MjPaiList) Swap(i, j int) {
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
		d.GangInfo = append(d.GangInfo[:index], d.GangInfo[index+1:]...)
		return nil
	} else {
		log.E("服务器错误：删除gangInfo的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌...")
	}

	return nil
}
