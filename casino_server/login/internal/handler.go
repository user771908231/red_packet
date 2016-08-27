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


///处理登陆协议
func HandlerREQQuickConn(args []interface{}) {
	m := args[0].(*bbproto.REQQuickConn)
	log.T("游戏登陆的时候发送的请求的协议内容login.handler.HandlerREQQuickConn()[%v]", m)
	a := args[1].(gate.Agent)
	//需要返回的结果

	//ip地址信息
	result := newACKQuickConn()
	var resultUser *bbproto.User

	//首先判断是否是微信登陆
	if m.GetWx() != nil && m.GetWx().GetOpenId() != "" {
		//log.T("用户是使用微信登陆")
		//微信登陆,如果是微信新用户,则创建一个user,并且保存
		openId := m.GetWx().GetOpenId()
		resultUser = userService.GetUserByOpenId(openId)
		if resultUser == nil {
			//重新生成一个并保存到数据库
			if m.GetWx().GetHeadUrl()==""||m.GetWx().GetNickName()==""||m.GetWx().GetOpenId()==""{
				//表示参数非法 返回错误
			}else{
				resultUser, _ = userService.NewUserAndSave(m.GetWx().GetOpenId(),m.GetWx().GetNickName(),m.GetWx().GetHeadUrl())
			}
		}
	}

	//处理客户端版本升级
	LatestClientVersion := int32(0) //当前已发布客户端版本, TODO:放到配置文件中
	result.ForceUpdate = new(int32)
	*result.ForceUpdate = 0
	if m.GetCurVersion() < LatestClientVersion {
		log.T("客户端需要升级, 版本为:%v", m.GetCurVersion() )
		*result.ForceUpdate = 0 //1=强制升级 0=可选升级
		result.DownloadUrl = new(string)
		*result.DownloadUrl = "http://d.tondeen.com/sjtexas.html" //TODO:放入配置文件中
	}
	result.CurVersion = new(int32)
	*result.CurVersion = LatestClientVersion

	//服务器停服维护公告
	result.IsMaintain = new(int32)
	result.MaintainMsg = new(string)
	*result.IsMaintain = 0 //TODO:从配置中读取停服维护公告
	if *result.IsMaintain == 1 {
		*result.MaintainMsg = "服务器正在例行维护中，请于今日5:00后再登录游戏!"
	}

	
	//如果得到的user ==nil 或者 用密码登陆的时候密码不正确
	if resultUser == nil || (resultUser.GetPwd() != "" && m.GetPwd() != resultUser.GetPwd()) {
		log.E("没有找到用户,返回登陆失败...")
		*result.AckResult = intCons.ACK_RESULT_ERROR
		a.WriteMsg(result)
		return
	} else {
		*result.AckResult = intCons.ACK_RESULT_SUCC                           //返回结果
		*result.CoinCnt = resultUser.GetDiamond()                             //用户金币--> 这里返回用户砖石的数量
		*result.UserName, _ = numUtils.Uint2String(resultUser.GetId())        //
		*result.UserId = resultUser.GetId()
		*result.NickName = resultUser.GetNickName()
		log.T("快速登录,有userId,没有密码时返回的信息:[%v]", result)
		a.WriteMsg(result)
	}
}


//空协议
func handlerNullMsg(args []interface{}) {
	//log.T("收到一条空消息")
	a := args[1].(gate.Agent)
	a.WriteMsg(&bbproto.NullMsg{})
}


//请求当前的公告
func handlerNotice(args []interface{}) {
	m :=  args[0].(*bbproto.Game_Notice)
	a :=  args[1].(gate.Agent)
	log.T("开始查询公告的信息m[%v]",m)
	tnotice := noticeServer.GetNoticeByType(m.GetNoticeType())
	a.WriteMsg(tnotice)
}



//返回登陆信息
func newACKQuickConn() * bbproto.ACKQuickConn{
	result := &bbproto.ACKQuickConn{}

	result.CoinCnt = new(int64)
	result.UserName = new(string)
	result.UserId = new(uint32)
	result.NickName = new(string)
	result.AckResult = new(int32)
	result.ReleaseTag = new(int32)
	//默认发布版本是1
	*result.ReleaseTag = 1			///todo  这里需要修改

	arrs := strings.Split(conf.Server.TCPAddr, ":")
	var ip string = arrs[0]
	var port string = arrs[1]

	ogRoomInfo := &bbproto.DDRoomInfo{}
	ogRoomInfo.RoomIp = &ip
	ogRoomInfo.RoomPort = &port

	thranSjjInfo := &bbproto.SHENJINGJSSInfo{}
	thranSjjInfo.RoomIP = &ip
	thranSjjInfo.RoomPort = &port

	oglist := make([]*bbproto.DDRoomInfo, 1)
	oglist[0] = ogRoomInfo

	tslist := make([]*bbproto.SHENJINGJSSInfo, 1)
	tslist[0] = thranSjjInfo
	result.JssList = tslist
	result.MatchSvrList = oglist


	//返回客户端服务器IP
	serverInfo := &bbproto.ServerInfo{}
	serverInfo.Ip = &ip

	portStr := int32(numUtils.String2Int(port))
	serverInfo.Port = &portStr
	svrlist := make([]*bbproto.ServerInfo, 1)
	svrlist[0] = serverInfo
	result.ServerList = svrlist

	return result
}
