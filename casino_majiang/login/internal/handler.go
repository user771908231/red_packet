package internal

import (
	"reflect"
	"casino_server/common/log"
	"casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_server/msg/bbprotogo"
	"casino_server/service/userService"
	"casino_server/conf/intCons"
	"casino_majiang/msg/funcsInit"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(mjproto.Game_QuickConn{}, handlerREQQuickConn)
}

func getReleaseTagByVersion(version int32) int32 {
	return 0

}

//处理登陆
func handlerREQQuickConn(args []interface{}) {
	m := args[0].(*mjproto.Game_QuickConn)
	a := args[1].(gate.Agent)
	log.T("游戏登陆的时候发送的请求的协议内容login.handler.HandlerREQQuickConn()[%v]", m)
	//需要返回的结果

	result := newProto.NewGame_AckQuickConn()

	//设置releasTag
	*result.ReleaseTag = getReleaseTagByVersion(m.GetCurVersion())                   ///todo  需要把这个值加入到配置文件读取

	//设置用户的信息
	var resultUser *bbproto.User

	//首先判断是否是微信登陆
	if m.GetWxInfo() != nil && m.GetWxInfo().GetOpenId() != "" {
		//log.T("用户是使用微信登陆")
		//微信登陆,如果是微信新用户,则创建一个user,并且保存
		openId := m.GetWxInfo().GetOpenId()
		resultUser = userService.GetUserByOpenId(openId)
		if resultUser == nil {
			//重新生成一个并保存到数据库
			if m.GetWxInfo().GetHeadUrl() == "" || m.GetWxInfo().GetNickName() == "" || m.GetWxInfo().GetOpenId() == "" {
				//表示参数非法 返回错误
			} else {
				resultUser, _ = userService.NewUserAndSave(m.GetWxInfo().GetOpenId(), m.GetWxInfo().GetNickName(), m.GetWxInfo().GetHeadUrl())
			}
		}
	}

	//处理客户端版本升级
	LatestClientVersion := int32(0) //当前已发布客户端版本, TODO:放到配置文件中
	result.IsUpdate = new(int32)
	*result.IsUpdate = 0
	if m.GetCurVersion() < LatestClientVersion {
		log.T("客户端需要升级, 版本为:%v", m.GetCurVersion())
		*result.IsUpdate = 0 //1=强制升级 0=可选升级
		result.DownloadUrl = new(string)
		*result.DownloadUrl = "http://d.tondeen.com/sjtexas.html" //TODO:放入配置文件中
	}
	result.CurrVersion = new(int32)
	*result.CurrVersion = LatestClientVersion

	//服务器停服维护公告
	result.IsMaintain = new(int32)
	result.MaintainMsg = new(string)
	*result.IsMaintain = 0 //TODO:从配置中读取停服维护公告
	if *result.IsMaintain == 1 {
		*result.MaintainMsg = "服务器正在例行维护中，请于今日5:00后再登录游戏!"
	}


	//如果得到的user ==nil 或者 用密码登陆的时候密码不正确
	if resultUser == nil {
		log.E("没有找到用户,返回登陆失败...")
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(result)
		return
	} else {
		*result.Header.Code = intCons.ACK_RESULT_SUCC                           //返回结果
		*result.UserId = resultUser.GetId()
		*result.NickName = resultUser.GetNickName()
		log.T("快速登录,有userId,没有密码时返回的信息:[%v]", result)
		a.WriteMsg(result)
	}
}
