package internal

import (
	"reflect"
	"casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_majiang/msg/funcsInit"
	"casino_common/common/log"
	"casino_common/common/consts"
	"casino_common/common/userService"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&mjproto.Game_QuickConn{}, handlerREQQuickConn)
	handler(&mjproto.Game_Login{}, handlerGame_Login)
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
	*result.ReleaseTag = getReleaseTagByVersion(m.GetCurrVersion())                   ///todo  需要把这个值加入到配置文件读取

	//处理客户端版本升级
	LatestClientVersion := int32(0) //当前已发布客户端版本, TODO:放到配置文件中
	result.IsUpdate = new(int32)
	*result.IsUpdate = 0
	if m.GetCurrVersion() < LatestClientVersion {
		log.T("客户端需要升级, 版本为:%v", m.GetCurrVersion())
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
	*result.Header.Code = consts.ACK_RESULT_SUCC                           //返回结果
	log.T("handlerREQQuickConn 协议返回的信息:[%v]", result)
	a.WriteMsg(result)

}


/**
	登陆的协议...
 */
func handlerGame_Login(args []interface{}) {
	m := args[0].(*mjproto.Game_Login)
	a := args[1].(gate.Agent)

	log.T("请求handlerGame_Login  m[%v]", m)
	weixin := m.GetWxInfo()

	//不是初次登录
	if weixin == nil {
		//判断uerId
		userId := m.GetHeader().GetUserId()
		user := userService.GetUserById(userId)
		if user == nil {
			//登陆失败
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = consts.ACK_RESULT_ERROR
			a.WriteMsg(ack)
		} else {
			//返回登陆成功的结果
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = consts.ACK_RESULT_SUCC
			*ack.UserId = user.GetId()
			*ack.NickName = user.GetNickName()
			*ack.Chip = user.GetDiamond()
			a.WriteMsg(ack)
		}
		return
	}

	//1,首先通过weixinInfo 在数据库中查找 用户是否存在，如果用户存在，则表示，登陆成功
	user := userService.GetUserByOpenId(weixin.GetOpenId())
	if user == nil {
		//表示数据库中不存在次用户，新增加一个人后返回
		if weixin.GetOpenId() == "" || weixin.GetHeadUrl() == "" || weixin.GetNickName() == "" {
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = consts.ACK_RESULT_ERROR
			a.WriteMsg(ack)
			return
		}

		//如果数据库中不存在用户，那么重新生成一个user
		user, _ = userService.NewUserAndSave(weixin.GetOpenId(), weixin.GetNickName(), weixin.GetHeadUrl(), weixin.GetSex(), weixin.GetCity())
		if user == nil {
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = consts.ACK_RESULT_ERROR
			a.WriteMsg(ack)
			return
		}
	}

	//返回登陆成功的结果
	ack := newProto.NewGame_AckLogin()
	*ack.Header.Code = consts.ACK_RESULT_SUCC
	*ack.UserId = user.GetId()
	*ack.NickName = user.GetNickName()
	*ack.Chip = user.GetDiamond()

	a.WriteMsg(ack)
	return

}

