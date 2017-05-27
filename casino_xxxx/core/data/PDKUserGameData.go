package data

type UserOutCards []*PDKOutCard //打出去的牌
type UserHandCards []*PokerCard //手牌

type PDKUserGameData struct {
	//账单
	OutCards  UserOutCards  //outCard
	HandCards UserHandCards //handCard
}

//玩家的账单
type PDKUserBill struct {
}
