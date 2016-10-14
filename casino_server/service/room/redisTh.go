package room

import (
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
	"casino_server/service/userService"
	"errors"
)

//关于thuser的redis存储
var REDIS_TH_USER_KEY_PRE = "redis_th_user_key"

func getRedisThUserKey(deskId int32, userId uint32) string {
	deskIdStr, _ := numUtils.Int2String(deskId)
	userIdStr, _ := numUtils.Uint2String(userId)
	result := strings.Join([]string{REDIS_TH_USER_KEY_PRE, deskIdStr, userIdStr}, "_")
	//log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]", deskId, gameNumber, userId, result)
	return result
}

//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThUser(deskId int32, userId uint32) *bbproto.ThServerUser {
	key := getRedisThUserKey(deskId, userId)
	p := redisUtils.GetObj(key, &bbproto.ThServerUser{})
	if p == nil {
		log.E("获取用户数据失败with key [%v]", key)
		return nil
	} else {
		return p.(*bbproto.ThServerUser)
	}

}


//把redis的user转化为thuser
func RedisThuserTransThuser(b *bbproto.ThServerUser) *ThUser {
	user := userService.GetUserById(b.GetUserId())
	t := NewThUser()
	t.UserId = b.GetUserId()      //用户id
	t.NickName = user.GetNickName()                //用户昵称
	t.deskId = b.GetDeskId()                 //用户所在的桌子的编号
	t.Status = b.GetStatus()                 //当前的状态,单局游戏的状态
	t.CSGamingStatus = b.GetCSGamingStatus()                  //是否进行锦标赛,这个字段其实 是在服务器crash之后,恢复数据的时候可以用到
	t.GameStatus = b.GetGameStatus()                 //用户的游戏状态
	t.IsBreak = b.GetIsBreak()                  //用户断线的状态,这里判断用户是否断线
	t.IsLeave = b.GetIsLeave()                  //用户是否处于离开的状态
	//log.T("t[%v],t.handcards[%v],b[%v],b.handcards[%v]",t==nil,t.HandCards==nil,b==nil,b.HandCards==nil)
	t.HandCards = b.GetHandCards()       //手牌
	//thCards           //恢复的时候初始化一次就行了
	//waiTime            time.Time             //等待时间
	//waitUUID           string                //等待标志
	t.PreCoin = b.GetPreCoin()                 //前注
	t.TotalBet = b.GetTotalBet()                 //计算用户总共押注的多少钱
	t.TotalBet4calcAllin = b.GetTotalBet4CalcAllin()                 //押注总额 ***注意,目前这个值是用来计算all in 的
	t.winAmount = b.GetWinAmount()                 //总共赢了多少钱
	t.winAmountDetail = b.GetWinAmountDetail()               //赢钱的细节, 主要是每个记录每个奖池赢了多少钱
	t.TurnCoin = b.GetTurnCoin()                 //单轮押注(总共四轮)的金额
	t.HandCoin = b.GetHandCoin()                 //用户下注多少钱、指单局
	t.RoomCoin = b.GetRoomCoin()                 //用户上分的金额
	t.RebuyCount = b.GetRebuyCount()              //重购的次数
	t.LotteryCheck = b.GetLotteryCheck()                  //这个字段用于判断是否可以开奖,默认是false:   1,如果用户操作弃牌,则直接设置为true,2,如果本局是all in,那么要到本轮次押注完成之后,才能设置为true
	t.TotalRoomCoin = b.GetTotalRoomCoin()
	t.IsShowCard = b.GetIsShowCard()
	t.MatchId = b.GetMatchId()                //得到matchId
	t.CloseCheck = b.GetCloseCheck()
	t.IsAutoReady = b.GetIsAutoReady()
	t.WaitRebuyFlag = b.GetWaitRebuyFlag()
	return t
}

//保存一个用户
func saveRedisThUser(user *bbproto.ThServerUser) error {
	//获取redis连接
	key := getRedisThUserKey(user.GetDeskId(), user.GetUserId())
	redisUtils.SetObj(key, user)
	return nil
}

