package internal

import (
	"reflect"
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
	"casino_server/service/userService"
	"casino_doudizhu/msg/protogo"
	"casino_doudizhu/msg/funcsInit"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&ddzproto.DdzLogin{}, handlerGame_Login)
}

/**
	登陆的协议...
 */
func handlerGame_Login(args []interface{}) {
	m := args[0].(*ddzproto.DdzLogin)
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
			*ack.Header.Code = intCons.ACK_RESULT_ERROR
			a.WriteMsg(ack)
		} else {
			//返回登陆成功的结果
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = intCons.ACK_RESULT_SUCC
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
			*ack.Header.Code = intCons.ACK_RESULT_ERROR
			a.WriteMsg(ack)
			return
		}

		//如果数据库中不存在用户，那么重新生成一个user
		user, _ = userService.NewUserAndSave(weixin.GetOpenId(), weixin.GetNickName(), weixin.GetHeadUrl(), weixin.GetSex(), weixin.GetCity())
		if user == nil {
			ack := newProto.NewGame_AckLogin()
			*ack.Header.Code = intCons.ACK_RESULT_ERROR
			a.WriteMsg(ack)
			return
		}
	}

	//返回登陆成功的结果
	ack := newProto.NewGame_AckLogin()
	*ack.Header.Code = intCons.ACK_RESULT_SUCC
	*ack.UserId = user.GetId()
	*ack.NickName = user.GetNickName()
	*ack.Chip = user.GetDiamond()

	a.WriteMsg(ack)
	return

}

