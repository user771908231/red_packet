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

func getRedisThUsrKey(deskId int32,gameNumber int32,userId uint32) string{
	deskIdStr,_ := numUtils.Int2String(deskId)
	gameNumberStr,_ := numUtils.Int2String(gameNumber)
	userIdStr,_ := numUtils.Uint2String(userId)

	result := strings.Join([]string{deskIdStr,gameNumberStr,userIdStr},"_")
	log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]",deskId,gameNumber,userId,result)
	return result
}

//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThUser(deskId int32,gameNumber int32,userId uint32) *bbproto.ThServerUser{
	key := getRedisThUsrKey(deskId,gameNumber,userId)

	//获取redis的连接
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//开始获取obj
	rthuser := &bbproto.ThServerUser{}
	err := conn.GetObj(key,rthuser)
	if err != nil {
		log.E("获取用户数据失败with key [%v]",key)
		return nil
	}

	return rthuser
}

//新建立一个user放在redis中
func NewRedisThuser(desk *ThDesk,userId uint32){
	user := NewThServerUser()
	//初始化值
	saveRedisThUser(user)
}

//返回一个初始化的user
func NewThServerUser()*bbproto.ThServerUser{
	user := &bbproto.ThServerUser{}
	user.Seat 	= new(int32)
	user.Status	= new(int32)
	user.BreakStatus= new(int32)
	user.WaiTime 	= new(string)
	user.WaitUUID	= new(string)
	user.DeskId 	= new(int32)
	user.TotalBet	= new(int64)
	user.TotalBet4CalcAllin	 = new(int64)
	user.WinAmount	= new(int64)
	user.TurnCoin	= new(int64)
	user.HandCoin	= new(int64)
	user.RoomCoin	= new(int64)
	user.UserId	= new(uint32)
	user.GameNumber = new(int32)

	return user
}

//保存一个用户
func saveRedisThUser(user *bbproto.ThServerUser) error{
	//获取redis连接
	key := getRedisThUsrKey(user.GetDeskId(),user.GetGameNumber(),user.GetUserId())
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//保存数据
	conn.SetObj(key,user)
	return nil
}

//删除一个用户
func DelRedisThUser(deskId int32,gameNumber int32,userId uint32) error{

	//得到key
	key := getRedisThUsrKey(deskId,gameNumber,userId)

	//获取redis连接
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//删除数据
	conn.Del(key)
	return nil

}

func UpdateRedisThuser(deskId int32,gameNumber int32,userId uint32,roomCoin,handCoin,turnCoin,totalBet,totalBet4CalcAllin,winAmount int64) error{
	//1,得到user
	user := GetRedisThUser(deskId,gameNumber,userId)
	if user == nil {
		return errors.New("没有找到用户")
	}

	//2,增加金额
	*user.RoomCoin = roomCoin
	*user.HandCoin = handCoin
	*user.TurnCoin = turnCoin
	*user.TotalBet = totalBet
	*user.TotalBet4CalcAllin = totalBet4CalcAllin
	*user.WinAmount = winAmount

	saveRedisThUser(user)
	return nil
}



//减少用户的金额

//关于thdesk的redis存储
var REDIS_TH_DESK_KEY_PRE = "redis_th_desk_key"
