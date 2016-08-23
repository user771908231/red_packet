package room

import (
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redis"
	"errors"
	"casino_server/conf/casinoConf"
)

//关于thuser的redis存储
var REDIS_TH_USER_KEY_PRE = "redis_th_user_key"

func getRedisThUserKey(deskId int32, gameNumber int32, userId uint32) string {
	deskIdStr, _ := numUtils.Int2String(deskId)
	gameNumberStr, _ := numUtils.Int2String(gameNumber)
	userIdStr, _ := numUtils.Uint2String(userId)

	result := strings.Join([]string{REDIS_TH_USER_KEY_PRE, deskIdStr, gameNumberStr, userIdStr}, "_")
	log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]", deskId, gameNumber, userId, result)
	return result
}

//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThUser(deskId int32, gameNumber int32, userId uint32) *bbproto.ThServerUser {
	key := getRedisThUserKey(deskId, gameNumber, userId)

	//获取redis的连接
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//开始获取obj
	rthuser := &bbproto.ThServerUser{}
	err := conn.GetObj(key, rthuser)
	if err != nil {
		log.E("获取用户数据失败with key [%v]", key)
		return nil
	}

	return rthuser
}

//新建立一个user放在redis中
func NewRedisThuser(user *ThUser) {

	//初始化值 user
	ruser := NewThServerUser()
	*ruser.UserId = user.UserId
	*ruser.DeskId = user.deskId
	*ruser.GameNumber = user.GameNumber
	*ruser.Seat = user.Seat
	*ruser.Status = user.Status
	*ruser.IsBreak = user.IsBreak
	*ruser.IsLeave = user.IsLeave
	ruser.HandCards = user.HandCards
	//*ruser.WaiTime = user.waiTime
	*ruser.WaitUUID = user.waitUUID
	*ruser.TotalBet = user.TotalBet
	*ruser.TotalBet4CalcAllin = user.TotalBet4calcAllin
	*ruser.WinAmount = user.winAmount
	ruser.WinAmountDetail = user.winAmountDetail
	*ruser.TurnCoin = user.TurnCoin
	*ruser.HandCoin = user.HandCoin
	*ruser.RoomCoin = user.RoomCoin

	//保存
	saveRedisThUser(ruser)
}

//返回一个初始化的user
func NewThServerUser() *bbproto.ThServerUser {
	user := &bbproto.ThServerUser{}
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
	return user
}

//保存一个用户
func saveRedisThUser(user *bbproto.ThServerUser) error {
	//获取redis连接
	key := getRedisThUserKey(user.GetDeskId(), user.GetGameNumber(), user.GetUserId())
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//保存数据
	conn.SetObj(key, user)
	return nil
}

//删除一个用户
func DelRedisThUser(deskId int32, gameNumber int32, userId uint32) error {

	//得到key
	key := getRedisThUserKey(deskId, gameNumber, userId)

	//获取redis连接
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//删除数据
	conn.Del(key)
	return nil

}

func UpdateRedisThuser(u *ThUser) error {
	//1,得到user
	user := GetRedisThUser(u.deskId, u.GameNumber, u.UserId)
	if user == nil {
		return errors.New("没有找到用户")
	}

	//2,增加金额
	*user.RoomCoin = u.RoomCoin
	*user.HandCoin = u.HandCoin
	*user.TurnCoin = u.TurnCoin
	*user.TotalBet =u.TotalBet
	*user.TotalBet4CalcAllin = u.TotalBet4calcAllin
	*user.WinAmount = u.winAmount
	*user.IsBreak = u.IsBreak
	*user.IsLeave = u.IsLeave
	*user.Status = u.Status
	saveRedisThUser(user)
	return nil
}



//减少用户的金额

//关于thdesk的redis存储
var REDIS_TH_DESK_KEY_PRE = "redis_th_desk_key"
