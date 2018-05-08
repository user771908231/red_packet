package main

import (
	"casino_paosangong/game"
	"casino_paosangong/gate"
	"casino_paosangong/login"
	"github.com/name5566/leaf"
	"casino_paosangong/conf"
	"casino_paosangong/service/paosangong"
	"casino_common/proto/ddproto"
	"casino_common/common/service/rpcService"
	"casino_common/common/service/robotService"
	"casino_paosangong/game/rpc_handler"
	"casino_common/common/service/userGameBillService"
)

func main() {
	//初始化配置
	conf.LoadConfig()
	//先初始化房间列表
	paosangong.InitRoomList()

	//监听牛牛rpc
	rpc_handler.LisenAndServeNiuniuRpc(conf.Server.NiuRpcAddr)
	//初始化大厅rpc
	rpcService.HallPool.Init(conf.Server.HallRpcAddr, 1)

	//初始化机器人
	paosangong.RobotManager = robotService.NewRobotManager(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN)

	//gamebill初始化
	userGameBillService.OnInit(int32(ddproto.CommonEnumGame_GID_PAOYAO), 6)

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
