package loginBonus

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/userService"
	"errors"
	"casino_server/utils"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/conf/casinoConf"
	"fmt"
	"casino_server/common/log"
	"casino_server/utils/timeUtils"
)

//转盘的格子
//转盘的奖励的中奖概率设置

var LoginTurntableBonus []int32 = []int32{1,2,3,4,5,6,7,8,9,10,11}
var LoginTurntablePro []int32 = []int32{0,10,20,30,40,50,60,70,80,90,100,110}
var LoginTurntableProMax = 200
var LoginTurntableCount  = 11		//转盘的格子最大数目


/**
	转盘奖励:
	1,每日只有一次
 */
func HandleLoginTurntableBonus(m *bbproto.LoginTurntableBonus,a gate.Agent) error{

	//返回值
	result := &bbproto.LoginTurntableBonus{}

	//1,检测参数是否正确
	err := checkBonusAble(m.GetHeader().GetUserId())
	if err != nil {
		var errMsg string = err.Error()
		result.Header = protoUtils.GetErrorHeaderWithMsg(&errMsg)
		a.WriteMsg(result)
		return err
	}

	//2,开始发放奖励
	var si int32 = 0
	pro := utils.Randn(LoginTurntableProMax)		//的到的概率
	for i := 0; i<LoginTurntableCount;i++  {
		if pro < LoginTurntablePro[i + 1] {
			si = int32(i)
			break
		}
	}
	var ba int32 = int32(LoginTurntableBonus[si])

	//计算奖励之后,保存到数据库
	updateTurntableBonus(m.GetHeader().GetUserId(),ba)

	//3,返回结果
	result.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
	result.BonusAmount = &ba
	result.SuccIndex = &si
	a.WriteMsg(result)		//给客户端发送成功的消息

	return nil
}



/**
	检测用户是否可以转动转盘
 */

func checkBonusAble(userId uint32) (error){

	//1,判断用户的ID是否合法
	if !userService.CheckUserIdRightful(userId) {
		//需要给客户端返回错误信息
		return errors.New("用户的ID不合法")
	}


	//2,判断今日是否可以再次领取
	user := userService.GetUserById(userId)
	if user.GetLoginTurntable() {
		return nil
	}else{
		return errors.New("已经领取过奖了")
	}
}

/**
	保存数据到数据库,并且更新用户的状态
 */
func updateTurntableBonus(userId uint32,amount int32) error{
	//1,对更新操作加锁,判断用户是否正确
	lock := userService.UserLockPools.GetUserLockByUserId(userId)
	lock.Lock()
	defer lock.Unlock()

	user := userService.GetUserById(userId)
	if user.LoginTurntable == false{
		log.E("用户[%v]领取转盘奖励失败,因为今天已经领取过了...",userId)
		return errors.New("领取失败,今日奖励已经领取过了")
	}
	//2,保存转盘记录

	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	s := c.Ref()
	defer c.UnRef(s)

	//2,创建user获得自增主键
	tid, err := c.NextSeq(casinoConf.DB_NAME, casinoConf.DBT_T_BONUS_TURNTABLE, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return nil, err
	}
	t := &mode.T_bonus_turntable{}
	t.Mid = bson.NewObjectId()
	t.Amount = amount
	t.Time = time.Now()
	t.UserId = userId
	t.Id =tid
	err = s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_BONUS_TURNTABLE).Insert(t)
	if err != nil {
		log.E("给用户发送转盘奖励的时候,出错:",err.Error())
		return err
	}
	//3,更新用户的记录
	var coinResult = user.GetCoin() + amount

	user.LoginTurntable = false;	//表示已经领取
	user.LoginTurntableTime = timeUtils.FormatYYYYMMDD(time.Now())
	user.Coin = &coinResult
	userService.SaveUser2Redis(user)

	return nil
}