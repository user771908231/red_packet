package DdzService

import (
	"casino_server/common/log"
	"casino_doudizhu/service/doudizhu"
	"errors"
	"casino_server/common/Error"
	"github.com/name5566/leaf/gate"
)

/*
	1,gate 作为路由入口
	2,service作为逻辑入口

*/


//创建房间
func HandlerCreateDesk(userId uint32, a gate.Agent) {
	room := doudizhu.GetFDdzRoom()
	desk := room.CreateDesk(userId)
	if desk == nil {
		log.E("创建房间失败...")
		//todo return error ack
		return
	}

	//自动进入 desk
	err := desk.EnterUser(userId, a)
	if err != nil {
		log.E("用户进入房间失败...")
	}
}

//进入房间的逻辑
func HandlerEnterDesk(userId uint32, key string, deskType int32, a gate.Agent) error {
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
	err := desk.EnterUser(userId, a)
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

//抢地主
func HandlerQiangDiZhu(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("米有找到desk")
	}

	err := desk.QiangDiZhu(userId, 0)
	if err != nil {
		log.E("玩家[%v]抢地主失败,err[%v]", userId, err)
		return Error.NewFailError("玩家抢地主出错")

	}

	return nil

}

//叫地主
func HandlerJiaoDiZhu(userId uint32) error {
	return nil
}

func HandlerShowHandPokers(userId uint32) error {
	return nil
}

func HandlerMenuZhua(userId uint32) error {
	return nil
}

func HandlerSeeCards(userId uint32) error {
	return nil

}

func HandlePull(userId uint32) error {
	return nil
}

func HandlerDissolveDesk(userId uint32) error {
	return nil
}

func HandlerLeaveDesk(userId uint32) error {
	return nil
}

func HandlerMessage() {

}

func HandlerGameRecord() {

}

func HandlerJiaBei(){

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

//pass的协议
func HandlerActPass(userId uint32) error {
	return nil
}
