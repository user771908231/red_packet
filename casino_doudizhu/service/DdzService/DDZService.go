package DdzService

import (
	"casino_server/common/log"
	"casino_doudizhu/service/doudizhu"
	"errors"
	"casino_server/common/Error"
	"github.com/name5566/leaf/gate"
	"casino_doudizhu/msg/protogo"
	"casino_doudizhu/msg/funcsInit"
	"casino_server/conf/intCons"
)

/*
	1,gate 作为路由入口
	2,service作为逻辑入口

*/


//创建房间
func HandlerCreateDesk(userId uint32, roominfo *ddzproto.RoomTypeInfo, a gate.Agent) {
	log.T("玩家[%v]创建房间roomInfo[%v]", userId, roominfo)
	room := doudizhu.GetFDdzRoom()
	desk := room.CreateDesk(userId, roominfo)
	if desk == nil {
		log.E("创建房间失败...")
		ret := newProto.NewGame_AckCreateRoom()
		*ret.Header.Code = intCons.ACK_RESULT_ERROR
		*ret.Header.Error = "创建房间失败"
		a.WriteMsg(ret)
		return
	}

	//自动进入 desk
	go HandlerEnterDesk(userId, desk.GetKey(), desk.GetRoomType(), a)
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

	//进入房间,失败返回失败的信息，成功放回成功的信息
	err, isReconnect := desk.EnterUser(userId, a)
	if err != nil {
		log.E("玩家[%v]进入desk[%v]失败err[%v]", userId, desk.GetDeskId(), err)
		return Error.NewError(-1, "玩家进入desk失败.")
	} else {
		desk.SendGameDeskInfo(userId, isReconnect)        //发送deskGameInfo
	}

	return nil
}

//准备游戏
func HandlerFDdzReady(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		log.E("玩家[%v]准备失败，没有找到对应的desk..", userId)
		return Error.NewError(-1, "没有找到玩家所在的desk")
	}

	err := desk.Ready(userId)
	if err != nil {
		log.E("玩家[%v]准备游戏的时候失败...")
		return Error.NewError(-1, "玩家准备游戏的时候失败...")
	}
	return nil
}

//抢地主
func HandlerQiangDiZhu(userId uint32, qiang bool) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}

	//开始抢地主
	var err error
	if desk.IsHuanLeDoudDiZhu() {
		err = desk.HLQiangDiZhu(userId,qiang)
	}

	//抢地主失败
	if err != nil {
		log.E("玩家[%v]抢地主失败,err[%v]", userId, err)
		return Error.NewFailError("玩家抢地主出错")
	}
	return nil
}

//叫地主
func HandlerJiaoDiZhu(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}

	//开始叫地主
	var err error
	if desk.IsHuanLeDoudDiZhu() {
		err = desk.HLJiaoDiZhu(userId)
	} else if desk.IsSiChuanDouDiZhu() {
		err = desk.SCJiaoDiZhu(userId)
	}

	//判断叫地主是否成功
	if err != nil {
		log.E("叫地主失败")
		return err
	}

	return nil
}

//不叫地主
func HandlerBuJiaoDiZhu(userId uint32) error {
	//不叫地主
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}

	var err error

	//开始不叫地主
	if desk.IsHuanLeDoudDiZhu() {
		err = desk.HLBuJiaoDiZhu(userId)
	}

	//判断叫地主是否成功
	if err != nil {
		log.E("玩家[%v]不叫地主的时候失败", userId)
		return err
	}
	return nil
}

//明牌的逻辑
/**
	1，明牌只在欢乐斗地主中出现
	2，明牌会把低分翻倍
 */
func HandlerShowHandPokers(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}

	err := desk.ShowHandPokers(userId)
	if err != nil {
		return err
	}
	return nil
}

//闷抓
/**
	1， 一般出现在四川斗地主中

 */
func HandlerMenuZhua(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	return nil
}

func HandlerSeeCards(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	return nil

}

func HandlePull(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	return nil
}

func HandlerDissolveDesk(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	return nil
}

func HandlerLeaveDesk(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	return nil
}

func HandlerMessage(m *ddzproto.DdzMessage) {
	log.T("请求发送信息[%v]", m)
	userId := m.GetHeader().GetUserId()
	desk := doudizhu.GetDdzDeskBySession(userId)

	if desk == nil {
		log.E("玩家[%v]聊天的时候没有找到desk", userId)
		return
	}

	result := newProto.NewGameMessage()
	*result.UserId = m.GetHeader().GetUserId()
	*result.Id = m.GetId()
	*result.Msg = m.GetMsg()
	*result.MsgType = m.GetMsgType()
	desk.BroadCastProtoExclusive(result, result.GetUserId())
}

//查询战绩
func HandlerGameRecord(userId uint32, a gate.Agent) error {

	//
	//返回战绩
	return nil
}

func HandlerJiaBei(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}
	//目前只有欢乐斗地主才有加倍的逻辑
	var err error
	if desk.IsHuanLeDoudDiZhu() {
		err = desk.ActJiaBei(userId)
	}
	//判断结果
	if err != nil {
		log.E("玩家[%v]加倍失败..err[%v]", userId, err)
		return err
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

//pass的协议
func HandlerActPass(userId uint32) error {
	desk := doudizhu.GetDdzDeskBySession(userId)
	if desk == nil {
		return Error.NewFailError("没有找到desk")
	}

	err := desk.ActPass(userId)
	if err != nil {
		log.E("玩家[%v]过牌的时候失败", userId)
		return err
	}
	return nil
}

