package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_mj_changsha/service/majiang"
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/common/userService"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.T("地址[%v]建立连接...", a.RemoteAddr())
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.T("地址[%v]断开连接", a.RemoteAddr())

	//通过session 对对应玩家做对应处理...
	agentData := a.UserData() //agentUsr 的数据只在这里有用
	if agentData != nil {
		userData := agentData.(*ddproto.GameSession)
		desk := majiang.MjroomManagerIns.GetMjDeskBySession(userData.GetUserId()) //断开连接时候的处理
		if desk != nil {
			//这里一般不存在desk==nil的情况
			desk.SetOfflineStatus(userData.GetUserId())
		}
		//保存用户的数据，这是为了保证redis-money,redis-user,mgo-user 的数据一致性.
		userService.SyncMgoUserMoney(userData.GetUserId()) //关闭链接的时候，同步玩家的游戏币
	}
}
