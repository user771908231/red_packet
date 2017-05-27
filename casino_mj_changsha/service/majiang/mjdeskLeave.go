package majiang

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
	"github.com/name5566/leaf/gate"
	"errors"
)

//离开房间
/**
	判断当前desk的状态
	1,如果是游戏已经开始那么不能离开
 */
func (d *MjDesk) Leave(userId uint32) error {

	//1
	if d.IsGaming() {
		return ERR_LEAVE_RUNNING //离开房间失败
	}

	//2,准备阶段的时候可以离开
	if d.IsNotPreparing() {
		return ERR_LEAVE_RUNNING //只有在准备的时候才可以离开
	}

	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("玩家[%v]离开的时候失败，没有找到对应的user", user.GetUserId())
		return ERR_LEAVE_ERROR
	}

	err := d.rmUser(user) // 用户自己选择离开朋友桌房间时
	if err != nil {
		return err
	}

	//发送离开的广播
	ack := new(ddproto.CommonAckLeaveDesk)
	ack.UserId = proto.Uint32(user.GetUserId())
	ack.IsExchange = proto.Bool(false)
	d.BroadCastProto(ack)
	user.WriteMsg(ack) //自己没有收到回复

	if d.canDissolve() {
		room := MjroomManagerIns.GetFMJRoom()
		if room == nil {
			log.E("解散房间的时候出错...")
			return nil
		}

		dissolveErr := room.DissolveDesk(d, false) //离开的时候，如果没有人了，开始解散房间
		if dissolveErr != nil {
			log.E("解散房间的时候失败%v", dissolveErr)
			return nil
		}

	}
	return nil
}

//金币场离开房间
func (d *MjDesk) LeaveCoin(userId uint32) error {
	log.T("锁日志: %v LeaveCoin(%v)v开始离开金币场的房间-等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v LeaveCoin(%v)v开始离开金币场的房间-释放所", d.DlogDes(), userId)
	}()

	//离开之后设置玩家是托管状态
	user := d.GetUserByUserId(userId)
	if d.IsPreparing() {
		log.T("%v ，玩家%v开始离开金币场的房间.游戏处于准备中，直接删除User", d.DlogDes(), userId)
		//如果游戏没有开始，直接离开
		err := d.rmUser(user) //用户离开金币场的时候
		if err != nil {
			log.E("玩家离开房间的时候出错:%v", err)
			return nil
		}
	} else {
		log.T("%d ，玩家%v开始离开金币场的房间.游戏正在进行中，设置user为托管模式", d.DlogDes(), userId)
		//如果游戏已经开始
		//离开不用删除user ,设置为离开之后，设置成托管状态
		user.setAgentMode(true) //离开的时候设置为托管模式
		user.IsLeave = proto.Bool(true)
	}
	return nil
}

//可以更换房间的条件
/**
	1,根据玩家是否已经胡牌来判断
		胡牌之后:可以马上换房间
		胡牌之前:
		游戏开始之前可以换
		依稀结束之后可以换
 */

func (d *MjDesk) canExchange(userId uint32) bool {
	//如果玩家已经胡牌来，可以直接换房间
	user := d.GetUserByUserId(userId)
	if user.IsHu() {
		return true
	}

	//如果没有胡牌 准备的阶段可以换房
	if d.IsPreparing() {
		return true
	}

	return false

}

func (d *MjDesk) ExchangeRoom(userId uint32, a gate.Agent) error {
	//更换房间只有在没有开始游戏的时候
	if !d.canExchange(userId) {
		return errors.New("离开房间失败...")
	}

	//先离开
	err := d.LeaveCoin(userId) //金币场  更换房间...
	if err != nil {
		return err
	}
	//回复离开房间的协议
	ack := &ddproto.CommonAckLeaveDesk{
		UserId:     proto.Uint32(userId),
		IsExchange: proto.Bool(false)}
	//user
	d.BroadCastProtoExclusive(ack, userId) //给其他人发送此人离开的协议

	go HandlerGame_EnterDesk(userId, "", d.GetRoomType(), d.GetRoomLevel(), ENTERTYPE_NORMAL, a)
	return nil
}

//删除一个user
func (d *MjDesk) rmUser(user *MjUser) error {
	for i, u := range d.Users {
		if u != nil && u.GetUserId() == user.GetUserId() {
			//更新session
			d.Users[i] = nil
			d.RmUserSession(u) //删除一个玩家

			if user.GetIsRobot() {
				MjroomManagerIns.RobotManger.ReleaseRobots(user.GetUserId())
			}
		}
	}
	return nil
}
