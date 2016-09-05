package room

import (
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
)

//关于thuser的redis存储
var REDIS_TH_USER_KEY_PRE = "redis_th_user_key"

func getRedisThUserKey(deskId int32, gameNumber int32, userId uint32) string {
	deskIdStr, _ := numUtils.Int2String(deskId)
	gameNumberStr, _ := numUtils.Int2String(gameNumber)
	userIdStr, _ := numUtils.Uint2String(userId)

	result := strings.Join([]string{REDIS_TH_USER_KEY_PRE, deskIdStr, gameNumberStr, userIdStr}, "_")
	//log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]", deskId, gameNumber, userId, result)
	return result
}

//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThUser(deskId int32, gameNumber int32, userId uint32) *bbproto.ThServerUser {
	key := getRedisThUserKey(deskId, gameNumber, userId)
	p := redisUtils.GetObj(key, &bbproto.ThServerUser{})
	if p == nil {
		log.E("获取用户数据失败with key [%v]", key)
		return nil
	} else {
		return p.(*bbproto.ThServerUser)
	}

}

//保存一个用户
func saveRedisThUser(user *bbproto.ThServerUser) error {
	//获取redis连接
	key := getRedisThUserKey(user.GetDeskId(), user.GetGameNumber(), user.GetUserId())
	redisUtils.SaveObj(key, user)
	return nil
}

//删除一个用户
func DelRedisThUser(deskId int32, gameNumber int32, userId uint32) error {
	return nil
}

//更新thuser的数据到redis中
func UpdateRedisThuser(u *ThUser) error {
	//1,得到user
	ruser := GetRedisThUser(u.deskId, u.GameNumber, u.UserId)
	if ruser == nil {
		log.T("数据库中没有找到thuser【%v】", u.UserId)
		ruser = bbproto.NewThServerUser()
	}

	//2,为user赋值
	*ruser.UserId = u.UserId
	*ruser.DeskId = u.deskId
	*ruser.GameNumber = u.GameNumber
	*ruser.Seat = u.Seat
	*ruser.Status = u.Status
	*ruser.IsBreak = u.IsBreak
	*ruser.IsLeave = u.IsLeave
	ruser.HandCards = u.HandCards
	*ruser.WaitUUID = u.waitUUID
	*ruser.TotalBet = u.TotalBet
	*ruser.TotalBet4CalcAllin = u.TotalBet4calcAllin
	*ruser.WinAmount = u.winAmount
	ruser.WinAmountDetail = u.winAmountDetail
	*ruser.TurnCoin = u.TurnCoin
	*ruser.HandCoin = u.HandCoin
	*ruser.RoomCoin = u.RoomCoin
	*ruser.GameNumber = u.GameNumber

	//3,保存到数据库
	saveRedisThUser(ruser)
	return nil
}



//---------------------------------------------------------------------------thdesk----------
var REDIS_TH_DESK_KEY_PRE = "redis_th_desk_key"

func getRedisThDeskKey(deskId int32, gameNumber int32) string {
	deskIdStr, _ := numUtils.Int2String(deskId)
	gameNumberStr, _ := numUtils.Int2String(gameNumber)

	result := strings.Join([]string{REDIS_TH_DESK_KEY_PRE, deskIdStr, gameNumberStr}, "_")
	//log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]", deskId, gameNumber, userId, result)
	return result
}


//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThDesk(deskId int32, gameNumber int32) *bbproto.ThServerDesk {
	key := getRedisThDeskKey(deskId, gameNumber)
	return GetRedisThDeskByKey(key)
}

func GetRedisThDeskByKey(key string) *bbproto.ThServerDesk {
	p := redisUtils.GetObj(key, &bbproto.ThServerDesk{})
	if p == nil {
		log.E("获取用户数据失败with key [%v]", key)
		return nil
	} else {
		return p.(*bbproto.ThServerDesk)
	}
}

func RedisDeskTransThdesk(rt *bbproto.ThServerDesk) *ThDesk {
	return nil
}


//更新thdesk的数据到redis中
func UpdateTedisThDesk(t *ThDesk) error {
	rt := GetRedisThDesk(t.Id, t.GameNumber)
	if rt == nil {
		log.T("数据库中没有找到thdesk【%v】,重新生成一个,并且保存到redis中", t.Id)
		rt = bbproto.NewThServerDesk()
	}
	//2,为user赋值
	*rt.Id = t.Id
	*rt.DeskOwner = t.DeskOwner
	*rt.RoomKey = t.RoomKey
	*rt.DeskType = t.GameType
	*rt.InitRoomCoin = t.InitRoomCoin
	*rt.JuCount = t.JuCount
	*rt.SmallBlindCoin = t.SmallBlindCoin
	*rt.BigBlindCoin = t.BigBlindCoin
	*rt.Dealer = t.Dealer
	*rt.BigBlind = t.BigBlind
	*rt.SmallBlind = t.SmallBlind
	*rt.RaiseUserId = t.RaiseUserId
	*rt.NewRoundFirstBetUser = t.NewRoundFirstBetUser
	*rt.BetUserNow = t.BetUserNow
	*rt.GameNumber = t.GameNumber
	*rt.UserCount = t.UserCount
	*rt.Status = t.Status
	*rt.BetAmountNow = t.BetAmountNow
	*rt.RoundCount = t.RoundCount
	*rt.Jackpot = t.Jackpot
	*rt.EdgeJackpot = t.EdgeJackpot
	*rt.MinRaise = t.MinRaise

	rt.AllInJackpot = t.GetServerProtoAllInJackPot()
	rt.PublicPai = t.PublicPai

	//3,保存到数据库
	saveRedisThDesk(rt)
	return nil

}

//保存一个thdesk
func saveRedisThDesk(t *bbproto.ThServerDesk) error {
	//获取redis连接
	key := getRedisThDeskKey(t.GetId(), t.GetGameNumber())
	redisUtils.SaveObj(key, t)
	return nil
}




