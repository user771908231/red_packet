package bbproto

//初始化盲注bean
func NewGame_TounamentBlindBean() *Game_TounamentBlindBean {
	ret := &Game_TounamentBlindBean{}
	ret.Ante = new(int32)
	ret.SmallBlind = new(int32)
	return ret
}


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