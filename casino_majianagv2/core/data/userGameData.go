package data

import "casino_majiang/service/majiang"

type MJUserGameData struct {
	*majiang.PlayerGameData
}

//todo 删除过胡的信息
func (d *MJUserGameData) DelGuoHuInfo() {

}

//todo 删除摸牌前杠牌的信息
func (d *MJUserGameData) DelPreMoGangInfo() {

}
