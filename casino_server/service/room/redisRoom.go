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