//删除一个用户
func DelRedisThUser(deskId int32, userId uint32) {
	redisUtils.Del(getRedisThUserKey(deskId, userId))
}

//更新thuser的数据到redis中
func UpdateRedisThuser(u *ThUser) error {
	desk := u.GetDesk()
	var gameNumber int32 = 0
	if desk != nil {
		gameNumber = desk.GameNumber
	}
	//1,得到user
	ruser := GetRedisThUser(u.deskId, u.UserId)
	if ruser == nil {
		log.T("数据库中没有找到thuser【%v】", u.UserId)
		ruser = bbproto.NewThServerUser()
	}

	//2,为user赋值
	*ruser.UserId = u.UserId
	*ruser.DeskId = u.deskId
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
	*ruser.GameNumber = gameNumber
	*ruser.TotalRoomCoin = u.TotalRoomCoin
	*ruser.LotteryCheck = u.LotteryCheck
	*ruser.IsShowCard = u.IsShowCard
	*ruser.MatchId = u.MatchId
	*ruser.CloseCheck = u.CloseCheck
	*ruser.IsAutoReady = u.IsAutoReady
	*ruser.WaitRebuyFlag = u.WaitRebuyFlag

	//3,保存到数据库
	saveRedisThUser(ruser)
	return nil
}



//---------------------------------------------------------------------------thdesk----------
var REDIS_TH_DESK_KEY_PRE = "redis_th_desk_key"

func getRedisThDeskKey(deskId int32) string {
	deskIdStr, _ := numUtils.Int2String(deskId)
	result := strings.Join([]string{REDIS_TH_DESK_KEY_PRE, deskIdStr}, "_")
	//log.T("通过deskId[%v],gameNumber[%v],userId[%v]得到的redis_key是[%v]", deskId, gameNumber, userId, result)
	return result
}


//通过桌子id,游戏编号,用户id来唯一确定一个thuser
func GetRedisThDesk(deskId int32) *bbproto.ThServerDesk {
	key := getRedisThDeskKey(deskId)
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
	ret := NewThDesk()
	ret.Id = rt.GetId()
	ret.MatchId = rt.GetMatchId()
	ret.DeskOwner = rt.GetDeskOwner()
	ret.RoomKey = rt.GetRoomKey()
	ret.CreateFee = rt.GetCreateFee()
	ret.GameType = rt.GetGameType()
	ret.InitRoomCoin = rt.GetInitRoomCoin()
	ret.JuCount = rt.GetJuCount()
	ret.JuCountNow = rt.GetJuCountNow()
	ret.PreCoin = rt.GetPreCoin()
	ret.SmallBlindCoin = rt.GetSmallBlindCoin()
	ret.BigBlindCoin = rt.GetBigBlindCoin()
	ret.blindLevel = rt.GetBlindLevel()
	//BeginTime            time.Time                    //游戏开始时间
	//EndTime              time.Time                    //游戏结束时间
	ret.RebuyCountLimit = rt.GetRebuyCountLimit()                        //重购的次数
	ret.RebuyBlindLevelLimit = rt.GetRebuyBlindLevelLimit()                        //rebuy盲注的限制

	ret.Dealer = rt.GetDealer()                       //庄家
	ret.BigBlind = rt.GetBigBlind()                       //大盲注
	ret.SmallBlind = rt.GetSmallBlind()                       //小盲注
	ret.RaiseUserId = rt.GetRaiseUserId()                       //加注的人的Id,一轮结束的判断需要按照这个人为准
	ret.NewRoundFirstBetUser = rt.GetNewRoundFirstBetUser()                       //新一轮,开始押注的第一个人//第一轮默认是小盲注,但是当小盲注弃牌之后,这个人要滑倒下一家去
	ret.BetUserNow = rt.GetBetUserNow()                       //当前押注人的Id

	ret.GameNumber = rt.GetGameNumber()                        //每一局游戏的游戏编号
	ret.PublicPai = rt.GetPublicPai()               //公共牌的部分
	ret.Status = rt.GetStatus()                        //牌桌的状态
	ret.BetAmountNow = rt.GetBetAmountNow()                        //当前的押注金额是多少
	ret.RoundCount = rt.GetRoundCount()                        //第几轮
	ret.Jackpot = rt.GetJackpot()                        //奖金池
	ret.EdgeJackpot = rt.GetEdgeJackpot()                        //边池
	ret.MinRaise = rt.GetMinRaise()                        //最低加注金额
	//ret.AllInJackpot         []*pokerService.AllInJackpot //allin的标记

	ret.SendFlop = rt.GetSendFlop()                         //公共底牌
	ret.SendTurn = rt.GetSendTurn()                         //第四章牌
	ret.SendRive = rt.GetSendRive()                         //第五章牌
	return ret
}


