package newProto

import (
	"github.com/golang/protobuf/proto"
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
)

//初始化-RoomInfo
func NewZjhRoomInfo() *ddproto.ZjhBaseRoomInfo {
	info := ddproto.ZjhBaseRoomInfo{
		Header:commonNewPorot.NewHeader(),
		RoomId:proto.Int32(0),
		RoomType:ddproto.ZjhEnumRoomType_ROOM_TYPE_NORMAL.Enum(),
		RoomLevel:proto.Int32(0),
		RoomTitle:proto.String(""),
		BaseValue:proto.Int64(0),
		MaxValue:proto.Int64(0),

	}
	return &info
}

//初始化RoomListACK
func NewZjhRoomListAck() *ddproto.ZjhAckRoomList {
	list := ddproto.ZjhAckRoomList{
		Header:commonNewPorot.NewHeader(),
		Rooms:[]*ddproto.ZjhBaseRoomInfo{},
	}
	return &list
}

//msg User Info
func NewZjhUserInfo() *ddproto.ZjhBaseUserInfo {
	user := &ddproto.ZjhBaseUserInfo{
		Uid:new(int32),
		NickName:new(string),
		Coin:new(int64),
		Chip:new(int32),
		State:ddproto.ZjhEnumUserState_USER_IS_STAND.Enum(),
	}

	return user
}

//自动进入房间ACK
func NewZjhDeskStateAck() *ddproto.ZjhDeskStateAck {
	desk_state := &ddproto.ZjhDeskStateAck{
		State:ddproto.ZjhEnumDeskState_DESK_IS_WAIT.Enum(),
		UserList:[]*ddproto.ZjhBaseUserInfo{},
		DeskId:new(int32),
		RoomInfo:NewZjhRoomInfo(),
	}
	return desk_state
}

func NewPlayerInfo() *ddproto.ZjhBasePlayerInfo {
	ret := new(ddproto.ZjhBasePlayerInfo)
	ret.BReady = new(int32)
	ret.Coin = new(int64)
	ret.IsCheckedPokers = new(bool)
	ret.IsFirst = new(bool)
	ret.NickName = new(string)
	ret.OnlineStatus = new(int32)
	ret.SeatIndex = new(int32)
	ret.Sex = new(int32)
	ret.WxInfo = commonNewPorot.NewWeixinInfo()
	ret.UserId = new(uint32)
	ret.Status = new(ddproto.ZjhEnumPlayerGameStatus)
	ret.IsWatch = new(bool)
	ret.IsDiscard = new(bool)
	ret.IsLost = new(bool)
	return ret
}

func NewWeixinInfo() *ddproto.WeixinInfo {
	ret := new(ddproto.WeixinInfo)
	ret.City = new(string)
	ret.HeadUrl = new(string)
	ret.NickName = new(string)
	ret.OpenId = new(string)
	ret.Sex = new(int32)
	ret.UnionId = new(string)
	return ret

}

func NewPoker() *ddproto.ClientBasePoker {
	p := new(ddproto.ClientBasePoker)
	p.Num = new(int32)
	p.Id = new(int32)
	p.Suit = new(ddproto.CommonEnumPokerColor)
	return p
}

func NewZjhDeskInfo() *ddproto.ZjhBaseDeskInfo {
	ret := new(ddproto.ZjhBaseDeskInfo)
	ret.GameStatus = new(int32)
	ret.PlayerNum = new(int32)
	ret.ActiveUserId = new(uint32)
	ret.ActionTime = new(int32)
	ret.NInitActionTime = new(int32)
	ret.InitRoomCoin = new(int64)
	ret.CurrPlayCount = new(int32)
	ret.TotalPlayCount = new(int32)
	ret.RoomNumber = new(string)
	ret.PlayRate = new(int32)
	ret.CurrRoundCount = new(int32)
	ret.RoomInfo = NewZjhRoomInfo()
	ret.TotalRoundCount = new(int32)
	ret.CuurBaseValue = new(int64)
	ret.CuurCoinCount = new(int64)
	ret.XuepinChipCount = new(int32)
	ret.XuepinBaseValue = new(int64)
	return ret
}

func NewZjhSendGameInfo() *ddproto.ZjhBcGameInfo {
	ret := new(ddproto.ZjhBcGameInfo)
	ret.Header = commonNewPorot.NewHeader()
	ret.ZjhDeskInfo = NewZjhDeskInfo()
	ret.IsReconnect = new(int32)
	ret.SenderUserId = new(uint32)
	return ret
}

//广播：新用户加入房间
func NewZjhBcNewPlayerEnter() *ddproto.ZjhBcNewPlayerEnter {
	ret := new(ddproto.ZjhBcNewPlayerEnter)
	ret.PlayerInfo = NewPlayerInfo()
	return ret
}

//广播：开局
func NewZjhBcOpening() *ddproto.ZjhBcOpening {
	ret := new(ddproto.ZjhBcOpening)
	ret.PlayerInfoList = []*ddproto.ZjhBasePlayerInfo{}
	ret.BaseAnte = new(int64)
	return ret
}
