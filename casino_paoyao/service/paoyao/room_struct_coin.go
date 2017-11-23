package paoyao

import (
	"errors"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/utils/db"
	"casino_paoyao/conf/config"
)

//查找一个空闲房间，如果没有空闲的则创建一个新房间
func (room *Room) GetFreeCoinDesk() (*Desk, error) {
	if room.GetRoomId() == 0 {
		return nil, errors.New("参数非法！")
	}
	for _,desk := range room.Desks{
		if desk == nil || desk.DeskOption == nil {
			continue
		}

		if len(desk.Users) < int(desk.DeskOption.GetGammerNum()) {
			return desk, nil
		}
	}

	//未找到有空闲位置的房间，则创建一个新房间
	new_desk, err := room.NewCoinDesk()
	if new_desk != nil && err == nil {
		//创建一个新房间
		room.Desks = append(room.Desks, new_desk)
		return new_desk, nil
	}

	return nil, errors.New("未找到空闲房间！")
}

//创建一个金币桌
func (room *Room) NewCoinDesk() (*Desk, error) {
	new_desk_id,err := db.GetNextSeq(config.DBT_PAOYAO_DESK)
	if err != nil {
		return nil, errors.New("get desk seq id fail.")
	}
	new_game_number,err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
	if err != nil {
		return nil, errors.New("get gamenumber seq id fail.")
	}
	desk_option := &ddproto.PaoyaoDeskOption{
		BoardsCout: proto.Int32(0),
		HasAnimation:proto.Bool(true),
	}

	new_desk := &Desk{
		Room: room,
		PaoyaoSrvDesk: &ddproto.PaoyaoSrvDesk{
			DeskId: proto.Int32(new_desk_id),
			Pwd: proto.String(""),
			GameNumber: proto.Int32(new_game_number),
			RoomId: proto.Int32(room.GetRoomId()),
			Status: ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY.Enum(),
			DeskOption: desk_option,
			CircleNo: proto.Int32(1),
			Owner: proto.Uint32(0),
			IsStart: proto.Bool(false),
			IsOnDissolve: proto.Bool(false),
			DissolveTime: proto.Int64(0),
			OneStartTime: proto.Int64(0),
			AllStartTime: proto.Int64(0),
			DaikaiUser: proto.Uint32(0),
			IsDaikai: proto.Bool(false),
			IsCoinRoom: proto.Bool(true),
			SurplusTime: proto.Int32(0),
		},
		Users: []*User{},
	}

	return new_desk, nil
}
