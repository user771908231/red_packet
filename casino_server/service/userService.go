package userService

import (
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbproto"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/name5566/leaf/gate"
	"fmt"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/mode"
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
	var e string
	e = "收到了消息"
	var c int32
	c = 1
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
	nuser.Name = "testName1"
	nuser.Mobile = "18081922618"

	//2,保存在数据库

	c, err := mongodb.Dial("localhost", 51668)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	////创建自增字段
	err = c.EnsureCounter("test", "t_user", "id")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}


	id, err := c.NextSeq("test", "t_user", "id")
	nuser.Id = uint32(id)

	s:=c.Ref()
	defer c.UnRef(s)
	s.DB("test").C("t_user").Insert(nuser)
	return  nuser

}
