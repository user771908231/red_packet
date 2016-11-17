package internal

import (
	"reflect"
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
	"casino_login/msg/protoInit"
	"casino_login/service"
	"casino_login/msg/protogo"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&loginproto.Game_QuickConn{}, handlerREQQuickConn)
}


//这里需要一个reload 配置文件的功能

//处理登陆
func handlerREQQuickConn(args []interface{}) {
	m := args[0].(*loginproto.Game_QuickConn)
	a := args[1].(gate.Agent)
	log.T("游戏登陆的时候发送的请求的协议内容login.handler.HandlerREQQuickConn()[%v]", m)

	gameId := m.GetGameId()

	//需要返回的结mahjong_desk.proto果
	result := newProto.NewGame_AckQuickConn()
	*result.ReleaseTag = service.GetReleaseTagByVersion(gameId, m.GetCurrVersion())//发布的版本
	*result.IsMaintain = service.GetIsMaintain(gameId)
	*result.CurrVersion = service.GetLatestClientVersion(gameId) ///处理客户端版本升级
	*result.GameServer.Ip = service.GetGameServerIp(gameId)
	*result.GameServer.Port = service.GetGameServerPort(gameId)
	*result.GameServer.Status = service.GetGameServerStatus(gameId)

	*result.IsUpdate = 0        //默认可选升级

	if m.GetCurrVersion() < result.GetCurrVersion() {
		*result.IsUpdate = service.GetIsUpdate(gameId)//1=强制升级 0=可选升级
		*result.DownloadUrl = service.GetDownloadUrl(gameId) //
	}

	//服务器停服维护公告
	if result.GetIsMaintain() == 1 {
		*result.MaintainMsg = service.GetMaintainMsg(gameId)
	}

	*result.Header.Code = intCons.ACK_RESULT_SUCC                           //返回结果
	log.T("handlerREQQuickConn 协议返回的信息:[%v]", result)
	a.WriteMsg(result)
}

