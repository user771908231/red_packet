package newProto

import "casino_doudizhu/msg/protogo"

func NewHeader() *ddzproto.ProtoHeader {
	ret := &ddzproto.ProtoHeader{}
	ret.UserId = new(uint32)
	ret.Code = new(int32)
	ret.Error = new(string)
	return ret
}

func NewGame_AckLogin() *ddzproto.DdzAckLogin {
	ret := &ddzproto.DdzAckLogin{}
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
func NewGame_AckCreateRoom() *ddzproto.DdzAckCreateRoom {
	ret := new(ddzproto.DdzAckCreateRoom)
	return ret
}

//进入房间成功
func NewGame_AckEnterRoom() *ddzproto.DdzAckEnterRoom {
	ret := new(ddzproto.DdzAckEnterRoom)
	ret.Header = NewHeader()
	return ret
}