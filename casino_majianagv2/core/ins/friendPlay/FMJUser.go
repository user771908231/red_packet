package friendPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/majiang"
	"casino_common/common/log"
	"casino_majianagv2/core/api"
	"github.com/name5566/leaf/gate"
)

type FMJUser struct {
	*skeleton.SkeletonMJUser
}

func NewFMJUser(desk api.MjDesk, userId uint32, a gate.Agent) *FMJUser {
	//骨架User
	suser := skeleton.NewSkeletonMJUser(desk, userId, a) //朋友桌user
	suser.Coin = 0                                       //默认金币是0
	suser.ActTimeoutCount = 0
	return &FMJUser{
		SkeletonMJUser: suser,
	}
}

func (u *FMJUser) SendOverTurn(p proto.Message) error {
	//如果是金币场有超时的处理...
	u.WriteMsg(p)
	return nil
}

//得到判定bean
func (u *FMJUser) GetCheckBean(p *majiang.MJPai, xueliuchenghe bool, remainPaiCoun int32) *majiang.CheckBean {
	bean := majiang.NewCheckBean()

	*bean.CheckStatus = majiang.CHECK_CASE_BEAN_STATUS_CHECKING
	*bean.UserId = u.GetUserId()
	bean.CheckPai = p

	//是否可以胡牌
	if u.IsCanInitCheckCaseHu(xueliuchenghe) {
		*bean.CanHu, _, _, _, _, _ = u.GetDesk().GetHuParser().GetCanHu(u.GameData.HandPai, p, false, 0)
	}
	//是否可以杠
	if u.IsCanInitCheckCaseGang(xueliuchenghe) {
		*bean.CanGang, _ = u.GameData.HandPai.GetCanGang(p, remainPaiCoun)
	}
	//是否可以碰
	if u.IsCanInitCheckCasePeng() {
		*bean.CanPeng = u.GameData.HandPai.GetCanPeng(p)
	}

	//是否可以吃牌
	if u.IsCanInitCheckCaseChi() {
		*bean.CanChi, bean.ChiCards = u.GameData.HandPai.GetCanChi(p)
	}

	log.T("得到用户[%v]对牌[%v]的check , bean[%v]", u.GetUserId(), p.LogDes(), bean)
	//判断过胡.如果有过胡，那么就不能再胡了
	if u.HadGuoHuInfo(p) {
		*bean.CanHu = false
	}

	if bean.GetCanGang() || bean.GetCanHu() || bean.GetCanPeng() || bean.GetCanChi() {
		return bean
	} else {
		return nil
	}
}

//是否已经有过胡了
func (u *FMJUser) HadGuoHuInfo(pai *majiang.MJPai) bool {
	if u.GameData.GuoHuInfo == nil || len(u.GameData.GuoHuInfo) <= 0 {
		return false
	}

	//目前只做成  牌一样的时候再判断

	//如果huinfo的牌和pai 一样，表示有guohu的info
	for _, info := range u.GameData.GuoHuInfo {
		if pai.GetClientId() == info.GetPai().GetClientId() {
			return true
		}
	}

	//没有过胡的信息
	return false
}

//判断用户是否可以杠
func (u *FMJUser) IsCanInitCheckCaseGang(xueliuchenghe bool) bool {
	//这里需要判断是否是 血流成河，目前暂时不判断...

	//1,普通规则
	if u.GetStatus().IsNotHu() {
		return true
	}

	//2,血流成河
	if u.GetStatus().IsHu() && xueliuchenghe {
		return true
	}

	//其他情况返回false
	return false
}

func (u *FMJUser) IsCanInitCheckCasePeng() bool {
	//1,普通规则
	if u.GetStatus().IsNotHu() {
		return true
	} else {
		return false;
	}
}

//目前主要是长沙麻将使用
func (u *FMJUser) IsCanInitCheckCaseChi() bool {
	return false
}

//判断用户是否可以杠
func (u *FMJUser) IsCanInitCheckCaseHu(xueliuchenghe bool) bool {
	return u.IsCanInitCheckCaseGang(xueliuchenghe)
}
