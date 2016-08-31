package userService

import (
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbprotogo"
	"casino_server/common/config"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/mode"
	"time"
	"casino_server/utils/db"
	"gopkg.in/mgo.v2"
	"casino_server/common/Error"
	"casino_server/utils/redisUtils"
)

var NEW_USER_DIAMOND_REWARD int64 = 20		//新用户登陆的时候,默认的砖石数量


/**
	1,create 一个user
	2,保存mongo
	3,缓存到redis
 */
func NewUserAndSave(openId, wxNickName, headUrl string) (*bbproto.User, error) {

	//1,创建user获得自增主键
	id, err := db.GetNextSeq(casinoConf.DBT_T_USER)
	if err != nil {
		return nil, err
	}
	//构造user
	nuser := &mode.T_user{}
	nuser.Mid = bson.NewObjectId()
	nuser.Id = uint32(id)
	nuser.Diamond = NEW_USER_DIAMOND_REWARD		//新用户注册的时候,默认的钻石数量
	if wxNickName == "" {
		nuser.NickName = config.RandNickname()
	} else {
		nuser.NickName = wxNickName
	}
	nuser.OpenId = openId
	nuser.HeadUrl = headUrl

	//2保存数据到数据库
	err = db.InsertMgoData(casinoConf.DBT_T_USER, nuser)
	if err != nil {
		log.E("保存用户的时候失败 error【%v】", err.Error())
		return nil, err
	}

	result, _ := Tuser2Ruser(nuser)
	return result, nil
}

func GetRedisUserKey(id uint32) string {
	idStr, _ := numUtils.Uint2String(id)
	return strings.Join([]string{casinoConf.DBT_T_USER, idStr}, "-")
}


//取session的rediskey
func GetRedisUserSeesionKey(userid uint32) string {
	idStr, _ := numUtils.Uint2String(userid)
	return strings.Join([]string{"agent_session", idStr}, "_")
}
/**
	根据用户id得到User的id
	1,首先从redis中查询user信息
	2,如果redis中不存在,则从mongo中查询
	3,如果mongo不存在,返回错误信息,客户端跳转到登陆界面

 */
func GetUserById(id uint32) *bbproto.User {
	//1,首先在 redis中去的数据
	key := GetRedisUserKey(id)
	result := redisUtils.GetObj(key, &bbproto.User{})
	if result == nil {
		log.E("redis中没有找到user[%v],需要在mongo中查询,并且缓存在redis中。", id)
		// 获取连接 connection
		tuser := &mode.T_user{}
		db.Query(func(d *mgo.Database) {
			d.C(casinoConf.DBT_T_USER).Find(bson.M{"id": id}).One(tuser)
		})

		if tuser.Id < casinoConf.MIN_USER_ID {
			log.T("在mongo中没有查询到user[%v].", id)
			result = nil
		} else {
			log.T("在mongo中查询到了user[%v],现在开始缓存", tuser)
			//把从数据获得的结果填充到redis的model中
			result, _ = Tuser2Ruser(tuser)
			if result != nil {
				SaveUser2Redis(result.(*bbproto.User))
			}
		}
	}

	//判断用户是否存在,如果不存在,则返回空
	if result == nil {
		return nil
	} else {
		ret := result.(*bbproto.User)
		ret.OninitLoginTurntableState()        //初始化登录转盘之后的奖励
		return ret
	}
}

//返回session信息
func GetUserSessionByUserId(id uint32) *bbproto.ThServerUserSession {
	result := redisUtils.GetObj(GetRedisUserSeesionKey(id), &bbproto.ThServerUserSession{})
	if result == nil {
		log.T("没有找到user[%v]的thserverUserSession", id)
		return nil
	} else {
		return result.(*bbproto.ThServerUserSession)
	}
}

//保存回话信息
func SaveUserSession(userData *bbproto.ThServerUserSession) {
	redisUtils.SaveObj(GetRedisUserSeesionKey(userData.GetUserId()), userData)
}

func GetUserByOpenId(openId  string) *bbproto.User {
	//1,首先在 redis中去的数据--登录考虑是否需要从redis中查询

	//2,从数据库中查询
	result := &bbproto.User{}
	tuser := &mode.T_user{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_USER).Find(bson.M{"openid": openId}).One(tuser)
	})

	if tuser == nil || tuser.Id < casinoConf.MIN_USER_ID {
		log.T("在mongo中没有查询到user[%v].", openId)
		result = nil
	} else {
		log.T("在mongo中查询到了user[%v],现在开始缓存", tuser)
		//把从数据获得的结果填充到redis的model中
		result, _ = Tuser2Ruser(tuser)
		if result != nil {
			SaveUser2Redis(result)
		}
	}

	//判断用户是否存在,如果不存在,则返回空
	if result == nil {
		return nil
	} else {
		//result.OninitLoginTurntableState()	//初始化登录转盘之后的奖励
		return result
	}
}


/**
	将用户model保存在redis中
 */
