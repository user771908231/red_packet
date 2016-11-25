package doudizhu

import "casino_common/common/Error"


//四川斗地主的逻辑都在这个文件里面

//四川叫地主
/**
	四川地主，当玩家叫地主之后，不再继续抢地主
 */
func (d *DdzDesk) SCJiaoDiZhu(userId uint32) error {
	//验证活动玩家
	err := d.CheckActiveUser(userId)
	if err != nil {
		return err
	}

	//验证用户是否为空
	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewFailError("玩家没找到，抢地主失败")
	}

	//抢地主
	d.SetDizhu(user.GetUserId())
	user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_QIANG)
	//表示地主都抢过来，又轮到第一家，抢地主的逻辑结束
	if user.IsQiangDiZhuQiang() && d.GetDizhuPaiUser() == user.GetUserId() {
		//todo 抢地主结束的操作
		return nil
	}

	//todo 抢地主结束的操作
	d.afterQiangDizhu()

	return nil
}

