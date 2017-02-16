package data

import "casino_majiang/service/majiang"

type MJUserGameData struct {
	*majiang.PlayerGameData
	preMoGangInfo *majiang.GangPaiInfo
}

//todo 删除过胡的信息
func (d *MJUserGameData) DelGuoHuInfo() {

}

//返回摸前的杠牌info
func (d *MJUserGameData) GetPreMoGangInfo() *majiang.GangPaiInfo {
	return d.preMoGangInfo
}

//todo 删除摸牌前杠牌的信息
func (d *MJUserGameData) DelPreMoGangInfo() {

}

//增加胡牌的信息
func (d *MJUserGameData) AddHuPaiInfo(hu *majiang.HuPaiInfo) {
	d.HuInfo = append(d.HuInfo, hu)
	d.HandPai.HuPais = append(d.HandPai.HuPais, hu.Pai) //增加胡牌
}
