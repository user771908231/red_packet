package coinPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/huParserIns"
	"casino_common/common/log"
	"casino_common/utils/timeUtils"
	"github.com/name5566/leaf/timer"
	"casino_common/common/Error"
	"github.com/golang/protobuf/proto"
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
)

//朋友桌麻将的desk
type CMJDesk struct {
	*skeleton.SkeletonMJDesk
	RobotEnterTimer *timer.Timer
	CoinLimit       int64 //金币场的金币最低限制
}

//创建一个朋友桌的desk
func NewCMJDesk(config data.SkeletonMJConfig) api.MjDesk {
	//判断创建条件：房卡，
	//desk 骨架
	desk := &CMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}
	desk.HuParser = huParserIns.NewChengDuHuParser()
	return desk
}

//准备
func (d *CMJDesk) Ready(userId uint32) error {
	d.Lock()
	defer d.Unlock()
	d.SkeletonMJDesk.Ready(userId)
	d.initEnterTimer() //房间进入一个人之后开始计划添加机器人
	d.begin()
	return nil
}

//todo
func (d *CMJDesk) initEnterTimer() {
	//房间中有人的时候，才会让机器人进来
	log.T("%v开始准备机器人的enterTimer,d.GetUserCount(%v),d.GetUserCountLimit(%v)", d.DlogDes(), d.GetUserCount(), d.GetMJConfig().PlayerCountLimit)
	if d.GetUserCount() > 0 && d.GetUserCount() < d.GetMJConfig().PlayerCountLimit {
		//5到10秒之后进入一个机器人
		if d.RobotEnterTimer != nil {
			d.RobotEnterTimer.Stop()
		}
		//停止之后，再来重新计时进入机器人
		d.RobotEnterTimer = d.AfterFunc(timeUtils.RandDuration(5, 10), func() {
			log.T("%v现在开始添加机器人.", d.DlogDes())
			d.enterRobot() //进入一个
		})
	}
}

//机器人进入房间
func (d *CMJDesk) enterRobot() {
	log.T("%v 开是添加机器人", d.DlogDes())
	//1,做异常处理
	defer Error.ErrorRecovery("添加机器人")

	//2,获取机器人
	robot := d.Room.GetRoomMgr().GetRobotManger().ExpropriationRobotByCoin(d.CoinLimit)
	if robot == nil {
		log.E("[%v]添加机器人的时候，没有找到合适的机器人...", d.DlogDes())
		return
	}

	//3，加入房间
	err := d.EnterUser(robot.GetId(), nil) //机器人进入房间
	if err != nil {
		//用户加入房间失败...
		log.E("机器人玩家[%v]加入房间失败errMsg[%v]", robot.GetId(), err)
		d.Room.GetRoomMgr().GetRobotManger().ReleaseRobots(robot.GetId())
	}

}

//todo
func (d *CMJDesk) begin() {

}

//金币场的玩家
func (d *CMJDesk) getCMJUser(user api.MjUser) *CMJUser {
	return user.(*CMJUser)
}

//金币场离开房间
func (d *CMJDesk) Leave(userId uint32) error {
	//离开之后设置玩家是托管状态
	user := d.GetUserByUserId(userId)
	if d.GetStatus().IsNotPreparing() {
		//如果游戏没有开始，直接离开
		err := d.rmUser(user)
		if err != nil {
			log.E("玩家离开房间的时候出错:%v", err)
			return nil
		}
	} else {
		//如果游戏已经开始
		//离开不用删除user ,设置为离开之后，设置成托管状态
		user.GetStatus().SetAgentMode(true) //离开的时候设置为托管模式
		user.GetStatus().IsLeave = true
	}

	//回复离开房间的回复
	ack := &ddproto.CommonAckLeaveDesk{
		UserId:     proto.Uint32(userId),
		IsExchange: proto.Bool(false)}
	user.WriteMsg(ack) //回复离开房间的回复
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
			go d.Room.EnterUser(userId, "") //进入房间
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
