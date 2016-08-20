package test

import (
	"testing"
	"fmt"
	"casino_server/utils/db"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
	"github.com/name5566/leaf/db/mongodb"
)

func TestTemp(t *testing.T) {
	//testN()
	//testMongoUtils()
	testMongoUtils2()
}

func testN() {
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
	}
	defer c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)

	ret := []mode.T_user{}
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"nickname":"dongbing333"}).Sort("id").Limit(20).All(&ret)
	fmt.Println("得到的结果:", ret)

	tuser := &mode.T_user{}
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"nickname": "dongbing333"}).One(tuser)
	fmt.Println("得到的结果2:", tuser)

}

func testMongoUtils() {
	fmt.Println("开始测试")

	ret := []mode.T_user{}
	users := db.GetList(ret, casinoConf.DBT_T_USER, bson.M{"nickname":"dongbing333"}, "id", 20)

	fmt.Println("得到的结果:", users)
}

func testMongoUtils2() {
	fmt.Println("开始测试")

	ret := []mode.T_user{}
	db.GetList2(func(s *mongodb.Session) {
		s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"nickname":"dongbing333"}).All(&ret)
	})
	fmt.Println("得到的结果.size():", len(ret))
	fmt.Println("得到的结果:", ret)
}




