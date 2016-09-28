package newProto

import (
	mjProto "casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
)

func SuccessHeader() *mjProto.ProtoHeader {
	header := new(mjProto.ProtoHeader)
	header.Code = new(int32)
	*header.Code = 0
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

func NewGame_AckCreateRoom() *mjProto.Game_AckCreateRoom {
	ret := &mjProto.Game_AckCreateRoom{}

	return ret
}

//返回一个麻将room
func NewMjRoom() *majiang.MjRoom {
	ret := &majiang.MjRoom{}
	return ret
}

//返回一个麻将
func NewMjDesk() *majiang.MjDesk {
	ret := &majiang.MjDesk{}
	return ret
}