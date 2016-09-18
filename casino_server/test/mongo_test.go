package test

import (
	"testing"
	"fmt"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/utils/db"
)

func TestMongoUtils(t *testing.T) {
	user := f(10180)
	fmt.Println("user:", user)

}

//找有某个人的战绩

func f(id uint32) *bbproto.User {
	fmt.Println("-----sf-s--")

	//1,首先在 redis中去的数据
	var buser *bbproto.User = nil
	result := redisUtils.GetObj("t_user-10180", &bbproto.User{})
	if result == nil {
		fmt.Println("redis中没有找到user")
		// 获取连接 connection
		tuser := &mode.T_user{}
		db.Query(func(d *mgo.Database) {
			d.C(casinoConf.DBT_T_USER).Find(bson.M{"id": id}).One(tuser)
		})

		if tuser.Id < casinoConf.MIN_USER_ID {
			result = nil
		} else {
			//把从数据获得的结果填充到redis的model中
			buser, _ = Tuser2Ruser(tuser)
			if buser != nil {
				fmt.Println("开始保存")
			}
		}
	}

	//判断用户是否存在,如果不存在,则返回空
	return buser
}

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