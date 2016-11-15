package DdzService

import (
	"casino_server/common/log"
	"casino_doudizhu/service/doudizhu"
	"errors"
	"casino_server/common/Error"
)

/*
	1,gate 作为路由入口
	2,service作为逻辑入口

*/
//进入房间的逻辑
func HandlerEnterDesk(userId uint32, key string, deskType int32) error {
	log.T("玩家[%v]进入斗地主的房间。", userId)
	//todo 目前只做朋友桌
	room := doudizhu.GetFDdzRoom()
	if room == nil {
		//返回失败的信息
		return errors.New("进入房间失败...")
	}

	desk := room.GetDeskByKey(key)
	if desk == nil {
		//返回失败的信息
		return errors.New("进入房间失败,没有找到合适的desk")
	}

	err := desk.EnterUser(userId)
	if err != nil {
		log.E("玩家[%v]进入desk[%v]失败err[%v]", userId, desk.GetDeskId(), err)
		return Error.NewError(-1, "玩家进入desk失败.")
	}

	return nil
}
