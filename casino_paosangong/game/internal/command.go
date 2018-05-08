package internal

import (
	"fmt"
	"casino_paosangong/conf"
	"casino_common/utils/numUtils"
	"casino_paosangong/service/paosangong"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("reload", "", commandReload)
	skeleton.RegisterCommand("robotadd", "增加机器人", commandInitRobots)
	skeleton.RegisterCommand("dessolveDesk", "强制解散某个房间，参数为房间号。", commandDessolveDesk)
}

//重新载入配置
func commandReload(args []interface{}) interface{} {
	//载入配置
	conf.LoadConfig()
	return "reload ok."
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func commandInitRobots(args []interface{}) interface{} {
	countS := args[0].(string)
	count := numUtils.String2Int(countS)
	if paosangong.RobotManager == nil {
		return fmt.Sprintf("找不到机器人管理实例..")
	}

	paosangong.RobotManager.NewRobotsAndSave(int32(count),1000, 50000)
	return fmt.Sprintf("增加机器人..%v", args)
}

//解散房间
func commandDessolveDesk(args []interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Sprintf("房号参数不全！")
	}
	passwd := args[0].(string)
	if len(passwd) != 6 {
		return fmt.Sprintf("房号参数异常！")
	}

	room, err := paosangong.Rooms.GetRoomById(0)
	if err != nil {
		return err.Error()
	}

	desk, err := room.GetDeskByPassword(passwd)
	if err != nil {
		return err.Error()
	}

	//解散房间
	err = room.RemoveFriendDesk(desk.GetDeskId())
	if err != nil {
		return fmt.Sprintf("解散失败！错误:%v", err.Error())
	}
	return fmt.Sprintf("解散房间%s成功！", desk.GetDeskNumber())
}