//更新thdesk的数据到redis中
func UpdateTedisThDesk(t *ThDesk) error {

	//1,检测参数
	if t == nil {
		log.E("需要备份数据到redis 的thdesk为nil,备份失败...")
		return errors.New("备份desk到redis失败...")
	}

	rt := GetRedisThDesk(t.Id)
	if rt == nil {
		log.T("数据库中没有找到thdesk【%v】,重新生成一个,并且保存到redis中", t.Id)
		rt = bbproto.NewThServerDesk()
	}
	//2,为user赋值
	*rt.Id = t.Id
	*rt.MatchId = t.MatchId
	*rt.DeskOwner = t.DeskOwner
	*rt.RoomKey = t.RoomKey
	*rt.CreateFee = t.CreateFee
	*rt.DeskType = t.GameType
	*rt.GameType = t.GameType
	*rt.InitRoomCoin = t.InitRoomCoin
	*rt.JuCount = t.JuCount
	*rt.JuCountNow = t.JuCountNow
	*rt.PreCoin = t.PreCoin

	*rt.SmallBlindCoin = t.SmallBlindCoin
	*rt.BigBlindCoin = t.BigBlindCoin
	*rt.BlindLevel = t.blindLevel
	*rt.RebuyCountLimit = t.RebuyCountLimit
	*rt.RebuyBlindLevelLimit = t.RebuyBlindLevelLimit
	//BeginTime            time.Time
	//EndTime              time.Time
	*rt.Dealer = t.Dealer
	*rt.BigBlind = t.BigBlind
	*rt.SmallBlind = t.SmallBlind
	*rt.RaiseUserId = t.RaiseUserId
	*rt.NewRoundFirstBetUser = t.NewRoundFirstBetUser
	*rt.BetUserNow = t.BetUserNow
	*rt.GameNumber = t.GameNumber
	rt.PublicPai = t.PublicPai
	*rt.UserCount = t.GetUserCount()
	*rt.Status = t.Status
	*rt.BetAmountNow = t.BetAmountNow
	*rt.RoundCount = t.RoundCount
	*rt.Jackpot = t.Jackpot
	*rt.EdgeJackpot = t.EdgeJackpot
	*rt.MinRaise = t.MinRaise
	rt.AllInJackpot = t.GetServerProtoAllInJackPot()
	*rt.SendFlop = t.SendFlop
	*rt.SendTurn = t.SendTurn
	*rt.SendRive = t.SendRive

	//所有人的id
	rt.UserIds = t.GetuserIds()
	rt.LeaveUserIds = t.GetLeaveUserIds()

	//3,保存到数据库
	saveRedisThDesk(rt)
	return nil

}

//保存一个thdesk
func saveRedisThDesk(t *bbproto.ThServerDesk) error {
	//获取redis连接
	key := getRedisThDeskKey(t.GetId())
	redisUtils.SetObj(key, t)
	return nil
}

//删除redis中的thdesk
func DelRedisThdesk(deskId int32) error {
	key := getRedisThDeskKey(deskId)
	redisUtils.Del(key)
	return nil
}



