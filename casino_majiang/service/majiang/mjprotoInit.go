package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
)

func NewMjpai() *MJPai {
	ret := &MJPai{}
	ret.Value = new(int32)
	ret.Des = new(string)
	ret.Flower = new(int32)
	ret.Index = new(int32)
	return ret
}

//返回一个麻将room
func NewMjRoom() *MjRoom {
	ret := &MjRoom{}
	return ret
}

//返回一个麻将
func NewMjDesk() *MjDesk {
	ret := &MjDesk{}
	ret.DeskId = new(int32)
	ret.RoomId = new(int32)
	ret.Status = new(int32)
	ret.Password = new(string)
	ret.Owner = new(uint32)
	ret.CreateFee = new(int64)
	ret.MjRoomType = new(int32)
	ret.BoardsCout = new(int32)
	ret.CapMax = new(int64)
	ret.CardsNum = new(int32)
	ret.Settlement = new(int32)
	ret.BaseValue = new(int64)
	ret.ZiMoRadio = new(int32)
	ret.DianGangHuaRadio = new(int32)
	ret.MJPaiCursor = new(int32)
	ret.ActiveUser = new(uint32)
	ret.Banker = new(uint32)
	ret.GameNumber = new(int32)
	ret.TotalPlayCount = new(int32)
	ret.CurrPlayCount = new(int32)
	ret.ActUser = new(uint32)
	ret.ActType = new(int32)
	ret.BeginTime = new(string)
	ret.EndTime = new(string)
	return ret
}

func NewMjUser() *MjUser {
	ret := &MjUser{}
	ret.UserId = new(uint32)
	ret.Coin = new(int64)
	ret.Status = new(int32)
	ret.IsBreak = new(bool)
	ret.IsLeave = new(bool)
	ret.DeskId = new(int32)
	ret.RoomId = new(int32)
	ret.DingQue = new(bool)
	ret.Bill = NewBill()
	ret.Statisc = NewMjUserStatisc()
	ret.IsBanker = new(bool)
	ret.WaitTime = new(int64)
	return ret
}

func NewBill() *Bill {
	ret := &Bill{}
	ret.WinAmount = new(int64)
	return ret
}

func NewMjUserStatisc() *MjUserStatisc {
	ret := &MjUserStatisc{}
	ret.WinCoin = new(int64)
	ret.CountZiMo = new(int32)
	ret.CountHu = new(int32)
	ret.CountDianPao = new(int32)
	ret.CountDianGang = new(int32)
	ret.CountAnGang = new(int32)
	ret.CountMingGang = new(int32)
	ret.CountChaJiao = new(int32)
	return ret
}

func NewBillBean() *BillBean {
	ret := &BillBean{}
	ret.UserId = new(uint32)
	ret.Des = new(string)
	ret.OutUserId = new(uint32)
	ret.Type = new(int32)
	ret.Amount = new(int64)
	return ret
}

func NewMjSession() *MjSession {
	s := &MjSession{}
	s.DeskId = new(int32)
	s.GameStatus = new(int32)
	s.RoomId = new(int32)
	s.UserId = new(uint32)
	return s
}

func NewMJHandPai() *MJHandPai {
	ret := &MJHandPai{}
	ret.QueFlower = new(int32)
	return ret
}

func NewPlayerGameData() *PlayerGameData {
	ret := &PlayerGameData{}
	ret.HandPai = NewMJHandPai()
	return ret

}

func NewGangPaiInfo() *GangPaiInfo {
	ret := &GangPaiInfo{}
	ret.ByWho = new(int32)
	ret.GangType = new(int32)
	ret.SendUserId = new(uint32)
	ret.GetUserId = new(uint32)
	return ret
}

func NewCheckCase() *CheckCase {
	ret := &CheckCase{}
	ret.UserIdOut = new(uint32)
	ret.CheckStatus = new(int32)
	return ret
}

func NewCheckBean() *CheckBean {
	ret := &CheckBean{}
	ret.UserId = new(uint32)
	ret.CanGang = new(bool)
	ret.CanHu = new(bool)
	ret.CanPeng = new(bool)
	ret.CheckStatus = new(int32)
	return ret
}

func NewHuPaiInfo() *HuPaiInfo {
	ret := &HuPaiInfo{}
	ret.ByWho = new(int32)
	ret.CardType = new(int32)
	ret.Fan = new(int32)
	ret.Score = new(int64)
	ret.HuType = new(int32)
	ret.GetUserId = new(uint32)
	ret.SendUserId = new(uint32)
	ret.HuDesc = new(string)
	return ret

}

//生成一张只有背面的牌
func NewBackPai() *mjproto.CardInfo {
	cardInfo := newProto.NewCardInfo()
	*cardInfo.Id = 0
	*cardInfo.Type = 0
	*cardInfo.Value = 0
	return cardInfo
}

func NewStatiscRound() *StatiscRound {
	ret := &StatiscRound{}
	ret.CountAnGang = new(int32)
	ret.GameNumber = new(int32)
	ret.Result = new(string)
	ret.WinAmount = new(int64)
	ret.CountHu = new(int32)
	ret.CountZiMo = new(int32)
	ret.CountDianPao = new(int32)
	ret.CountAnGang = new(int32)
	ret.CountMingGang = new(int32)
	ret.CountDianGang = new(int32)
	ret.CountChaJiao = new(int32)
	ret.CountChaJiao = new(int32)
	ret.CountBaGnag = new(int32)
	return ret
}