package TUserDao

import (
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/utils/db"
	"casino_server/common/log"
)

func FindUsers(limit int) []mode.T_user {
	var users []mode.T_user
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_USER).Find(bson.M{}).Sort("-id").All(&users)
	})
	log.T("找到的user[%v]", users)
	return users
}

func FindUserByUserId(userId  uint32) mode.T_user {
	var user mode.T_user
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_USER).Find(bson.M{"id": userId}).One(&user)
	})
	log.T("找到的user[%v]", user)
	return user
}
