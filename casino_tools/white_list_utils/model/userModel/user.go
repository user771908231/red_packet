package userModel

import (
	"casino_tools/white_list_utils/conf/config"
	"casino_common/utils/db"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"casino_tools/white_list_utils/conf"
)

const ADMIN_TABLE_NAME string = "t_super_admin"

type User struct {
	Id       uint32
	NickName string
	PassWd   string
}

//通过id获取用户信息
func GetUserById(id uint32) *User {
	var err error = nil
	user_row := new(User)
	err = db.C(ADMIN_TABLE_NAME).Find(bson.M{
		"id": id,
	}, user_row)
	if err != nil {
		return user_row
	}
	return nil
}

//验证密码密码
func Login(user_name string, passwd string) *User {
	var err error = nil
	user_row := new(User)
	h := md5.New()
	h.Write([]byte(passwd))
	passwd = hex.EncodeToString(h.Sum(nil))
	err = db.C(ADMIN_TABLE_NAME).Find(bson.M{
		"nickname": user_name,
		"passwd":   passwd,
	}, user_row)
	if err != nil {
		return nil
	}
	return user_row
}

//使用配置的用户名和密码
func LoginByConfFile(user_name string, passwd string) *User {
	if user_name == "" || passwd == "" {
		return nil
	}

	if user_name != conf.Server.UserName || passwd != conf.Server.PassWord {
		return nil
	}

	return &User{
		Id: 10080,
		NickName: user_name,
		PassWd: passwd,
	}
}

//插入一个新用户
func (user *User) Insert() error {
	id, err := db.GetNextSeq(config.DB_USER_SEQ)
	if err != nil {
		return errors.New("获取user_id自增键失败！")
	}
	user.Id = uint32(id)
	h := md5.New()
	h.Write([]byte(user.PassWd))
	user.PassWd = hex.EncodeToString(h.Sum(nil))
	err = db.C(ADMIN_TABLE_NAME).Insert(user)
	return err
}

//编辑用户资料
func (user *User) Save() error {
	err := db.C(ADMIN_TABLE_NAME).Update(bson.M{"id": user.Id}, user)
	return err
}
