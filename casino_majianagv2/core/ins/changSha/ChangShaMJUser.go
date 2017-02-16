package changSha

import (
	"casino_majiang/service/majiang"
	"casino_common/common/log"
	"casino_majianagv2/core/ins/skeleton"
)

type ChangShaMJUser struct {
	*skeleton.SkeletonMJUser
}

/**
	补牌的逻辑问题
 */
func (u *ChangShaMJUser) GetCanChangShaGang(pai *majiang.MJPai) bool {
	log.T("桌子[%v]玩家[%v]开始判断是否能长沙杠pai[%v]", u.GetDesk().GetMJConfig().DeskId, u.GetUserId(), pai.LogDes())
	if u.GameData == nil || u.GameData.HandPai == nil || u.GameData.HandPai.Pais == nil || len(u.GameData.HandPai.Pais) <= 0 {
		log.T("桌子[%v]玩家[%v]判断是否能杠错误 玩家手牌为空", u.GetDesk().GetMJConfig().DeskId, u.GetUserId())
		return false
	}
	//1 手牌中和pai 花色相同的去掉 by 彬哥
	newPais := []*majiang.MJPai{}
	pais := u.GameData.HandPai.Pais
	for i := 0; i < len(pais); i++ {
		fp := pais[i]
		if fp != nil && fp.GetClientId() != pai.GetClientId() {
			newPais = append(newPais, fp)
		}
	}

	//添加in牌
	inpai := u.GetGameData().GetHandPai().GetInPai()
	if inpai != nil && inpai.GetClientId() != pai.GetClientId() {
		newPais = append(newPais, inpai)
	}

	//2 看剩下的牌有无听 by 彬哥
	jiaoPais := u.GetDesk().GetSkeletonMJDesk().HuParser.GetJiaoPais(newPais)
	if jiaoPais != nil && len(jiaoPais) > 0 {
		return true
	}
	return false

}
