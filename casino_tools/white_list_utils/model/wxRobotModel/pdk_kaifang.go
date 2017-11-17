package wxRobotModel

import (
	"github.com/songtianyi/wechat-go/wxweb"
	"casino_common/proto/ddproto"
	"fmt"
	"github.com/golang/protobuf/proto"
	"casino_common/common/service/rpcService"
	"casino_common/common/consts"
	"casino_common/common/log"
	"golang.org/x/net/context"
)

//跑得快配置
var PdkConf = CreateConfig{
	[][]string{
		[]string{"十五张跑得快", "十五张", "15张"},
		[]string{"经典跑得快", "十六张", "16张"},
	},
	[][]string{
		[]string{"20局", "二十局"},
		[]string{"10局", "十局"},
	},
	[][]string{
		[]string{"2人", "二人", "两人"},
		[]string{"3人", "三人"},
	},
	[][]string{
		[]string{"首出黑桃3","首出黑桃三"},
		[]string{"不出黑桃3"},
	},
	[][]string{
		[]string{"不抓鸟"},
		[]string{"红桃10抓鸟", "红桃十抓鸟", "抓鸟"},
	},
	[][]string{
		[]string{"不显示余牌"},
		[]string{"显示余牌"},
	},
}

//经典跑得快开房
func DoPdkKaifang(group_info *T_WECHAT_GROUP_INFO, session *wxweb.Session, msg *wxweb.ReceivedMessage, contact *wxweb.User) {
	var opt_gamer_number int = 2
	var opt_room_type ddproto.PdkEnumRoomType = ddproto.PdkEnumRoomType_PDK_T_FIFTEEN_PDK
	var opt_circle_num int = 20
	var opt_ht3_shouchu bool = false
	var opt_ht10_zhuaniao bool = false
	var opt_show_yupai bool = false

	//解析关键词
	pdk_keywords := PdkConf.GetKeywords(group_info.GameType, msg.Content)
	for _, v := range pdk_keywords {
		switch v {
		case "十五张跑得快":
			opt_room_type = ddproto.PdkEnumRoomType_PDK_T_FIFTEEN_PDK
		case "经典跑得快":
			opt_room_type = ddproto.PdkEnumRoomType_PDK_T_NORMAL_PDK
		case "10局":
			opt_circle_num = 10
		case "20局":
			opt_circle_num = 20
		case "2人":
			opt_gamer_number = 2
		case "3人":
			opt_gamer_number = 3
		case "首出黑桃3":
			opt_ht3_shouchu = true
		case "红桃10抓鸟":
			opt_ht10_zhuaniao = true
		case "显示余牌":
			opt_show_yupai = true
		}
	}

	//检查是否有空闲房间
	ex_room := group_info.GetFreeRoom(opt_gamer_number, pdk_keywords)
	if ex_room != nil {
		session.SendText(fmt.Sprintf("房间:%v\n玩法:%v\n空闲位置:%d", ex_room.GetPassword(), ex_room.GetTips(), opt_gamer_number-len(ex_room.GetUsers())), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	rpc_req := &ddproto.PdkReqCreateDesk{
		Header: &ddproto.ProtoHeader{
			UserId: proto.Uint32(group_info.OwnerId),
		},
		RoomTypeInfo: &ddproto.PdkBaseRoomTypeInfo{
			RoomType: opt_room_type.Enum(),
			BoardsCount: proto.Int32(int32(opt_circle_num)),
			UserCountLimit:proto.Int32(int32(opt_gamer_number)),
			IsDaikai: proto.Bool(true),
			IsShowCardsNum:proto.Bool(opt_show_yupai),
			IsZhuaNiao:proto.Bool(opt_ht10_zhuaniao),
			IsSpadeThree: proto.Bool(opt_ht3_shouchu),
		},
	}
	res,err := rpcService.GetPdk().CreateRoom(context.Background(), rpc_req)
	log.T("rpc req:%v res:%v res-err:%v", rpc_req, res, err)

	if err != nil {
		session.SendText(fmt.Sprintf("开房失败，请联系管理员(错误信息:%s)", err), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	if res == nil {
		session.SendText("开房失败，请联系管理员(错误信息:res is nil.)", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	if res.Header.GetCode() == consts.ACK_RESULT_SUCC {
		ex_room := group_info.GetFreeRoom(opt_gamer_number, pdk_keywords)
		if ex_room != nil {
			session.SendText(fmt.Sprintf("房间:%v\n玩法:%v\n空闲位置:%d", ex_room.GetPassword(), ex_room.GetTips(), opt_gamer_number-len(ex_room.GetUsers())), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		}else {
			session.SendText(fmt.Sprintf("开房失败，请联系管理员。"), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		}
		return
	}else {
		session.SendText(fmt.Sprintf("开房失败，错误:%v", res.Header.GetError()), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}
}
