package paosangong

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
			if u.GetIsRobot() && u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY && u.GetIsReady() == false {
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
		if desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY {
			<-time.After(timeUtils.RandDuration(2, 6))
			new_user.DoReadyCoin()
		}
	}()
}

//抢庄
func (user *User) DoRobotQiangzhuang() {
	go func() {
		Error.ErrorRecovery("DoRobotQiangzhuang()")
		<-time.After(timeUtils.RandDuration(3, 7))
		switch {
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_6:
			user.DoQiangzhuang(4)
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_3:
			user.DoQiangzhuang(2)
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_1:
			user.DoQiangzhuang(1)
		default:
			user.DoQiangzhuang(-1)
		}
	}()
}

//加倍
func (user *User) DoRobotJiabei() {
	go func() {
		Error.ErrorRecovery("DoRobotJiabei()")
		<-time.After(timeUtils.RandDuration(8, 15))
		//加倍倍数
		const double_score int64 = 5
		switch {
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_6:
			//如果比庄家大
			if user.IsBigThanBanker() {
				user.DoJiabei(5 * double_score)
			}else {
				user.DoJiabei(3 * double_score)
			}
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_4:
			//如果比庄家大
			if user.IsBigThanBanker() {
				user.DoJiabei(4 * double_score)
			}else {
				user.DoJiabei(2 * double_score)
			}
		case user.Pokers.GetType() >= ddproto.NiuniuEnum_PokerType_NIU_1:
			//如果比庄家大
			if user.IsBigThanBanker() {
				user.DoJiabei(3 * double_score)
			}else {
				user.DoJiabei(1 * double_score)
			}
		default:
			user.DoJiabei(1 * double_score)
		}
	}()
}

//继续游戏
func (user *User) DoRobotReady(isEnter bool) {
	go func() {
		Error.ErrorRecovery("DoRobotReady()")
		var sleep_time int32 = 2
		//继续游戏 如果有动画，则11秒后才准备
		if !isEnter && user.Desk.DeskOption.GetHasAnimation() {
			sleep_time = 13
		}
		<-time.After(timeUtils.RandDuration(sleep_time, sleep_time + 6))
		rand_ex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

		switch {
		case rand_ex<10:
			//平均美10局换一个机器人
			user.DoLeaveDesk()
			<-time.After(2*time.Second)
			user.Desk.AddRobot()
		case rand_ex<20:
			//如果真人大于1，直接退出
			if len(user.Desk.Users) - int(user.Desk.GetRobotNumber()) > 1 {
				user.DoLeaveDesk()
			}
		default:
			//有真人才准备
			if len(user.Desk.Users) - int(user.Desk.GetRobotNumber()) > 0 {
				user.DoReadyCoin()
			}
		}
	}()
}
