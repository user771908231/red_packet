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
	ret.Error = new(string)
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

func NewGame_AckReady() *mjProto.Game_AckReady {
	ret := &mjProto.Game_AckReady{}
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewGame_AckEnterRoom() *mjProto.Game_AckEnterRoom {
	ret := &mjProto.Game_AckEnterRoom{}
	ret.Header = NewHeader()
	return ret
}

func NewGame_BroadcastBeginDingQue() *mjProto.Game_BroadcastBeginDingQue {
	ret := &mjProto.Game_BroadcastBeginDingQue{}
	return ret
}

func NewGame_DingQue() *mjProto.Game_DingQue {
	que := &mjProto.Game_DingQue{}
	que.Header = NewHeader()
	que.Color = new(int32)
	que.UserId = new(uint32)
	return que
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
	ret.SenderUserId = new(uint32)
	return ret
}

//返回deskGameInfo
func NewDeskGameInfo() *mjProto.DeskGameInfo {
	ret := &mjProto.DeskGameInfo{}
	ret.GameStatus = new(int32)
	ret.PlayerNum = new(int32)
	ret.ActiveUserId = new(uint32)
	ret.ActionTime = new(int32)
	ret.DelayTime = new(int32)
	ret.NInitActionTime = new(int32)
	ret.NInitDelayTime = new(int32)
	ret.InitRoomCoin = new(int64)
	ret.CurrPlayCount = new(int32)
	ret.TotalPlayCount = new(int32)
	ret.RoomNumber = new(string)
	ret.Banker = new(uint32)
	ret.RemainCards = new(int32)
	return ret
}

//返回一个playerInfo
func NewPlayerInfo() *mjProto.PlayerInfo {
	info := &mjProto.PlayerInfo{}
	info.IsBanker = new(bool)
	info.Coin = new(int64)
	info.NickName = new(string)
	info.UserId = new(uint32)
	info.IsOwner = new(bool)
	info.BReady = new(int32)
	info.BDingQue = new(int32)
	info.BExchanged = new(int32)
	info.NHuPai = new(int32)
	info.NickName = new(string)
	info.QuePai = new(int32)
	return info
}

//麻将card
func NewPlayerCard() *mjProto.PlayerCard {
	card := &mjProto.PlayerCard{}
	card.UserId = new(uint32)
	return card
}

func NewCardInfo() *mjProto.CardInfo {
	cardInfo := &mjProto.CardInfo{}
	cardInfo.Id = new(int32)
	cardInfo.Type = new(int32)
	cardInfo.Value = new(int32)
	return cardInfo
}

func NewComposeCard() *mjProto.ComposeCard {
	ret := &mjProto.ComposeCard{}
	ret.Type = new(int32)
	ret.Value = new(int32)
	return ret
}

func NewGame_OverTurn() *mjProto.Game_OverTurn {
	ret := &mjProto.Game_OverTurn{}
	ret.ActType = new(int32)
	ret.CanGang = new(bool)
	ret.CanPeng = new(bool)
	ret.CanHu = new(bool)
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	ret.NextUserId = new(uint32)
	ret.Time = new(int32)
	ret.PaiCount = new(int32)
	return ret
}

func NewGame_DealCards() *mjProto.Game_DealCards {
	ret := &mjProto.Game_DealCards{}
	ret.Header = NewHeader()
	ret.DealerUserId = new(uint32)
	return ret
}

func NewGame_Opening() *mjProto.Game_Opening {
	ret := &mjProto.Game_Opening{}
	ret.Header = NewHeader()
	return ret
}

func NewGame_AckQuickConn() *mjProto.Game_AckQuickConn {
	ret := &mjProto.Game_AckQuickConn{}
	ret.CurrVersion = new(int32)
	ret.DownloadUrl = new(string)
	ret.Header = NewHeader()
	ret.IsMaintain = new(int32)
	ret.IsUpdate = new(int32)
	ret.ReleaseTag = new(int32)
	ret.VersionInfo = new(string)
	return ret
}

func NewGame_DingQueEnd() *mjProto.Game_DingQueEnd {
	ret := &mjProto.Game_DingQueEnd{}
	return ret
}

func NewDingQueEndBean() *mjProto.DingQueEndBean {
	ret := &mjProto.DingQueEndBean{}
	ret.UserId = new(uint32)
	ret.Flower = new(int32)
	return ret
}

func NewGame_AckActHu() *mjProto.Game_AckActHu {
	ret := &mjProto.Game_AckActHu{}
	ret.Header = NewHeader()
	ret.UserIdIn = new(uint32)
	ret.UserIdOut = new(uint32)
	ret.HuType = new(int32)
	ret.IsZiMo = new(bool)
	return ret
}

func NewGame_AckSendOutCard() *mjProto.Game_AckSendOutCard {
	ret := &mjProto.Game_AckSendOutCard{}
	ret.Header = NewHeader()
	ret.Result = new(int32)
	ret.UserId = new(uint32)
	return ret
}

func NewGame_AckActGang() *mjProto.Game_AckActGang {
	ret := &mjProto.Game_AckActGang{}
	ret.GangType = new(int32)
	ret.Header = NewHeader()
	ret.UserIdIn = new(uint32)
	ret.UserIdOut = new(uint32)
	return ret
}

func NewGame_AckActPeng() *mjProto.Game_AckActPeng {
	ret := &mjProto.Game_AckActPeng{}
	ret.Header = NewHeader()
	ret.UserIdIn = new(uint32)
	ret.UserIdOut = new(uint32)
	ret.JiaoInfos = nil
	return ret
}

func NewGame_AckLogin() *mjProto.Game_AckLogin {
	ret := &mjProto.Game_AckLogin{}
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

func NewWeixinInfo() *mjProto.WeixinInfo {
	ret := &mjProto.WeixinInfo{}
	ret.City = new(string)
	ret.HeadUrl = new(string)
	ret.NickName = new(string)
	ret.OpenId = new(string)
	ret.Sex = new(int32)
	ret.UnionId = new(string)
	return ret
}

func NewGame_SendCurrentResult() *mjProto.Game_SendCurrentResult {
	ret := &mjProto.Game_SendCurrentResult{}
	ret.Header = NewHeader()
	return ret

}

func NewGame_SendEndLottery() *mjProto.Game_SendEndLottery {
	ret := &mjProto.Game_SendEndLottery{}
	ret.Header = NewHeader()
	return ret
}

func NewWinCoinInfo() *mjProto.WinCoinInfo {
	ret := &mjProto.WinCoinInfo{}
	ret.CardTitle = new(string)
	ret.Coin = new(int64)
	ret.HuCount = new(int32)
	ret.IsDealer = new(bool)
	ret.UserId = new(uint32)
	ret.NickName = new(string)
	ret.WinCoin = new(int64)
	return ret
}

func NewEndLotteryInfo() *mjProto.EndLotteryInfo {
	ret := &mjProto.EndLotteryInfo{}
	ret.UserId = new(uint32)
	ret.NickName = new(string)
	ret.BigWin = new(bool)
	ret.IsOwner = new(bool)
	ret.WinCoin = new(int64)
	ret.CountHu = new(int32)
	ret.CountZiMo = new(int32)
	ret.CountDianPao = new(int32)
	ret.CountAnGang = new(int32)
	ret.CountMingGang = new(int32)
	ret.CountDianGang = new(int32)
	ret.CountChaJiao = new(int32)
	return ret
}

func NewGame_AckDissolveDesk() *mjProto.Game_AckDissolveDesk {
	ret := &mjProto.Game_AckDissolveDesk{}
	ret.DeskId = new(int32)
	ret.Header = NewHeader()
	ret.PassWord = new(string)
	ret.UserId = new(uint32)
	return ret
}

func NewGame_AckGameRecord() *mjProto.Game_AckGameRecord {
	ret := &mjProto.Game_AckGameRecord{}
	ret.UserId = new(uint32)
	ret.Header = NewHeader()
	return ret
}

//麻将战绩的bean
func NewBeanGameRecord() *mjProto.BeanGameRecord {
	ret := &mjProto.BeanGameRecord{}
	ret.Header = NewHeader()
	ret.DeskId = new(int32)
	ret.BeginTime = new(string)
	ret.Id = new(int32)
	return ret
}

func NewBeanUserRecord() *mjProto.BeanUserRecord {
	result := &mjProto.BeanUserRecord{}
	result.NickName = new(string)
	result.Header = NewHeader()
	result.UserId = new(uint32)
	result.WinAmount = new(int64)
	return result
}

func NewGame_SendMessage() *mjProto.Game_SendMessage {
	result := &mjProto.Game_SendMessage{}
	result.Header = NewHeader()
	result.Id = new(int32)
	result.Msg = new(string)
	result.MsgType = new(int32)
	result.UserId = new(uint32)
	return result
}

func NewGame_AckActGuo() *mjProto.Game_AckActGuo {
	ret := &mjProto.Game_AckActGuo{}
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewGame_AckExchangeCards() *mjProto.Game_AckExchangeCards {
	ret := &mjProto.Game_AckExchangeCards{}
	ret.Header = NewHeader()
	ret.UserId = new(uint32)
	return ret
}

func NewGame_ExchangeCardsEnd() *mjProto.Game_ExchangeCardsEnd {
	ret := &mjProto.Game_ExchangeCardsEnd{}
	ret.Header = NewHeader()
	ret.ExchangeNum = new(int32)
	ret.ExchangeType = new(int32)
	return ret
}

func NewGame_BroadcastBeginExchange() *mjProto.Game_BroadcastBeginExchange {
	ret := &mjProto.Game_BroadcastBeginExchange{}
	return ret
}