package internal

import (
	"reflect"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/service/userService"
	"casino_server/conf/intCons"
	"runtime"
	"strings"
	"strconv"
	"fmt"
	"casino_server/conf"
	"casino_server/utils/numUtils"
	"casino_server/service/noticeServer"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&bbproto.REQQuickConn{}, HandlerREQQuickConn)
	handleMsg(&bbproto.NullMsg{}, handlerNullMsg)
	handleMsg(&bbproto.Game_Notice{}, handlerNotice)
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


///处理联众游戏,登陆的协议
func HandlerREQQuickConn(args []interface{}) {
	m := args[0].(*bbproto.REQQuickConn)
	log.T("联众游戏登陆的时候发送的请求的协议内容login.handler.HandlerREQQuickConn()[%v]", m)
	a := args[1].(gate.Agent)
	//需要返回的结果

	//ip地址信息
	result := &bbproto.ACKQuickConn{}
	result.CoinCnt = new(int64)
	result.UserName = new(string)
	result.UserId = new(uint32)
	result.NickName = new(string)
	result.AckResult = new(int32)

	var resultUser *bbproto.User

	arrs := strings.Split(conf.Server.TCPAddr, ":")
	var ip string = arrs[0]
	var port string = arrs[1]

	ogRoomInfo := &bbproto.OGRoomInfo{}
	ogRoomInfo.RoomIp = &ip
	ogRoomInfo.RoomPort = &port

	thranSjjInfo := &bbproto.ThranJSSInfo{}
	thranSjjInfo.RoomIP = &ip
	thranSjjInfo.RoomPort = &port

	oglist := make([]*bbproto.OGRoomInfo, 1)
	oglist[0] = ogRoomInfo

	tslist := make([]*bbproto.ThranJSSInfo, 1)
	tslist[0] = thranSjjInfo
	result.JssList = tslist
	result.MatchSvrList = oglist

	// 通过userId来判断是登录还是注册,如果userId ==0 ,重新注册一个,如果userId !=0 从数据库查询
	if m.GetUserId() == 0 {
		resultUser, _ = userService.NewUserAndSave()
	} else {
		resultUser = userService.GetUserById(m.GetUserId())
	}


	//如果得到的user ==nil 或者 用密码登陆的时候密码不正确
	if resultUser == nil || (resultUser.GetPwd() != "" && m.GetPwd() != resultUser.GetPwd()) {
		log.E("没有找到用户,返回登陆失败...")
		*result.AckResult = intCons.ACK_RESULT_ERROR
		a.WriteMsg(result)
		return
	} else {
		*result.AckResult = intCons.ACK_RESULT_SUCC                                //返回结果
		*result.CoinCnt = resultUser.GetCoin()                                //用户金币
		*result.UserName, _ = numUtils.Uint2String(resultUser.GetId())        //
		*result.UserId = resultUser.GetId()
		*result.NickName = resultUser.GetNickName()
		log.T("快速登录,有userId,没有密码时返回的信息:[%v]", result)
		a.WriteMsg(result)
	}
}


//空协议
func handlerNullMsg(args []interface{}) {
	log.T("收到一条空消息")
	a := args[1].(gate.Agent)
	a.WriteMsg(&bbproto.NullMsg{})
}


//请求当前的公告
func handlerNotice(args []interface{}) {
	m :=  args[0].(*bbproto.Game_Notice)
	a :=  args[1].(gate.Agent)
	log.T("查询公告type[%v]",m.GetNoticeType())
	tnotice := noticeServer.GetNoticeByType(m.GetNoticeType())
	log.T("查询到了notice[%v]",tnotice)
	a.WriteMsg(tnotice)
}