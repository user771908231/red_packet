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

	//判断desk是否存在
	desk := room.GetDeskByKey(key)
	if desk == nil {
		//返回失败的信息
		return errors.New("进入房间失败,没有找到合适的desk")
	}

	//进入房间
	err := desk.EnterUser(userId)
	if err != nil {
		log.E("玩家[%v]进入desk[%v]失败err[%v]", userId, desk.GetDeskId(), err)
		return Error.NewError(-1, "玩家进入desk失败.")
	}

	return nil
}

//准备游戏
func HandlerFDdzReady(user uint32) error {
	desk := doudizhu.GetDdzDeskBySession(user)
	if desk == nil {
		log.E("玩家[%v]准备失败，没有找到对应的desk..")
		return Error.NewError(-1, "没有找到玩家所在的desk")
	}

	err := desk.Ready(user)
	if err != nil {
		log.E("玩家[%v]准备游戏的时候失败...")
		return Error.NewError(-1, "玩家准备游戏的时候失败...")
	}

	return nil
}

//开始出牌
func HandlerActOut(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewError(-1, "没有找到desk")
	}

	//获取到outpai
	outpai := doudizhu.GetOutPais()
	//开始出牌
	err := desk.ActOut(userId, outpai)
	if err != nil {
		log.E("出牌的时候失败...")
	}
	return nil
}
