package userService

import (
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbproto"
	"fmt"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/mode"
	"casino_server/common/config"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
	"errors"
	"casino_server/conf/StrCons"
	"casino_server/conf/intCons"
)

/**

判断id是否正确
 */
func CheckUserId(userId uint32) int8{
	if userId > casinoConf.MAX_USER_ID || userId < casinoConf.MIN_USER_ID {
		return intCons.LOGIN_WAY_QUICK
	}else{
		return intCons.LOGIN_WAY_LOGIN
	}
}


/**
快速登录
	快速登录模式需要 为用户分配一个id,并且返回给用户
 */
func QuickLogin(user *bbproto.ReqAuthUser) (*mode.User,error){
	//1,判断入参是否正确
	uuid := user.Uuid
	log.Debug("header.code %v",uuid)
	if uuid == nil {
		log.Error("登录的时候uuid 为空,无法登陆")
	}
	//2,为用户分配id
	nuser,err := newUserAndSaveToDB()
	if err !=nil{
		log.E(err.Error())
		return nil,err
	}else{
		return nuser,nil
	}
}


/**
普通登录
 */
func Login(user *bbproto.ReqAuthUser) (*mode.User, error){
	//1,检测参数
	userId := user.Header.UserId
	log.Debug("需要登陆的userId %v",userId)

	//2,通过userId 在mongo 中查询user的信息
	dbUser := &mode.User{}
	//1,获取数据库连接和回话
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()
	s:=c.Ref()
	defer c.UnRef(s)

	err2 := s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"id":userId}).One(&dbUser)
	if err2 != nil{
		log.E(err2.Error())
		return nil,errors.New(StrCons.STR_ERR_LOGIN_NOT_FOUND_USER)
	}
	return dbUser,nil
}

/**
 动态生成一个 User,并且保存到数据库
 */
func newUserAndSaveToDB() (*mode.User,error){
	//1,获取数据库连接和回话
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	s:=c.Ref()
	defer c.UnRef(s)

	//2,创建user获得自增主键
	nuser := &mode.User{}
	id, err := c.NextSeq(casinoConf.DB_NAME, casinoConf.DBT_T_USER, casinoConf.DB_ENSURECOUNTER_KEY)
	nuser.Id = uint32(id)
	nuser.NickName = config.RandNikeName()
	nuser.Mid = bson.NewObjectId()

	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Insert(nuser)
	return  nuser,nil

}
