package paoyao

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

//加入机器人
func (desk *Desk) AutoJoinRobot() {
	go func() {
		Error.ErrorRecovery("AutoJoinRobot()")
		rand_join := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

		switch desk.GetRobotNumber() {
		case 0:
			rand_join = 100
		case 1:
			if rand_join < 30 {
				rand_join = 100
			}
		case 2:
			if rand_join < 10 {
				rand_join = 100
			}
		}

		if rand_join != 100 {
			return
		}

		//加入机器人
		desk.AddRobot()
	}()
}

//加入机器人
func (desk *Desk) AddRobot() {
	go func() {
		Error.ErrorRecovery("AddRobot()")
		//加入机器人
		new_robot := RobotManager.ExpropriationRobotByRange(7000, 100000)
		if new_robot == nil {
			log.E("没有机器人了。房间%v加入机器人失败。", desk.GetDeskId())
			return
		}
		new_user,err := desk.AddUser(new_robot.GetId(), nil)
		if new_user == nil {
			log.E("房间%d加入机器人失败,错误：%v", desk.GetDeskId(), err)
			return
		}
		//设为机器人
		*new_user.IsRobot = true

		//自动准备
		if desk.GetStatus() == ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY {
			new_user.DoReadyCoin()
		}
	}()
}

//抢庄
func (user *User) DoRobotQiangzhuang() {
	go func() {
		Error.ErrorRecovery("DoRobotQiangzhuang()")
		<-time.After(timeUtils.RandDuration(3, 7))
		//todo
	}()
}

//加倍
func (user *User) DoRobotJiabei() {
	go func() {
		Error.ErrorRecovery("DoRobotJiabei()")
		<-time.After(timeUtils.RandDuration(3, 7))
		//todo

	}()
}

//继续游戏
func (user *User) DoRobotReady() {
	go func() {
		Error.ErrorRecovery("DoRobotReady()")
		var sleep_time int32 = 3
		//如果有动画，则11秒后才准备
		if user.Desk.DeskOption.GetHasAnimation() {
			sleep_time = 11
		}
		<-time.After(timeUtils.RandDuration(sleep_time, sleep_time + 5))
		rand_ex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

		switch {
		case rand_ex<10:
			//平均美10局换一个机器人
			user.DoLeaveDesk()
			<-time.After(2*time.Second)
			user.Desk.AddRobot()
		default:
			user.DoReadyCoin()
		}
	}()
}
