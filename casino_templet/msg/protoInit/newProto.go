package newProto

import (
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
)

func NewGame_AckQuickConn() *ddproto.AckQuickConn {
	ret := &ddproto.AckQuickConn{}
	ret.Header = commonNewPorot.NewHeader()
	return ret
}

func NewServerInfo() *ddproto.ServerInfo {
	server := new(ddproto.ServerInfo)
	server.Ip = new(string)
	server.Port = new(int32)
	server.Status = new(int32)
	return server
}
