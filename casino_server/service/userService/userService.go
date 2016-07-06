package userService

import (
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbprotogo"
	"fmt"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/common/config"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
	"errors"
	"casino_server/conf/StrCons"
	"casino_server/conf/intCons"
	"casino_server/utils/redis"
	"casino_server/utils/numUtils"
	"strings"
)


//const-config
var(
	//更新的方式,是只更新到redis 还是redis和mongo都需要更新
	UPDATE_TYPE_ONLY_REDIS	int  = 1
	UPDATE_TYPE_REAIS_MONGO int  = 2

)


/**

判断id是否正确
 */
func CheckUserId(userId uint32) int8 {
	if userId > casinoConf.MAX_USER_ID || userId < casinoConf.MIN_USER_ID {
		return intCons.LOGIN_WAY_QUICK
	} else {
		return intCons.LOGIN_WAY_LOGIN
	}
}


/**
快速登录
	1,快速登录模式需要 为用户分配一个id,并且返回给用户
	2,登陆成功之后,需要为agent 绑定userData
 */
func QuickLogin(user *bbproto.ReqAuthUser) (*bbproto.User, error) {
	//1,判断入参是否正确
	uuid := user.Uuid
	log.Debug("header.code %v", uuid)
	if uuid == nil {
		log.Error("登录的时候uuid 为空,无法登陆")
	}
	//2,为用户分配id
	nuser, err := newUserAndSave()
	if err != nil {
		log.E(err.Error())
		return nil, err
	} else {
		return nuser, nil
	}
}


/**
普通登录
 */
func Login(user *bbproto.ReqAuthUser) (*bbproto.User, error) {
	//1,检测参数
	userId := user.Header.UserId
	log.Debug("需要登陆的userId %v", userId)

	//2,通过userId 在mongo 中查询user的信息
	dbUser := &bbproto.User{}
	//1,获取数据库连接和回话
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()
	s := c.Ref()
	defer c.UnRef(s)

	err2 := s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"id":userId}).One(&dbUser)
	if err2 != nil {
		log.E(err2.Error())
		return nil, errors.New(StrCons.STR_ERR_LOGIN_NOT_FOUND_USER)
	}
	return dbUser, nil
}

/**
	1,create 一个user
	2,保存mongo
	3,缓存到redis
 */
func newUserAndSave() (*bbproto.User, error) {
	//1,获取数据库连接和回话
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	s := c.Ref()
	defer c.UnRef(s)

	//2,创建user获得自增主键
	id, err := c.NextSeq(casinoConf.DB_NAME, casinoConf.DBT_T_USER, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return nil, err
	}
	userId := uint32(id)
	Nickname := config.RandNickname()

	//构造user
	nuser := &bbproto.User{}
	nuser.Id = &userId
	nuser.NickName = &Nickname
	nuser.Balance = &intCons.NUM_INT32_0
	err = s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Insert(nuser)
	if err != nil {
		return nil,err
	}
	return nuser, nil

}

func AddGold(userId uint32, ) {

}

func GetRedisUserKey(id uint32) string {
	idStr, _ := numUtils.Uint2String(id)
	return strings.Join([]string{idStr, casinoConf.DBT_T_USER}, "-")
}

/**
	根据用户id得到User的id
	1,首先从redis中查询user信息
	2,如果redis中不存在,则从mongo中查询
	3,如果mongo不存在,返回错误信息,客户端跳转到登陆界面

 */
func GetUserById(id uint32) *bbproto.User {
	//1,首先在 redis中去的数据
	log.T("在redis中查询user[%v]是否存在.", id)
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	key := GetRedisUserKey(id)
	result := &bbproto.User{}
	conn.GetObj(key, result)
	if result == nil || result.GetId() == 0 {

		log.E("redis中没有找到user[%v],需要在mongo中查询,并且缓存在redis中。", id)
		// 获取连接 connection
		c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
		if err != nil {
			result = nil
		}
		defer c.Close()
		s := c.Ref()
		defer c.UnRef(s)

		//从数据库中查询user
		user := &bbproto.User{}
		s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"id": id}).One(user)
		if user.GetId() < casinoConf.MIN_USER_ID {
			result = nil
		}
	}

	return result
}

/**
	将用户model保存在redis中
 */
func SaveUser2Redis(u *bbproto.User) {
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()
	key := GetRedisUserKey(u.GetId())
	conn.SetObj(key, u)
}


func UpsertUser2Mongo(u *bbproto.User){
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).UpsertId(bson.M{"id": u.GetId()},u)

}
/**
	更新用用户余额的信息
 */
func UpUserBalance(userId uint32, amount int32,utype int) error {

	//1,获得锁
	l := userLockPool.getUserLockByUserId(userId)
	l.Lock()
	defer l.Unlock()


	//2,跟新redis中的值
	//由于用户user相关的都会存在redis 中的,所以肯定会更新redis
	user := GetUserById(userId)
	var b int32 = user.GetBalance()
	b += amount
	user.Balance = &b
	SaveUser2Redis(user)	//保存user

	//3,更新mongo 中的值
	if utype == UPDATE_TYPE_REAIS_MONGO {
		UpsertUser2Mongo(user)
	}

	return nil
}
