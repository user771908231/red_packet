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
)

//朋友桌麻将的desk
type CMJDesk struct {
	*skeleton.SkeletonMJDesk
	RobotEnterTimer *timer.Timer
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

//离开房间
func (d *CMJDesk) Leave(userId uint32) error {
	return nil
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
	robot := MjroomManagerIns.RobotManger.ExpropriationRobotByCoin(d.GetCoinLimit())
	if robot == nil {
		log.E("[%v]添加机器人的时候，没有找到合适的机器人...", d.DlogDes())
		return
	}

	//3，加入房间
	err := d.enterUser(robot.GetId(), nil) //机器人进入房间
	if err != nil {
		//用户加入房间失败...
		log.E("机器人玩家[%v]加入房间失败errMsg[%v]", robot.GetId(), err)
		MjroomManagerIns.RobotManger.ReleaseRobots(robot.GetId())
	}

}

//todo
func (d *CMJDesk) begin() {

}

//金币场的玩家
func (d *CMJDesk) getCMJUser(user api.MjUser) *CMJUser {
	return user.(*CMJUser)
}
