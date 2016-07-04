package protoUtils

import (
	"casino_server/msg/bbprotogo"
	"casino_server/conf/intCons"
)

/**
	得到一个返回正确的header
 */
func GetSuccHeader() *bbproto.ProtoHeader {
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_SUCC
	return result
}

func GetSuccHeaderwithMsgUserid(userId *uint32, msg *string) *bbproto.ProtoHeader {
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_SUCC
	result.Error = msg
	result.UserId = userId
	return result
}

/**
	得到一个返回错误的header
 */

func GetErrorHeader() *bbproto.ProtoHeader {
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_FAIL
	return result
}

func GetErrorHeaderWithMsg(msg *string) *bbproto.ProtoHeader {
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_FAIL
	result.Error = msg
	return result
}

func GetErrorHeaderWithMsgUserid(userId *uint32, msg *string) *bbproto.ProtoHeader {
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_FAIL
	result.Error = msg
	result.UserId = userId
	return result
}



