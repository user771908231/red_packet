package bbproto

//返回一个奖励的bean
func NewGame_TounamentRewardsBean() *Game_TounamentRewardsBean {
	ret := &Game_TounamentRewardsBean{}
	ret.IconPath = new(string)
	ret.Rewards = new(string)
	return ret
}


//初始化一个rankBean
func NewGame_TounamentRankBean() *Game_TounamentRankBean {
	ret := &Game_TounamentRankBean{}
	ret.Coin = new(int64)
	ret.Place = new(int32)
	ret.PlayerImage = new(string)
	ret.PlayerName = new(string)
	return ret
}

func NewGame_TounamentRank() *Game_TounamentRank {
	ret := &Game_TounamentRank{}
	ret.MatchId = new(int32)
	return ret
}
//初始化一个

func NewGame_MatchItem() *Game_MatchItem {
	ret := &Game_MatchItem{}
	ret.CostFee = new(int64)
	ret.Person = new(int32)
	ret.Status = new(int32)
	ret.Person = new(int32)
	ret.Time = new(string)
	ret.Title = new(string)
	ret.Type = new(int32)
	ret.CanInto = new(bool)
	ret.MatchId = new(int32)
	return ret
}

func NewGame_MatchList() *Game_MatchList {
	ret := &Game_MatchList{}
	ret.Result = new(int32)
	ret.HelpMessage = new(string)
	return ret
}

func NewGame_AckLogin() *Game_AckLogin {
	ret := &Game_AckLogin{}
	ret.GameStatus = new(int32)
	ret.MatchId = new(int32)
	ret.Notice = new(string)
	ret.RoomPassword = new(string)
	ret.TableId = new(int32)
	ret.CostCreateRoom = new(int64)
	ret.CostRebuy = new(int64)
	ret.Championship = new(bool)
	ret.Chip = new(int64)
	return ret

}

func NewGame_TestResult() *Game_TestResult {
	ret := &Game_TestResult{}
	ret.Tableid = new(int32)
	ret.Rank = new(int32)
	ret.CanRebuy = new(bool)
	ret.RebuyCount = new(int32)
	ret.RankUserCount = new(int32)
	return ret
}

func NewGame_InitCard() *Game_InitCard {
	//设置默认值
	ret := &Game_InitCard{}
	ret.Tableid = new(int32)
	ret.ActionTime = new(int32)
	ret.DelayTime = new(int32)
	ret.NextUser = new(int32)
	ret.MinRaise = new(int64)
	ret.CurrPlayCount = new(int32)
	ret.TotalPlayCount = new(int32)
	return ret
}

func NewGame_EndLottery() *Game_EndLottery {
	gel := &Game_EndLottery{}
	gel.UserId = new(uint32)
	gel.Coin = new(int64)
	gel.BigWin = new(bool)
	gel.Owner = new(bool)
	gel.Rolename = new(string)
	gel.Seat = new(int32)
	return gel
}

func NewGame_PreCoin() *Game_PreCoin {
	ret := &Game_PreCoin{}
	ret.Pool = new(int64)
	ret.Tableid = new(int32)
	return ret
}

func NewGame_BlindCoin() *Game_BlindCoin {
	blindB := &Game_BlindCoin{}
	blindB.Banker = new(int32)
	blindB.Bigblindseat = new(int32)
	blindB.Smallblindseat = new(int32)
	return blindB
}

func NewGame_SendGameInfo() *Game_SendGameInfo {
	result := &Game_SendGameInfo{}
	result.TablePlayer = new(int32)
	result.Tableid = new(int32)
	result.BankSeat = new(int32)
	result.ChipSeat = new(int32)
	result.ActionTime = new(int32)
	result.DelayTime = new(int32)
	result.GameStatus = new(int32)
	//result.NRebuyCount
	//result.NAddonCount
	result.Pool = new(int64)
	result.MinRaise = new(int64)
	result.NInitActionTime = new(int32)
	result.NInitDelayTime = new(int32)
	result.Seat = new(int32)
	result.TurnMax = new(int64)
	result.Result = new(int32)
	result.OwnerSeat = new(int32)
	result.PreCoin = new(int64)
	result.SmallBlind = new(int64)
	result.BigBlind = new(int64)
	result.RoomType = new(int32)
	result.InitRoomCoin = new(int64)
	result.CurrPlayCount = new(int32)
	result.TotalPlayCount = new(int32)
	result.SenderUserId = new(uint32)
	result.ShowCsBeginButton = new(bool)
	return result
}

func NewGame_AckRebuy() *Game_AckRebuy {
	ret := &Game_AckRebuy{}
	ret.CurrChip = new(int64)
	ret.RemainCount = new(int32)
	ret.Result = new(int32)
	return ret
}

func NewGame_ACKLeaveDesk() *Game_ACKLeaveDesk {
	ret := &Game_ACKLeaveDesk{}
	ret.Result = new(int32)
	return ret
}

func NewThServerUserSession() *ThServerUserSession {
	result := &ThServerUserSession{}
	result.DeskId = new(int32)
	result.MatchId = new(int32)
	result.UserId = new(uint32)
	result.GameType = new(int32)
	result.GameStatus = new(int32)
	result.RoomKey = new(string)
	result.IsLeave = new(bool)
	result.IsBreak = new(bool)
	return result
}

