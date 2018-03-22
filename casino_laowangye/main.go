package main

import (
	"casino_laowangye/game"
	"casino_laowangye/gate"
	"casino_laowangye/login"
	"github.com/name5566/leaf"
	"casino_laowangye/conf"
	"casino_laowangye/service/laowangye"
	//"casino_common/proto/ddproto"
	"casino_common/common/service/rpcService"
	//"casino_common/common/service/robotService"
	//"casino_laowangye/game/rpc_handler"
	//"casino_common/common/service/userGameBillService"
)

func main() {
	//初始化配置
	conf.LoadConfig()
	//先初始化房间列表
	laowangye.InitRoomList()

	//监听老王爷rpc
	//rpc_handler.LisenAndServeNiuniuRpc(conf.Server.LwyRpcAddr)
	//初始化大厅rpc
	rpcService.HallPool.Init(conf.Server.HallRpcAddr, 1)

	//初始化机器人
	//laowangye.RobotManager = robotService.NewRobotManager(ddproto.CommonEnumGame_GID_LAOWANGYE)

	//gamebill初始化
	//userGameBillService.OnInit(int32(ddproto.CommonEnumGame_GID_LAOWANGYE), 6)

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
