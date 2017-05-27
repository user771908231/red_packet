package internal

import (
	"reflect"
	"casino_paodekuai/core/data"
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_common/common/consts"
	"src/github.com/golang/protobuf/proto"
	"casino_common/common/Error"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&ddproto.PdkReqCreateDesk{}, createDesk) //创建房间
}

//创建一个Desk
func createDesk(args []interface{}) {
	//
	m := args[0].(*ddproto.PdkReqCreateDesk)
	a := args[1].(gate.Agent)
	//获得参数...
	userId := m.GetHeader().GetUserId()
	roomType := int32(0)
	roomLevel := int32(0)

	room, err := roomMgr.GetRoom(roomType, roomLevel)
	if err != nil {
		ack := &ddproto.PdkAckCreateDesk{
			Header:&ddproto.ProtoHeader{
				UserId:proto.Uint32(userId),
				Code:proto.Int32(Error.GetErrorCode(err)),
				Error:proto.String(Error.GetErrorMsg(err)),
			},
		}
		a.WriteMsg(ack)
	}

	//得到创建的参数
	cfg := &data.PdkDeskCfg{
	}

	//根据参数创建房间
	desk, err := room.CreateDesk(cfg)
	if err != nil {
		ack := &ddproto.PdkAckCreateDesk{
			Header:&ddproto.ProtoHeader{
				UserId:proto.Uint32(userId),
				Code:proto.Int32(Error.GetErrorCode(err)),
				Error:proto.String(Error.GetErrorMsg(err)),
			},
		}
		a.WriteMsg(ack)
	}

	//让房间加入一个玩家
	err = desk.EnterUser(userId)
	if err != nil {
		ack := &ddproto.PdkAckCreateDesk{
			Header:&ddproto.ProtoHeader{
				UserId:proto.Uint32(userId),
				Code:proto.Int32(Error.GetErrorCode(err)),
				Error:proto.String(Error.GetErrorMsg(err)),
			},
		}
		a.WriteMsg(ack)
	}

	gameinfo := &ddproto.PdkBcGameInfo{
		Header:&ddproto.ProtoHeader{
			UserId:proto.Uint32(userId),
			Code:proto.Int32(consts.ACK_RESULT_SUCC),
		},
	}

	//回复消息
	a.WriteMsg(gameinfo)
}

func actReady(args []interface{}) {
	//获得参数
	userId := uint32(0)

	//获得desk
	desk, err := roomMgr.GetDeskBySession(userId)
	if err != nil {
		return
	}

	//开始出牌
	err = desk.ActReady(userId)
	if err != nil {

	}
}

//打牌
func actOut(args []interface{}) {

	//m := args[0] //取出牌

	//获得参数
	userId := uint32(0) //玩家id
	paiIds := []int32{}

	//获得desk
	desk, err := roomMgr.GetDeskBySession(userId)
	if err != nil {
		return
	}

	//开始出牌
	err = desk.ActOut(userId, paiIds)
	if err != nil {

	}
}