func NewThServerUser() *ThServerUser {
	user := &ThServerUser{}
	user.Seat = new(int32)
	user.Status = new(int32)
	user.BreakStatus = new(int32)
	user.WaiTime = new(string)
	user.WaitUUID = new(string)
	user.DeskId = new(int32)
	user.TotalBet = new(int64)
	user.TotalBet4CalcAllin = new(int64)
	user.WinAmount = new(int64)
	user.TurnCoin = new(int64)
	user.HandCoin = new(int64)
	user.RoomCoin = new(int64)
	user.UserId = new(uint32)
	user.GameNumber = new(int32)
	user.IsBreak = new(bool)
	user.IsLeave = new(bool)
	user.TotalRoomCoin = new(int64)
	user.LotteryCheck = new(bool)
	user.IsShowCard = new(bool)
	user.MatchId = new(int32)
	user.CloseCheck = new(int32)
	user.IsRebuy = new(bool)
	user.WaitRebuyFlag = new(bool)
	return user
}

func NewThServerAllInJack() *ThServerAllInJackpot {
	ret := &ThServerAllInJackpot{}
	ret.AllInAmount = new(int64)
	ret.Jackpopt = new(int64)
	ret.ThroundCount = new(int32)
	ret.UserId = new(uint32)
	return ret
}

func NewThServerDesk() *ThServerDesk {
	result := &ThServerDesk{}
	result.Id = new(int32)
	result.MatchId = new(int32)
	result.DeskOwner = new(uint32)
	result.RoomKey = new(string)
	result.CreateFee = new(int64)
	result.DeskType = new(int32)
	result.GameType = new(int32)
	result.InitRoomCoin = new(int64)
	result.JuCount = new(int32)
	result.JuCountNow = new(int32)
	result.PreCoin = new(int64)
	result.SmallBlindCoin = new(int64)
	result.BigBlindCoin = new(int64)
	result.BlindLevel = new(int32)
	result.RebuyCountLimit = new(int32)
	result.RebuyBlindLevelLimit = new(int32)
	result.Dealer = new(uint32)
	result.BigBlind = new(uint32)
	result.SmallBlind = new(uint32)
	result.RaiseUserId = new(uint32)
	result.NewRoundFirstBetUser = new(uint32)
	result.BetUserNow = new(uint32)
	result.GameNumber = new(int32)
	result.UserCount = new(int32)
	result.UserCountOnline = new(int32)
	result.Status = new(int32)
	result.BetAmountNow = new(int64)
	result.RoundCount = new(int32)
	result.Jackpot = new(int64)
	result.EdgeJackpot = new(int64)
	result.MinRaise = new(int64)

	result.SendFlop = new(bool)
	result.SendTurn = new(bool)
	result.SendRive = new(bool)
	return result
}

func NewGame_SendChangeDeskOwner() *Game_SendChangeDeskOwner {
	result := &Game_SendChangeDeskOwner{}
	result.DeskId = new(int32)
	result.OldOwner = new(uint32)
	result.NewOwner = new(uint32)
	result.OldOwnerSeat = new(int32)
	result.NewOwnerSeat = new(int32)
	return result
}

func NewRUNNING_DESKKEYS() *RUNNING_DESKKEYS {
	ret := &RUNNING_DESKKEYS{}
	return ret
}

func NewGame_TounamentPlayerRank() *Game_TounamentPlayerRank {
	ret := &Game_TounamentPlayerRank{}
	ret.Message = new(string)
	ret.PlayerRank = new(int32)
	return ret
}

func NewGame_TounamentBlind() *Game_TounamentBlind {
	ret := &Game_TounamentBlind{}
	return ret
}

//
func NewGame_TounamentBlindBean() *Game_TounamentBlindBean {
	ret := &Game_TounamentBlindBean{}
	ret.Ante = new(string)
	ret.BlindLevel = new(string)
	ret.CanRebuy = new(bool)
	ret.RaiseTime = new(string)
	ret.SmallBlind = new(string)
	return ret
}

// func New
func NewGame_ChampionshipGameOver() *Game_ChampionshipGameOver {
	ret := &Game_ChampionshipGameOver{}
	ret.Coin = new(int64)
	ret.HeadUrl = new(string)
	ret.UserName = new(string)
	return ret
}

func NewGame_AckGameRecord() *Game_AckGameRecord {
	ret := &Game_AckGameRecord{}
	ret.UserId = new(uint32)
	ret.Result = new(int32)
	return ret
}

func NewACKQuickConn() *ACKQuickConn {
	result := &ACKQuickConn{}
	result.CoinCnt = new(int64)
	result.UserName = new(string)
	result.UserId = new(uint32)
	result.NickName = new(string)
	result.AckResult = new(int32)
	result.ReleaseTag = new(int32)
	return result

}

func NewGame_AckCreateDesk() *Game_AckCreateDesk {
	result := &Game_AckCreateDesk{}
	result.Result = new(int32)
	result.Password = new(string)
	result.DeskId = new(int32)
	result.CreateFee = new(int64)
	result.UserBalance = new(int64)
	return result
}

func NewGame_BeanUserRecord() *Game_BeanUserRecord {
	ru := &Game_BeanUserRecord{}
	ru.UserId = new(uint32)
	ru.NickName = new(string)
	ru.WinAmount = new(int64)
	return ru
}

func NewGame_BeanGameRecord() *Game_BeanGameRecord {
	r := &Game_BeanGameRecord{}
	r.DeskId = new(int32)
	r.Id = new(int32)
	r.BeginTime = new(string)
	return r
}
