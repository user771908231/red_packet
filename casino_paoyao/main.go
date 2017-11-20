package main

import (
	"casino_paoyao/game"
	"casino_paoyao/gate"
	"casino_paoyao/login"
	"github.com/name5566/leaf"
	"casino_paoyao/conf"
	//"casino_paoyao/service/paoyao"
	//"casino_common/proto/ddproto"
	//"casino_common/common/service/rpcService"
	//"casino_common/common/service/robotService"
	//"casino_paoyao/game/rpc_handler"
)

func main() {
	//初始化配置
	conf.LoadConfig()

	//todo 先初始化房间列表
	//paoyao.InitRoomList()

	//todo 监听刨幺rpc
	//rpc_handler.LisenAndServeNiuniuRpc(conf.Server.PaoyaoRpcAddr)
	//todo 初始化大厅rpc
	//rpcService.HallPool.Init(conf.Server.HallRpcAddr, 1)

	//todo 初始化机器人
	//paoyao.RobotManager = robotService.NewRobotManager(ddproto.CommonEnumGame_GID_PAOYAO)

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
