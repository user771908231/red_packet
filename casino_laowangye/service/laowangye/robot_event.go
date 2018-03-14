package laowangye

import (
	"math/rand"
	"time"
	"casino_common/common/service/robotService"
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/common/Error"
	"casino_common/utils/timeUtils"
)

//机器人管理器
var RobotManager *robotService.RobotsManager

//当前牌桌机器人数量
func (desk *Desk) GetRobotNumber() (num int32) {
	for _,u := range desk.Users {
		if u != nil && u.GetIsRobot() {
			num++
		}
	}
	return
}

//进房时触发，加入机器人
func (desk *Desk) AutoJoinRobot() {
	go func() {
		Error.ErrorRecovery("AutoJoinRobot()")
		rand_join := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

		addNum := 0
		switch desk.GetRobotNumber() {
		case 0:
			addNum = 3
		case 1:
			addNum = 2
		case 2:
			addNum = 1
		case 3:
			if rand_join < 20 {
				addNum = 1
			}
		}

		//加入机器人
		for i:=0;i<addNum;i++ {
			desk.AddRobot()
		}

		for _,u := range desk.Users {
			if u.GetIsRobot() && u.Desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY && u.GetIsReady() == false {
				u.DoRobotReady(true)
			}
		}
	}()
}

//加入机器人
func (desk *Desk) AddRobot() {
	go func() {
		Error.ErrorRecovery("AddRobot()")
		//加入机器人
		new_robot := RobotManager.ExpropriationRobotByRange(7000, 1000000)
		if new_robot == nil {
			log.E("没有机器人了。房间%v加入机器人失败。", desk.GetDeskId())
			return
		}
		new_user,err := desk.AddUser(new_robot.GetId(), nil)
		if new_user == nil {
			log.E("房间%d加入机器人失败,错误：%v", desk.GetDeskId(), err)
			RobotManager.ReleaseRobots(new_robot.GetId())
			return
		}
		//设为机器人
		*new_user.IsRobot = true

		//自动准备
		if desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
			<-time.After(timeUtils.RandDuration(2, 6))
			new_user.DoReadyCoin()
		}
	}()
}

//抢庄
func (user *User) DoRobotQiangzhuang() {

}

//加倍
func (user *User) DoRobotJiabei() {

}

//继续游戏
func (user *User) DoRobotReady(isEnter bool) {

}
