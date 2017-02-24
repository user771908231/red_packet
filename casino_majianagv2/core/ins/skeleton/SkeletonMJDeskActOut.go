package skeleton

import "casino_majianagv2/core/api"

func (d *SkeletonMJDesk) ActOut(userId uint32, cardId int32, auto bool) error {
	return nil
}

//得到下一个摸牌的人...
func (d *SkeletonMJDesk) GetNextMoPaiUser() api.MjUser {

	//首先找，刚刚杠过牌的User，否则找下一个User
	for _, u := range d.GetSkeletonMJUsers() {
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
		//log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v],循环检测index[%v],user.IsNotHu(%v),user.CanMoPai[%v]", d.GetActiveUser(), activeIndex, i, user.IsNotHu(), user.CanMoPai(d.IsXueLiuChengHe()))
		if user != nil && user.GetStatus().CanMoPai(d.IsXueLiuChengHe()) {
			activeUser = user
			break
		}
	}
	//找到下一个操作的user
	return activeUser
}
