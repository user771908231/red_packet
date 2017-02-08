package userModel

import (
	"casino_common/utils/db"
	"errors"
	"crypto/md5"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_super/conf/config"
	"encoding/hex"
)

const ADMIN_TABLE_NAME string = "t_admin"

type User struct {
	Id uint32
	NickName string
	PassWd string
}

//通过id获取用户信息
func GetUserById(id uint32) *User {
	var err error = nil
	user_row := new(User)
	db.Query(func(d *mgo.Database) {
		err = d.C(ADMIN_TABLE_NAME).Find(bson.M{
			"id": id,
		}).One(user_row)
	})
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
	db.Query(func(d *mgo.Database) {
		err = d.C(ADMIN_TABLE_NAME).Find(bson.M{
			"nickname": user_name,
			"passwd": passwd,
		}).One(user_row)
	})
	if err != nil {
		return nil
	}
	return user_row
}

//插入一个新用户
func (user *User)Insert() error {
	id, err := db.GetNextSeq(config.DB_USER_SEQ)
	if err != nil {
		return errors.New("获取user_id自增键失败！")
	}
	user.Id = uint32(id)
	h := md5.New()
	h.Write([]byte(user.PassWd))
	user.PassWd = hex.EncodeToString(h.Sum(nil))
	db.Query(func(d *mgo.Database) {
		err = d.C(ADMIN_TABLE_NAME).Insert(user)
	})
	return err
}

//编辑用户资料
func (user *User)Save() error {
	var err error = nil
	db.Query(func(d *mgo.Database) {
		err = d.C(ADMIN_TABLE_NAME).Update(bson.M{
			"id": user.Id,
		},user)
	})
	return err
}
