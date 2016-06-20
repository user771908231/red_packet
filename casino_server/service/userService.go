package service

import (
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbproto"
	"github.com/name5566/leaf/gate"
	"fmt"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/mode"
	"casino_server/common/config"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
)

/**

判断id是否正确
 */
func CheckUserId(userId uint32) int8{
	if userId > casinoConf.MAX_USER_ID || userId < casinoConf.MIN_USER_ID {
		return casinoConf.LOGIN_WAY_QUICK
	}else{
		return casinoConf.LOGIN_WAY_LOGIN
	}
}


/**
快速登录
	快速登录模式需要 为用户分配一个id,并且返回给用户
 */
func QuickLogin(user *bbproto.ReqAuthUser,a gate.Agent){
	//1,判断入参是否正确
	uuid := user.Uuid
	log.Debug("header.code %v",uuid)
	if uuid == nil {
		log.Error("登录的时候uuid 为空,无法登陆")
	}

	//2,为用户分配id
	nuser := newUser()
	//3,返回数据给用户
	//给发送者回应一个 Hello 消息
	e := string("收到了消息")
	c := int32(1)
	var header bbproto.ProtoHeader
	header.UserId = &(nuser.Id)
	header.Error = &e
	header.Code = &c;
	var data bbproto.ReqAuthUser
	data.Header = &header
	a.WriteMsg(&data)

}


/**
普通登录

 */
func Login(user *bbproto.ReqAuthUser){

	//1,检测参数
	userId := user.Header.UserId
	log.Debug("需要登陆的userId %v",userId)



}

/**
 动态生成一个 User
 */
func newUser() mode.User{
	//1,创建user
	nuser := mode.User{}
	nuser.NickName = config.RandNikeName()
	nuser.Mid = bson.NewObjectId()

	//2,获取数据库连接和回话
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	s:=c.Ref()
	defer c.UnRef(s)

	//活的自增主键
	id, err := c.NextSeq(casinoConf.DB_NAME, casinoConf.DBT_T_USER, casinoConf.DB_ENSURECOUNTER_KEY)

	nuser.Id = uint32(id)

	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Insert(nuser)
	return  nuser

}