func SaveUser2Redis(u *bbproto.User) {
	redisUtils.SaveObj(GetRedisUserKey(u.GetId()), u)
}

/**
	保存数据到redis和mongo中
 */
func SaveUser2RedisAndMongo(u *bbproto.User) {
	SaveUser2Redis(u)
	UpsertRUser2Mongo(u)
}


//把redis中的数据刷新到数据库
func FlashUser2Mongo(userId uint32) error {
	u := GetUserById(userId)
	UpsertRUser2Mongo(u)
	return nil
}

func UpsertRUser2Mongo(u *bbproto.User) {
	//把bbproto.User转化为  model.User
	tuser, _ := Ruser2Tuser(u)        //
	UpsertTUser2Mongo(tuser)
}

//保存用户到mongo
func UpsertTUser2Mongo(tuser *mode.T_user) {
	//得到数据库连接池
	if tuser.Mid == "" {
		db.InsertMgoData(casinoConf.DBT_T_USER, tuser)
	} else {
		db.UpdateMgoData(casinoConf.DBT_T_USER, tuser)
	}
}

/**
	mongo中User模型转化为 redis中的user模型
 */
func Tuser2Ruser(tu *mode.T_user) (*bbproto.User, error) {
	result := &bbproto.User{}
	if tu.Mid.Hex() != "" {
		hesStr := tu.Mid.Hex()
		result.Mid = &hesStr
		//log.T("获得t_user.mid %v",hesStr)
	}

	result.Id = &tu.Id
	result.NickName = &tu.NickName
	result.Coin = &tu.Coin
	result.Diamond = &tu.Diamond
	result.OpenId = &tu.OpenId
	result.HeadUrl = &tu.HeadUrl
	return result, nil
}

/**
	redis中的user模型转化为mongdo的User模型
	把Redis_user 转化为mongo_t_user的时候喂自动为其分配objectId,方存储
 */

func Ruser2Tuser(ru *bbproto.User) (*mode.T_user, error) {
	result := &mode.T_user{}

	if ru.Mid != nil {
		result.Mid = bson.ObjectIdHex(ru.GetMid())
	} else {
		result.Mid = bson.NewObjectId()
	}

	result.Id = ru.GetId()
	result.NickName = ru.GetNickName()
	result.Coin = ru.GetCoin()
	result.HeadUrl = ru.GetHeadUrl()
	result.OpenId = ru.GetOpenId()
	result.Diamond = ru.GetDiamond()

	return result, nil
}

/**
	判断用户id是否合法,todo 这里是否在数据库中判断?由于之后的查询会在数据库中查询,所以这里可以先不用判断
	if userId  > casinoConf.MAX_USER_ID || userId < casinoConf.MIN_USER_ID {

 */
func CheckUserIdRightful(userId uint32) bool {
	u := GetUserById(userId)
	if u == nil {
		return false
	} else {
		return true
	}
}



//更新用户的钻石之后,在放回用户当前的余额,更新用户钻石需要同事更新redis和mongo的数据
func UpdateUserDiamond(userId uint32, diamond int64) (int64, error) {
	//1,获取锁
	//lock := UserLockPools.GetUserLockByUserId(userId)
	//lock.Lock()
	//defer lock.Unlock()

	//2,修改用户redis和mongo中的数据
	user := GetUserById(userId)
	if user == nil {
		return -1, Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_CREATE_DESK_USER_NOTFOUND), "用户不存在")
	}


	//见转砖石的时候,如果user.getDianmond + diamond < 0 表示余额不足
	//判断用户的钻石是否足够
	if user.GetDiamond() + diamond <= 0 {
		return user.GetDiamond(), Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_CREATE_DESK_DIAMOND_NOTENOUGH), "余额不足")
	}

	//修改并且更新用户数据
	*user.Diamond += diamond
	SaveUser2RedisAndMongo(user)

	//3,返回数据
	return user.GetDiamond(), nil
}

//craete钻石交易记录

func CreateDiamonDetail(userId uint32, detailsType int32, diamond int64, remainDiamond int64, memo string) error {

	//1,获得的交易记录自增主键
	id, err := db.GetNextSeq(casinoConf.DBT_T_USER_DIAMOND_DETAILS)
	if err != nil {
		return Error.NewError(0, err.Error())
	}

	//2,构造交易记录
	detail := &mode.T_user_diamond_details{}
	detail.Id = uint32(id)
	detail.UserId = userId
	detail.Diamond = diamond
	detail.ReaminDiamond = remainDiamond
	detail.DetailsType = detailsType
	detail.DetailTime = time.Now()
	detail.Memo = memo

	//3,保存数据
	err = db.InsertMgoData(casinoConf.DBT_T_USER_DIAMOND_DETAILS, detail)
	if err != nil {
		log.E("保存用户交易记录的时候失败 error【%v】", err.Error())
		return Error.NewError(0, "创建交易记录失败")
	}
	return nil
}


