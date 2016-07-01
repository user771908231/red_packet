package protoUtils

import (
	"casino_server/msg/bbprotogo"
	"casino_server/conf/intCons"
)

/**
	得到一个返回正确的header
 */
func GetSuccHeader()*bbproto.ProtoHeader{
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_SUCC
	return result
}


/**
	得到一个返回错误的header
 */

func GetErrorHeader()*bbproto.ProtoHeader{
	result := &bbproto.ProtoHeader{}
	result.Code = &intCons.CODE_FAIL
	return result
}
