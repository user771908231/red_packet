package bbproto

//初始化盲注bean
func NewGame_TounamentBlindBean() *Game_TounamentBlindBean{
	ret := &Game_TounamentBlindBean{}
	ret.Ante = new(int32)
	ret.SmallBlind = new(int32)

	return ret
}


//返回一个奖励的bean
func NewGame_TounamentRewardsBean() *Game_TounamentRewardsBean{
	ret := &Game_TounamentRewardsBean{}
	ret.IconPath = new(string)
	ret.Rewards = new(string)
	return ret
}


//初始化一个rankBean
func NewGame_TounamentRankBean() *Game_TounamentRankBean{
	ret := &Game_TounamentRankBean{}
	ret.Coin = new(int64)
	ret.Place = new(int32)
	ret.PlayerImage = new(string)
	ret.PlayerName  = new(string)
	return ret
}

//初始化一个

func NewGame_MatchItem()*Game_MatchItem{
	ret := &Game_MatchItem{}
	ret.CostFee = new(int64)
	ret.Person  = new(int32)
	ret.Status  = new(int32)
	ret.Person  = new(int32)
	ret.Time    = new(string)
	ret.Title   = new(string)
	ret.Type    = new(int32)
	return ret
}