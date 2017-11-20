package internal

import (
	"fmt"
	"casino_paoyao/conf"
	"casino_common/utils/numUtils"
	"casino_paoyao/service/paoyao"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("reload", "", commandReload)
	skeleton.RegisterCommand("robotadd", "增加机器人", commandInitRobots)
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
	if paoyao.RobotManager == nil {
		return fmt.Sprintf("找不到机器人管理实例..")
	}

	paoyao.RobotManager.NewRobotsAndSave(int32(count),1000, 50000)
	return fmt.Sprintf("增加机器人..%v", args)
}
