package internal
//
//import (
//	"reflect"
//	"casino_common/proto/ddproto"
//	"casino_common/common/userService"
//	"casino_common/proto/funcsInit"
//	"casino_common/common/consts"
//	"casino_common/common/log"
//	"github.com/name5566/leaf/gate"
//	"casino_laowangye/service/zjhService"
//	"errors"
//)
//
//func handler(m interface{}, h interface{}) {
//	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
//}
//
//func init() {
//	handler(&ddproto.CommonReqGameLogin{}, HandlerGame_Login)
//}
//
////登录处理
//func HandlerGame_Login(args []interface{}) {
//	m := args[0].(*ddproto.CommonReqGameLogin)
//	a := args[1].(gate.Agent)
//
//	log.T("请求handlerGame_Login  m[%v]", m)
//	weixin := m.GetWxInfo()
//
//	log.Printf("weixin :%v",weixin)
//	//不是初次登录
//	if weixin == nil {
//		//判断uerId
//		userId := m.GetHeader().GetUserId()
//		user := userService.GetUserById(userId)
//		log.Printf("user:%v",user)
//		if user == nil {
//			//登陆失败
//			log.E("没有找到玩家[%v]的信息，login失败", userId)
//			ack := commonNewPorot.NewCommonAckGameLogin()
//			*ack.Header.Code = consts.ACK_RESULT_ERROR
//			a.WriteMsg(ack)
//		} else {
//			//返回登陆成功的结果
//			ack := commonNewPorot.NewCommonAckGameLogin()
//
//			*ack.Header.Code = consts.ACK_RESULT_SUCC
//			*ack.UserId = user.GetId()
//			*ack.NickName = user.GetNickName()
//			*ack.Chip = user.GetDiamond()
//			a.WriteMsg(ack)
//			//断线重连处理
//			ReconnectProcess(user.GetId(),a)
//		}
//		return
//	}
//
//	//1,首先通过weixinInfo 在数据库中查找 用户是否存在，如果用户存在，则表示，登陆成功
//	user := userService.GetUserByOpenId(weixin.GetOpenId())
//	log.Printf("log2:%v",user)
//	if user == nil {
//		//表示数据库中不存在次用户，新增加一个人后返回
//		if weixin.GetOpenId() == "" || weixin.GetHeadUrl() == "" || weixin.GetNickName() == "" {
//			ack := commonNewPorot.NewCommonAckGameLogin()
//
//			*ack.Header.Code = consts.ACK_RESULT_ERROR
//			a.WriteMsg(ack)
//			//断线重连处理
//			ReconnectProcess(user.GetId(),a)
//			return
//		}
//		//如果数据库中不存在用户，那么重新生成一个user
//		user, _ = userService.NewUserAndSave(weixin.GetUnionId(), weixin.GetOpenId(), weixin.GetNickName(), weixin.GetHeadUrl(), weixin.GetSex(), weixin.GetCity())
//		if user == nil {
//			ack := commonNewPorot.NewCommonAckGameLogin()
//
//			*ack.Header.Code = consts.ACK_RESULT_ERROR
//			a.WriteMsg(ack)
//			//断线重连处理
//			ReconnectProcess(user.GetId(),a)
//			return
//		}
//	}
//
//	//返回登陆成功的结果
//	ack := commonNewPorot.NewCommonAckGameLogin()
//	*ack.Header.Code = consts.ACK_RESULT_SUCC
//	*ack.UserId = user.GetId()
//	*ack.NickName = user.GetNickName()
//	*ack.Chip = user.GetDiamond()
//
//	a.WriteMsg(ack)
//	//断线重连处理
//	ReconnectProcess(user.GetId(),a)
//	return
//}
//
////断线重连处理
//func ReconnectProcess(uid uint32, agent gate.Agent) error {
//	desk,user,err := zjhService.GetGameStateByUid(uid)
//	log.Printf("login err:%v",err)
//	if err == nil {
//		//1.先删除原来的连接里保存的userData
//		user.Agent.SetUserData(nil)
//		//2.更新连接
//		user.Agent = agent
//		//3.向客户端发送断线状态
//		desk.SendGameDeskInfo(uid, 2)
//		//4.如果正在游戏
//		if user.GetData().GetBill() > 0 {
//			//1.更新游戏在线状态
//			user.Status = ddproto.ZjhEnumUserStatus_zjh_S_GAMING.Enum()
//			//2.如果轮到断线用户操作则发送操作OverTurn
//			if desk.GetNextUserId() == uid {
//				if desk.IsXuepinOver() {
//					desk.SendNextOT(false)
//				}else {
//					desk.SendNextOT(true)
//				}
//			}
//		}else {
//			user.Status = ddproto.ZjhEnumUserStatus_zjh_S_SITED.Enum()
//		}
//		return errors.New("断线重连处理成功！")
//	}
//	return nil
//}
//
