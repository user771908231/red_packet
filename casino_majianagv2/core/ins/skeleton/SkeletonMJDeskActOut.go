package skeleton

import (
	"casino_majianagv2/core/api"
	"casino_majiang/service/majiang"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
)

func (d *SkeletonMJDesk) ActOut(userId uint32, cardId int32, auto bool) error {
	return nil
}

//得到下一个摸牌的人...
func (d *SkeletonMJDesk) GetNextMoPaiUser() api.MjUser {

	//首先找，刚刚杠过牌的User，否则找下一个User
	for _, u := range d.GetUsers() {
		if u != nil && u.GetGameData().GetPreMoGangInfo() != nil {
			return u
		}
	}

	//log.T("查询下一个玩家...当前的activeUser[%v]", d.GetActiveUser())
	var activeUser api.MjUser = nil
	activeIndex := -1
	for i, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == d.GetMJConfig().ActiveUser {
			activeIndex = i
			break
		}
	}
	//log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v]", d.GetActiveUser(), activeIndex)
	if activeIndex == -1 {
		return nil
	}

	for i := activeIndex + 1; i < activeIndex+int(d.GetUserCount()); i++ {
		user := d.GetUsers()[i%int(d.GetUserCount())]
		if user != nil && user.GetStatus().CanMoPai(d.IsXueLiuChengHe()) {
			activeUser = user
			break
		}
	}
	//找到下一个操作的user
	return activeUser
}

//初始化checkCase
//如果出错 设置checkCase为nil
func (d *SkeletonMJDesk) InitCheckCase(p *majiang.MJPai, outUser api.MjUser) error {

	checkCase := majiang.NewCheckCase()
	checkCase.DianPaoCount = proto.Int32(0) //设置点炮的次数为0
	*checkCase.UserIdOut = outUser.GetUserId()
	*checkCase.CheckStatus = majiang.CHECK_CASE_STATUS_CHECKING //正在判定
	checkCase.CheckMJPai = p
	checkCase.PreOutGangInfo = outUser.GetGameData().GetPreMoGangInfo()
	d.CheckCase = checkCase

	//初始化checkbean
	for _, checkUser := range d.GetSkeletonMJUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			//添加checkBean
			bean := checkUser.GetCheckBean(p, d.GetRemainPaiCount())
			if bean != nil {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
	}

	log.T("判断最终的checkCase[%v]", checkCase)
	if checkCase.CheckB == nil || len(checkCase.CheckB) == 0 {
		d.CheckCase = nil
	}

	return nil
}
