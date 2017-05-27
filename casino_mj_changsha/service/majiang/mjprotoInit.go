package majiang

import (
	mjproto	"casino_mj_changsha/msg/protogo"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/game"
	"github.com/name5566/leaf/module"
	"casino_common/common/service/countService"
	"github.com/golang/protobuf/proto"
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
func NewMjRoom(s *module.Skeleton) *MjRoom {
	ret := new(MjRoom)
	ret.PMjRoom = NewPMjRoom()
	ret.Skeleton = s
	return ret
}

func NewPMjRoom() *PMjRoom {
	ret := new(PMjRoom)
	ret.RoomType = new(int32)
	return ret
}

//返回一个麻将
func NewMjDesk() *MjDesk {
	ret := &MjDesk{}
	ret.PMjDesk = NewPMjDesk()
	ret.GameDesk = game.NewGameDesk()
	return ret
}

func NewPMjDesk() *PMjDesk {
	ret := new(PMjDesk)
	ret.RoomType = new(int32)
	ret.MjRoomType = new(int32)
	ret.CapMax = new(int64)
	ret.CardsNum = new(int32)
	ret.Settlement = new(int32)
	ret.ZiMoRadio = new(int32)
	ret.DianGangHuaRadio = new(int32)
	ret.MJPaiCursor = new(int32)
	ret.ActiveUser = new(uint32)
	ret.TotalPlayCount = new(int32)
	ret.CurrPlayCount = new(int32)
	ret.ActUser = new(uint32)
	ret.ActType = new(int32)
	ret.HuRadio = new(int32)
	ret.NextBanker = new(uint32)
	ret.FangCountLimit = new(int32)
	ret.RoomLevel = new(int32)
	return ret
}

func NewMjUser() *MjUser {
	ret := new(MjUser)
	ret.PMjUser = NewPMjUser()
	ret.GameUser = game.NewGameUser()
	ret.Log = new(countService.T_game_log)
	return ret
}

func NewPMjUser() *PMjUser {
	ret := new(PMjUser)
	ret.DingQue = new(bool)
	ret.Bill = NewBill()
	ret.Statisc = NewMjUserStatisc()
	ret.IsBanker = new(bool)
	ret.Exchanged = new(bool)
	ret.Ready = new(bool)
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
	ret.CountBaGang = new(int32)
	ret.CountAnGang = new(int32)
	ret.CountMingGang = new(int32)
	ret.CountChaDaJiao = new(int32)
	ret.CountBeiAnGang = new(int32)
	ret.CountBeiZiMo = new(int32)
	ret.CountBeiBaGang = new(int32)
	ret.CountBeiChaHuaZhu = new(int32)
	ret.CountBeiChaJiao = new(int32)
	ret.CountChaHuaZhu = new(int32)
	ret.CountCatchBird = new(int32)
	ret.CountCaughtBird = new(int32)
	return ret
}

func NewBillBean() *BillBean {
	ret := &BillBean{}
	ret.UserId = new(uint32)
	ret.Des = new(string)
	ret.OutUserId = new(uint32)
	ret.Type = new(int32)
	ret.Amount = new(int64)
	ret.IsBird = new(bool)
	return ret
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
	ret.CanChi = new(bool)
	ret.CanGuo = proto.Bool(true) //默认可以过，特殊情况不能过，比如在有些规则里面，最后n张能胡牌必须胡牌
	return ret
}

func NewHuPaiInfo() *HuPaiInfo {
	ret := &HuPaiInfo{}
	ret.ByWho = new(int32)
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
	ret.CountBaGnag = new(int32)
	ret.CountBeiAnGang = new(int32)
	ret.CountBeiBaGang = new(int32)
	ret.CountBeiZiMo = new(int32)
	ret.CountBeiChaJiao = new(int32)
	ret.CountBeiChaHuaZhu = new(int32)
	ret.CountChaDaJiao = new(int32)
	ret.CountChaHuaZhu = new(int32)
	ret.CountCatchBird = new(int32)
	ret.CountCaughtBird = new(int32)
	ret.Round = new(int32)
	return ret
}

//生成一个过胡的info
func NewGuoHuInfo() *GuoHuInfo {
	ret := &GuoHuInfo{}
	ret.FanShu = new(int32)
	ret.GetUserId = new(uint32)
	ret.SendUserId = new(uint32)
	return ret
}

//生成一个叫的info
func NewJiaoInfo() *mjproto.JiaoInfo {
	ret := &mjproto.JiaoInfo{}
	ret.OutCard = newProto.NewCardInfo()
	ret.PaiInfos = nil
	return ret
}

//生成一个叫牌的info
func NewJiaoPaiInfo() *mjproto.JiaoPaiInfo {
	ret := &mjproto.JiaoPaiInfo{}
	ret.HuCard = newProto.NewCardInfo()
	ret.Count = new(int32)
	ret.Fan = new(int32)
	return ret
}
