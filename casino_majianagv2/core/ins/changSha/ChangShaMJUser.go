package changSha

import (
	"casino_majiang/service/majiang"
	"casino_common/common/log"
	"casino_majianagv2/core/ins/skeleton"
)

type ChangShaMJUser struct {
	*skeleton.SkeletonMJUser
	changshaGang bool
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
	jiaoPais := u.GetDesk().GetHuParser().GetJiaoPais(newPais)
	if jiaoPais != nil && len(jiaoPais) > 0 {
		return true
	}
	return false

}

//目前主要是长沙麻将使用
func (u *ChangShaMJUser) IsCanInitCheckCaseChi() bool {
	//由于湖南麻将是倒倒胡，所以不用考虑中间间隔胡牌玩家的情况
	//log.T("判断能不能吃..u.d.getIndexByUserId(u.d.GetCheckCase().GetUserIdOut()) %v", u.d.getIndexByUserId(u.d.GetCheckCase().GetUserIdOut()))
	//log.T("判断能不能吃..u.d.GetUserCountLimit() %v", u.d.GetUserCountLimit())
	//log.T("判断能不能吃..u.d.getIndexByUserId(u.GetUserId()) %v", u.d.getIndexByUserId(u.GetUserId()))
	if (u.GetSkeletonMJDesk().GetIndexByUserId(u.GetSkeletonMJDesk().GetCheckCase().GetUserIdOut())+1)%int(u.GetSkeletonMJDesk().GetMJConfig().PlayerCountLimit) == u.GetSkeletonMJDesk().GetIndexByUserId(u.GetUserId()) {
		return true
	}
	return false
}
