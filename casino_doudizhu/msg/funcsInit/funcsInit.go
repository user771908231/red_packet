package newProto

import (
	"casino_doudizhu/msg/protogo"
)

func NewDdzQuickConn() *ddzproto.DdzQuickConn {
	ret := new(ddzproto.DdzQuickConn)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	ret.ChannelId = new(string)
	ret.CurrVersion = new(int32)
	ret.GameId = new(int32)
	ret.LanguageId = new(int32)
	return ret
}

func NewDdzAckQuickConn() *ddzproto.DdzAckQuickConn {
	ret := new(ddzproto.DdzAckQuickConn)
	ret.CurrVersion = new(int32)
	ret.Header = NewHeader()
	ret.DownloadUrl = new(string)
	ret.GameServer = NewServerInfo()
	ret.IsMaintain = new(int32)
	ret.IsUpdate = new(int32)
	ret.MaintainMsg = new(string)
	ret.ReleaseTag = new(int32)
	ret.VersionInfo = new(string)
	return ret
}

func NewServerInfo() *ddzproto.ServerInfo {
	ret := new(ddzproto.ServerInfo)
	ret.Ip = new(string)
	ret.Port = new(int32)
	ret.Status = new(int32)
	return ret
}

func NewWeixinInfo() *ddzproto.WeixinInfo {
	ret := new(ddzproto.WeixinInfo)
	ret.City = new(string)
	ret.HeadUrl = new(string)
	ret.NickName = new(string)
	ret.OpenId = new(string)
	ret.Sex = new(int32)
	ret.UnionId = new(string)
	return ret

}

func NewDdzLogin() *ddzproto.DdzLogin {
	ret := new(ddzproto.DdzLogin)
	ret.Header = NewHeader()
	ret.ProtoVersion = new(int32)
	ret.UserId = new(uint32)
	ret.WxInfo = NewWeixinInfo()
	return ret
}

func NewDdzCreateRoom() *ddzproto.DdzCreateRoom {
	ret := new(ddzproto.DdzCreateRoom)
	ret.Header = NewHeader()
	ret.RoomTypeInfo = NewRoomTypeInfo()
	return ret
}

func NewRoomTypeInfo() *ddzproto.RoomTypeInfo {
	ret := new(ddzproto.RoomTypeInfo)
	ret.BaseValue = new(int64)
	ret.BoardsCount = new(int32)
	ret.CapMax = new(int64)
	ret.IsJiaoFen = new(bool)
	ret.RoomType = new(ddzproto.DDZRoomType)
	return ret
}

func NewHeader() *ddzproto.ProtoHeader {
	ret := new(ddzproto.ProtoHeader)
	ret.UserId = new(uint32)
	ret.Code = new(int32)
	ret.Error = new(string)
	return ret
}

func NewGame_AckLogin() *ddzproto.DdzAckLogin {
	ret := new(ddzproto.DdzAckLogin)
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

//发送过的回复
func NewDdzActGuoAck() *ddzproto.DdzActGuoAck {
	ret := new(ddzproto.DdzActGuoAck)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}


//发送发牌
func NewDdzOverTurn() *ddzproto.DdzOverTurn {
	ret := new(ddzproto.DdzOverTurn)
	return ret
}


//准备回复
func NewDdzAckReady() *ddzproto.DdzAckReady {
	ret := new(ddzproto.DdzAckReady)
	ret.UserId = new(uint32)
	ret.Msg = new(string)
	return ret
}


//发送消息
func NewGameMessage() *ddzproto.DdzSendMessage {
	ret := new(ddzproto.DdzSendMessage)
	ret.Header = NewHeader()
	ret.Id = new(int32)
	ret.Msg = new(string)
	ret.MsgType = new(int32)
	ret.UserId = new(uint32)
	return ret
}

func NewDdzEnterRoom() *ddzproto.DdzEnterRoom {
	ret := new(ddzproto.DdzEnterRoom)
	ret.Header = NewHeader()
	ret.MatchId = new(int32)
	ret.PassWord = new(string)
	ret.TableId = new(int32)
	ret.UserId = new(uint32)
	return ret
}


func NewDdzSendGameInfo() *ddzproto.DdzSendGameInfo {
	ret := new(ddzproto.DdzSendGameInfo)
	ret.Header = NewHeader()
	ret.IsReconnect = new(int32)
	ret.SenderUserId = new(uint32)
	return ret
}

func NewDdzReady() *ddzproto.DdzReady {
	ret := new(ddzproto.DdzReady)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.IsShowHandPokers = new(bool)
	return ret
}

func NewDdzOpening() *ddzproto.DdzOpening {
	ret := new(ddzproto.DdzOpening)
	ret.Header = NewHeader()
	return ret
}

