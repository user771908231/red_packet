package paosangong

import (
	"errors"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/utils/db"
	"casino_paosangong/conf/config"
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

		if len(desk.Users) < int(desk.DeskOption.GetMaxUser()) {
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
	new_desk_id,err := db.GetNextSeq(config.DBT_NIUNIU_DESK)
	if err != nil {
		return nil, errors.New("get desk seq id fail.")
	}
	new_game_number,err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
	if err != nil {
		return nil, errors.New("get gamenumber seq id fail.")
	}
	desk_option := &ddproto.NiuniuDeskOption{
		MinUser: proto.Int32(2),
		MaxUser: proto.Int32(room.GetMaxDeskGammer()),
		MaxCircle: proto.Int32(0),
		HasFlower:proto.Bool(true),
		BankRule: room.BankRule.Enum(),
		IsFlowerPlay:proto.Bool(true),
		IsJiaoFenJiaBei:proto.Bool(true),
		HasAnimation:proto.Bool(true),
		IsCoinRoom:proto.Bool(true),
		BaseScore:proto.Int64(room.GetBaseChip()),
		IsPrivate:proto.Bool(true),
		MinEnterScore:proto.Int64(room.GetEnterCoin()),
		MaxQzScore:proto.Int64(4),
		CoinFee:proto.Int64(room.GetCoinFee()),
		HasCuopai: proto.Bool(false),
	}

	//明牌抢庄模式，开启搓牌动画
	if room.GetBankRule() == ddproto.NiuniuEnumBankerRule_QIANG_ZHUANG {
		*desk_option.HasCuopai = true
	}

	new_desk := &Desk{
		Room: room,
		NiuniuSrvDesk: &ddproto.NiuniuSrvDesk{
			DeskId: proto.Int32(new_desk_id),
			DeskNumber: proto.String(""),
			GameNumber: proto.Int32(new_game_number),
			RoomId: proto.Int32(room.GetRoomId()),
			LastWiner: proto.Uint32(0),
			Status: ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY.Enum(),
			DeskOption: desk_option,
			CircleNo: proto.Int32(1),
			Owner: proto.Uint32(0),
			CurrBanker: proto.Uint32(0),
			IsStart: proto.Bool(false),
			IsOnDissolve: proto.Bool(false),
			DissolveTime: proto.Int64(0),
			OneStartTime: proto.Int64(0),
			AllStartTime: proto.Int64(0),
			DaikaiUser: proto.Uint32(0),
			IsDaikai: proto.Bool(false),
			IsOnGamming: proto.Bool(false),
			IsCoinRoom: proto.Bool(true),
			StartTime: proto.Int64(0),
		},
		Users: []*User{},
	}

	//新增快照索引
	new_desk.NewSnapShot()

	return new_desk, nil
}
