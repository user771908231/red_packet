package newProto

import "casino_login/msg/protogo"

func NewGame_AckQuickConn() *loginproto.Game_AckQuickConn {
	ret := &loginproto.Game_AckQuickConn{}
	ret.CurrVersion = new(int32)
	ret.DownloadUrl = new(string)
	ret.Header = NewHeader()
	ret.IsMaintain = new(int32)
	ret.IsUpdate = new(int32)
	ret.ReleaseTag = new(int32)
	ret.VersionInfo = new(string)
	ret.GameServer = NewServerInfo()
	return ret
}

func NewServerInfo() *loginproto.ServerInfo {
	server := new(loginproto.ServerInfo)
	server.Ip = new(string)
	server.Port = new(int32)
	server.Status = new(int32)
	return server
}

func NewHeader() *loginproto.ProtoHeader {
	ret := new(loginproto.ProtoHeader)
	ret.Error = new(string)
	ret.Code = new(int32)
	ret.UserId = new(uint32)
	return ret
}