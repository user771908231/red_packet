package newProto

import (
	mjProto "casino_majiang/msg/protogo"
)

func SuccessHeader() *mjProto.ProtoHeader {
	header := NewHeader()
	*header.Code = 0
	return header
}

func ErrorHeader() *mjProto.ProtoHeader {
	header := NewHeader()
	*header.Code = -1
	return header
}

func MakeHeader(header *mjProto.ProtoHeader, code int32, error string) {
	if ( header == nil ) {
		header = new(mjProto.ProtoHeader)
	}
	if ( header.Code == nil ) {
		header.Code = new(int32)
	}
	if ( header.Error == nil ) {
		header.Error = new(string)
	}

	*header.Code = code
	*header.Error = error
}

func NewHeader() *mjProto.ProtoHeader {
	ret := &mjProto.ProtoHeader{}
	ret.UserId = new(uint32)
	ret.Code = new(int32)
	return ret
}

func NewGame_AckCreateRoom() *mjProto.Game_AckCreateRoom {
	ret := &mjProto.Game_AckCreateRoom{}
	ret.Header = NewHeader()
	ret.DeskId = new(int32)
	ret.Password = new(string);
	ret.UserBalance = new(int64)
	ret.CreateFee = new(int64)
	return ret
}

func NewGame_AckEnterRoom() *mjProto.Game_AckEnterRoom {
	ret := &mjProto.Game_AckEnterRoom{}
	ret.Header = NewHeader()
	return ret
}

func NewRoomTypeInfo() *mjProto.RoomTypeInfo {
	ret := &mjProto.RoomTypeInfo{}
	ret.BaseValue = new(int64)
	ret.BoardsCout = new(int32)
	ret.CapMax = new(int64)
	ret.CardsNum = new(int32)
	ret.MjRoomType = new(mjProto.MJRoomType)
	ret.PlayOptions = NewPlayOptions()
	ret.Settlement = new(int32)
	return ret
}

func NewPlayOptions() *mjProto.PlayOptions {
	ret := &mjProto.PlayOptions{}
	ret.DianGangHuaRadio = new(int32)
	ret.HuRadio = new(int32)
	ret.ZiMoRadio = new(int32)
	return ret
}

func NewGame_SendGameInfo() *mjProto.Game_SendGameInfo {
	ret := &mjProto.Game_SendGameInfo{}
	ret.Header = NewHeader()
	return ret
}