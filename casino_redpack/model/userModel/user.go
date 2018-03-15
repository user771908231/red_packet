package userModel

import (
	"casino_redpack/conf/config"
	"casino_common/utils/db"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

const USER_TABLE_NAME string = "t_redpack_user"

type User struct {
	Id       uint32
	NickName string
	HeadUrl  string
	OpenId   string
	UnionId  string
	PassWd   string
}

//通过id获取用户信息
func GetUserById(id uint32) *User {
	var err error = nil
	user_row := new(User)
	err = db.C(USER_TABLE_NAME).Find(bson.M{
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
	err = db.C(USER_TABLE_NAME).Find(bson.M{
		"nickname": user_name,
		"passwd":   passwd,
	}, user_row)
	if err != nil {
		return nil
	}
	return user_row
}

//验证帐户 密码
func TableValues(user_name string , passwd_one string , passwd_two string) error {
	var err error = nil

	if VerificationRegexp(user_name) {
		return err
	}
	if passwd_one != passwd_two {
		return err
	}
	return err
}

//正则

func VerificationRegexp(str string)  bool{
	r,_:=regexp.Compile("[^A-Za-z0-9]")
	b:=r.MatchString(str)
	if b {
		return true
	}
	return false
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
	err = db.C(USER_TABLE_NAME).Insert(user)
	return err
}

//编辑用户资料
func (user *User) Save() error {
	err := db.C(USER_TABLE_NAME).Update(bson.M{"id": user.Id}, user)
	return err
}
