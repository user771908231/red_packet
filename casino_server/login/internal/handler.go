package internal

import (
	"reflect"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/service/userService"
	"casino_server/conf/StrCons"
	"casino_server/conf/intCons"
	"runtime"
	"strings"
	"strconv"
	"fmt"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/gamedata"
	"casino_server/conf"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&bbproto.Reg{},handleProtHello)
	handleMsg(&bbproto.ReqAuthUser{},handleReqAuthUser)
	handleMsg(&bbproto.REQQuickConn{},HandlerREQQuickConn)
	handleMsg(&bbproto.NullMsg{},handlerNullMsg)
}

/**
	处理注册消息的方法
	此方法可能暂时没有使用,而使用handleReqAuthUser
 */
func handleProtHello(args []interface{}){
	log.T("进入login.handler.handleProtHello()")
	a := args[1].(gate.Agent)
	var data bbproto.Reg
	var n string = "hi leaf"
	data.Name = &n
	a.WriteMsg(&data)
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}


/*

登陆注册的流程:
1,获取到请求信息
2,如果快速登录,则随机分配一个userId,并且自动注册
3,如果用户已经注册了,则走登录流程
 */
func handleReqAuthUser(args []interface{}){
	log.Debug("进入login.handler.handleReqAuthUser()")
	m := args[0].(*bbproto.ReqAuthUser)
	a := args[1].(gate.Agent)

	//判断是快速登录还是普通登录
	var resUser *bbproto.User
	var e   error
	loginWay := userService.CheckUserId(m.GetHeader().GetUserId())
	switch loginWay {
	case intCons.LOGIN_WAY_QUICK:
		log.T("快速登录模式")
		resUser,e = userService.QuickLogin(m)
	case intCons.LOGIN_WAY_LOGIN:
		log.T("普通登录模式")
		resUser,e = userService.Login(m)
	default:
		log.T("没有找到合适的登录方式")
	}

	//判断返回的信息,并且返回信息
	resReqUser := &bbproto.ReqAuthUser{}
	if e != nil {
		log.E(e.Error())
		resReqUser.Header = protoUtils.GetErrorHeaderWithMsg(&StrCons.STR_POINT_ERR_LOGIN_FAIL)
	}else{
		resReqUser.Header = protoUtils.GetSuccHeaderwithMsgUserid(resUser.Id,&StrCons.STR_POINT_ERR_LOGIN_SUCC)
		//增加用户锁
		userService.UserLockPools.AddUserLockByUserId(resReqUser.GetHeader().GetUserId())
	}

	//登录后在服务器需要做的操作
	agentUser :=gamedata.NewAgentUser(resReqUser.GetHeader().GetUserId())
	a.SetUserData(agentUser)

	//把数据返回给客户端
	a.WriteMsg(resReqUser)
}


func HandlerREQQuickConn(args []interface{}){
	log.Debug("进入login.handler.HandlerREQQuickConn()")
	m := args[0].(*bbproto.REQQuickConn)
	log.T("m:【%v】",m)
	a := args[1].(gate.Agent)

	arrs := strings.Split(conf.Server.TCPAddr,":")
	var ip string = arrs[0]
	var port string = arrs[1]

	ogRoomInfo := &bbproto.OGRoomInfo{}
	ogRoomInfo.RoomIp = &ip
	ogRoomInfo.RoomPort =&port

	thranSjjInfo := &bbproto.ThranJSSInfo{}
	thranSjjInfo.RoomIP = &ip
	thranSjjInfo.RoomPort = &port

	oglist := make([]*bbproto.OGRoomInfo,1)
	oglist[0] = ogRoomInfo

	tslist := make([]*bbproto.ThranJSSInfo,1)
	tslist[0] = thranSjjInfo

	result := &bbproto.ACKQuickConn{}
	result.JssList = tslist
	result.MatchSvrList = oglist

	var coin int64 = 999888
	result.CoinCnt = &coin
	log.T("回复的消息:[%v]",result)
	a.WriteMsg(result)
}


func handlerNullMsg(args []interface{}){
	log.T("收到一条空消息")
	//a := args[1].(gate.Agent)
	//retData := &bbproto.MatchList_SendMatchlistSvrInfo{}
	//sinfo := bbproto.MatchList_SvrInfo{}
	//var ip string = "192.168.199.120"
	//var port string = "3563"
	//sinfo.IP =
	//retData.SvrInfo.
	//
	//a.WriteMsg(retData)

}