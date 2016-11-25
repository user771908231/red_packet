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
	ret.CreateFee = new(int64)
	ret.DeskId = new(int32)
	ret.Header = NewHeader()
	ret.Password = new(string)
	ret.UserBalance = new(int64)
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
	ret.UserId = new(uint32)
	ret.ActType = new(ddzproto.ActType)
	ret.CanDouble = new(bool)
	ret.CanOutCards = new(bool)
	ret.PullOrPush = new(int32)
	ret.Time = new(int32)
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

func NewDdzDealCards() *ddzproto.DdzDealCards {
	ret := new(ddzproto.DdzDealCards)
	ret.Header = NewHeader()
	return ret
}

func NewDdzJiaoDiZhu() *ddzproto.DdzJiaoDiZhu {
	ret := new(ddzproto.DdzJiaoDiZhu)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.Jiao = new(bool)
	ret.Score = new(int32)
	return ret
}

func NewDdzJiaoDiZhuAck() *ddzproto.DdzJiaoDiZhuAck {
	ret := new(ddzproto.DdzJiaoDiZhuAck)
	ret.Jiao = new(bool)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzRobDiZhu() *ddzproto.DdzRobDiZhu {
	ret := new(ddzproto.DdzRobDiZhu)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzRobDiZhuAck() *ddzproto.DdzRobDiZhuAck {
	ret := new(ddzproto.DdzRobDiZhuAck)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	ret.Rob = new(bool)
	return ret
}

func NewDdzDouble() *ddzproto.DdzDouble {
	ret := new(ddzproto.DdzDouble)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.Double = new(int32)
	return ret
}

func NewDdzDoubleAck() *ddzproto.DdzDoubleAck {
	ret := new(ddzproto.DdzDoubleAck)
	ret.Double = new(int32)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzShowHandPokers() *ddzproto.DdzShowHandPokers {
	ret := new(ddzproto.DdzShowHandPokers)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzShowHandPokersAck() *ddzproto.DdzShowHandPokersAck {
	ret := new(ddzproto.DdzShowHandPokersAck)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzMenuZhua() *ddzproto.DdzMenuZhua {
	ret := new(ddzproto.DdzMenuZhua)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzMenuZhuaAck() *ddzproto.DdzMenuZhuaAck {
	ret := new(ddzproto.DdzMenuZhuaAck)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzSeeCards() *ddzproto.DdzSeeCards {
	ret := new(ddzproto.DdzSeeCards)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzSeeCardsAck() *ddzproto.DdzSeeCardsAck {
	ret := new(ddzproto.DdzSeeCardsAck)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzPull() *ddzproto.DdzPull {
	ret := new(ddzproto.DdzPull)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	ret.Act = new(int32)
	return ret
}

func NewDdzPullAck() *ddzproto.DdzPullAck {
	ret := new(ddzproto.DdzPullAck)
	ret.Act = new(int32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzOutCards() *ddzproto.DdzOutCards {
	ret := new(ddzproto.DdzOutCards)
	ret.Header = NewHeader()
	return ret
}

func NewDdzOutCardsAck() *ddzproto.DdzOutCardsAck {
	ret := new(ddzproto.DdzOutCardsAck)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.CardType = new(ddzproto.DdzPaiType)
	ret.RemainPokers = new(int32)
	return ret
}

func NewDdzActGuo() *ddzproto.DdzActGuo {
	ret := new(ddzproto.DdzActGuo)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzStartPlay() *ddzproto.DdzStartPlay {
	ret := new(ddzproto.DdzStartPlay)
	ret.Header = NewHeader()
	ret.FootRate = new(int32)
	ret.Dizhu = new(uint32)
	return ret
}

func NewDdzSendCurrentResult() *ddzproto.DdzSendCurrentResult {
	ret := new(ddzproto.DdzSendCurrentResult)
	ret.Header = NewHeader()
	return ret
}

func NewDdzSendEndLottery() *ddzproto.DdzSendEndLottery {
	ret := new(ddzproto.DdzSendEndLottery)
	ret.Header = NewHeader()
	return ret
}

func NewDdzDissolveDesk() *ddzproto.DdzDissolveDesk {
	ret := new(ddzproto.DdzDissolveDesk)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewDdzAckDissolveDesk() *ddzproto.DdzAckDissolveDesk {
	ret := new(ddzproto.DdzAckDissolveDesk)
	ret.UserId = new(uint32)
	ret.DeskId = new(int32)
	ret.PassWord = new(string)
	return ret
}

func NewDdzLeaveDesk() *ddzproto.DdzLeaveDesk {
	ret := new(ddzproto.DdzLeaveDesk)
	ret.Header = NewHeader()
	return ret
}

func NewDdzAckLeaveDesk() *ddzproto.DdzAckLeaveDesk {
	ret := new(ddzproto.DdzAckLeaveDesk)
	ret.Header = NewHeader()
	return ret
}

func NewDdzMessage() *ddzproto.DdzMessage {
	ret := new(ddzproto.DdzMessage)
	ret.Header = NewHeader()
	ret.DeskId = new(int32)
	ret.UserId = new(uint32)
	ret.Id = new(int32)
	ret.Msg = new(string)
	ret.MsgType = new(int32)
	return ret
}

func NewDdzSendMessage() *ddzproto.DdzSendMessage {
	ret := new(ddzproto.DdzSendMessage)
	ret.MsgType = new(int32)
	ret.Msg = new(string)
	ret.Id = new(int32)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzGameRecord() *ddzproto.DdzGameRecord {
	ret := new(ddzproto.DdzGameRecord)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.Id = new(int32)
	return ret
}

func NewDdzAckGameRecord() *ddzproto.DdzAckGameRecord {
	ret := new(ddzproto.DdzAckGameRecord)
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

func NewDdzDeskInfo() *ddzproto.DdzDeskInfo {
	ret := new(ddzproto.DdzDeskInfo)
	ret.GameStatus = new(int32)
	ret.PlayerNum = new(int32)
	ret.ActiveUserId = new(uint32)
	ret.ActionTime = new(int32)
	ret.NInitActionTime = new(int32)
	ret.InitRoomCoin = new(int64)
	ret.CurrPlayCount = new(int32)
	ret.TotalPlayCount = new(int32)
	ret.RoomNumber = new(string)
	ret.DiZhuUserId = new(uint32)
	ret.FootRate = new(int32)
	ret.PlayRate = new(int32)
	return ret
}

func NewPlayerInfo() *ddzproto.PlayerInfo {
	ret := new(ddzproto.PlayerInfo)
	ret.IsDiZhu = new(bool)
	ret.Coin = new(int64)
	ret.NickName = new(string)
	ret.Sex = new(int32)
	ret.UserId = new(uint32)
	ret.IsOwner = new(bool)
	ret.BReady = new(int32)
	ret.OnlineStatus = new(int32)
	return ret
}

func NewPoker() *ddzproto.Poker {
	p := new(ddzproto.Poker)
	p.Num = new(int32)
	p.Id = new(int32)
	p.Suit = new(ddzproto.PokerColor)
	return p
}

