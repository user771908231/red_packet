package internal

import (
	"fmt"
	"casino_common/utils/numUtils"
	"casino_mj_changsha/service/majiang"
	"github.com/golang/protobuf/proto"
	"casino_common/common/service/taskService/task"
	"casino_mj_changsha/conf"
	"casino_common/common/service/pushService"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("robotadd", "增加机器人 第一个参数是 机器人人数", commandInitRobots)
	skeleton.RegisterCommand("exmjp", "换牌,level(c,f),deskid,p1,p2", commandExhangeMjPai)
	skeleton.RegisterCommand("setpaicursor", "setpaicursor 参数level(c,f), index", commandSetpaicursor)
	skeleton.RegisterCommand("reload", "reload config.", ReloadConfig)
}

//重新载入配置
func ReloadConfig(args []interface{}) interface{} {
	//载入json配置
	conf.LoadJsonConfig()
	//载入任务配置
	task.InitTask()
	//载入push service初始化配置
	pushService.PoolInit(conf.Server.HallTcpAddr)
	return "reload ok."
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func commandInitRobots(args []interface{}) interface{} {
	countS := args[0].(string)
	count := numUtils.String2Int(countS)
	m := majiang.MjroomManagerIns.RobotManger
	for ; count > 0; count -- {
		m.NewRobotAndSave()
	}
	return fmt.Sprintf("增加机器人..%v", args)
}

//换牌器
func commandExhangeMjPai(args []interface{}) interface{} {
	deskId := int32(numUtils.String2Int(args[2].(string)))
	p1 := numUtils.String2Int(args[3].(string))
	p2 := numUtils.String2Int(args[4].(string))
	var desk *majiang.MjDesk
	//金币场
	desk = majiang.MjroomManagerIns.GetFMJRoom().GetDeskByDeskId(deskId)
	if desk != nil {
		desk.ExchangeMJPai(p1, p2)
	}
	return "换牌完成"
}

func commandSetpaicursor(args []interface{}) interface{} {
	deskId := int32(numUtils.String2Int(args[2].(string)))
	index := numUtils.String2Int(args[3].(string))
	var desk *majiang.MjDesk
	//金币场
	desk = majiang.MjroomManagerIns.GetFMJRoom().GetDeskByDeskId(deskId)
	if desk != nil {
		desk.MJPaiCursor = proto.Int32(int32(index))
	}
	return "坐标设置完成"
}
