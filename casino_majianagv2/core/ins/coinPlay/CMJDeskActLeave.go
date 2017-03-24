package coinPlay

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
	"github.com/name5566/leaf/gate"
)

//金币场离开房间
func (d *CMJDesk) Leave(userId uint32) error {
	log.T("%v ，玩家%v开始离开金币场的房间", d.DlogDes(), userId)
	//离开之后设置玩家是托管状态
	user := d.GetUserByUserId(userId)
	if d.GetStatus().IsPreparing() {
		log.T("%v ，玩家%v开始离开金币场的房间.游戏处于准备中，直接删除User", d.DlogDes(), userId)
		//如果游戏没有开始，直接离开
		err := d.rmUser(user)
		if err != nil {
			log.E("玩家离开房间的时候出错:%v", err)
			return nil
		}
	} else {
		log.T("%v ，玩家%v开始离开金币场的房间.游戏正在进行中，设置user为托管模式", d.DlogDes(), userId)
		//如果游戏已经开始
		//离开不用删除user ,设置为离开之后，设置成托管状态
		user.GetStatus().SetAgentMode(true) //离开的时候设置为托管模式
		user.GetStatus().IsLeave = true
	}
	return nil
}

//只有金币场才有离开房间的逻辑
func (d *CMJDesk) ExchangeRoom(userId uint32, a gate.Agent) error {
	//更换房间只有在没有开始游戏的时候

	if d.canExchange(userId) {
		//先离开
		err := d.Leave(userId)
		if err != nil {
			//打印离开失败的日志
			log.E("")
		} else {
			ack := new(ddproto.CommonAckLeaveDesk)
			ack.UserId = proto.Uint32(userId)
			ack.IsExchange = proto.Bool(true)
			//d.BroadCastProto(ack)
			a.WriteMsg(ack)
			//进入房间
			go d.Room.EnterUser(userId, "", a) //进入房间
		}
	}
	return nil
}

func (d *CMJDesk) canExchange(userId uint32) bool {
	//如果玩家已经胡牌来，可以直接换房间
	user := d.GetUserByUserId(userId)
	if user.GetStatus().IsHu() {
		return true
	}

	//如果没有胡牌 准备的阶段可以换房
	if d.GetStatus().IsPreparing() {
		return true
	}

	return false

}
