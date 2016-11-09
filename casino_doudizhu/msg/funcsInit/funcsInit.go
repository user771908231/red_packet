package newProto

import "casino_doudizhu/msg/protogo"

func NewHeader() *ddzproto.ProtoHeader {
	ret := &ddzproto.ProtoHeader{}
	ret.UserId = new(uint32)
	ret.Code = new(int32)
	ret.Error = new(string)
	return ret
}