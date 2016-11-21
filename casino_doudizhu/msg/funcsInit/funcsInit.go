package newProto

import "casino_doudizhu/msg/protogo"

func NewHeader() *ddzproto.ProtoHeader {
	ret := &ddzproto.ProtoHeader{}
	ret.UserId = new(uint32)
	ret.Code = new(int32)
	ret.Error = new(string)
	return ret
}

func NewGame_AckLogin() *ddzproto.Game_AckLogin {
	ret := &ddzproto.Game_AckLogin{}
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.NickName = new(string)
	ret.RoomPassword = new(string)
	ret.CostCreateRoom = new(int64)
	ret.CostRebuy = new(int64)
	ret.Championship = new(bool)
	ret.Chip = new(int64)
	ret.MailCount = new(int32)
	ret.Notice = new(string)
	ret.GameStatus = new(int32)
	return ret
}

//创建房间
func NewGame_AckCreateRoom() *ddzproto.Game_AckCreateRoom {
	ret := new(ddzproto.Game_AckCreateRoom)
	return ret
}

//进入房间成功
func NewGame_AckEnterRoom() *ddzproto.Game_AckEnterRoom {
	ret := new(ddzproto.Game_AckEnterRoom)
	ret.Header = NewHeader()
	return ret
}