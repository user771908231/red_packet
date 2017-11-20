package main

import (
	"log"
	"testing"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_paoyao/service/niuniuService"
	"casino_paoyao/service/niuniu"
	"reflect"
	"casino_common/test"
)

func init() {
	test.TestInit()
}

var gater *Agent

type Agent struct {
}

func (gate *Agent) Close() {
}
func (gate *Agent) Destroy() {
}

func (gate *Agent) WriteMsg(msg interface{}) {
	val := reflect.ValueOf(msg)
	log.Printf("%v : %v", val.Type().String(), msg)
}
func (gate *Agent) SetUserData(data interface{}) {
}
func (gate *Agent) UserData() interface{} {
	return nil
}
func (gate *Agent) RemoteAddr() interface{} {
	return nil
}

//测试创建房间、进房
func TestCreateDeskHandler(t *testing.T) {
	//创建房间
	paoyaoService.CreateDeskHandler(&ddproto.PaoyaoCreateDeskReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(1),
		},
		Option: &ddproto.PaoyaoniuDeskOption{
			MinUser: proto.Int32(2),
			MaxUser: proto.Int32(6),
			MaxCircle: proto.Int32(10),
			HasFlower: proto.Bool(false),
			BankRule: ddproto.PaoyaoniuEnumBankerRule_DING_ZHUANG.Enum(),
			IsFlowerPlay:proto.Bool(true),
			IsJiaoFenJiaBei: proto.Bool(true),
		},
	},gater)

	user,_ := paoyao.FindUserById(1)
	desk_number := user.Desk.GetPwd()

	//进房
	paoyaoService.EnterDeskHandler(&ddproto.PaoyaoEnterDeskReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(1),
		},
		DeskNumber: proto.String(desk_number),
	},gater)

	//第二个用户进房
	paoyaoService.EnterDeskHandler(&ddproto.PaoyaoEnterDeskReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(2),
		},
		DeskNumber: proto.String(desk_number),
	},gater)

	//第二个用户准备
	paoyaoService.ReadyHandler(&ddproto.PaoyaoSwitchReadyReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(2),
		},
		IsReady: proto.Bool(true),
	}, gater)

	//房主开始游戏
	paoyaoService.ReadyHandler(&ddproto.PaoyaoSwitchReadyReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(1),
		},
		IsReady: proto.Bool(true),
	}, gater)

	paoyaoService.JiabeiHandler(&ddproto.PaoyaoJiabeiReq{
		Header: &ddproto.ProtoHeader{
			UserId:proto.Uint32(2),
		},
		Score: proto.Int32(2),
	}, gater)

}
