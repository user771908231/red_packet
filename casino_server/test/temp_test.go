package test

import (
	"testing"
	"fmt"
	"casino_server/utils/db"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2"
)

func TestTemp(t *testing.T) {
	//testN()
	//testMongoUtils()
	testMongoUtils2()
}

func testN() {
	c, err := db.GetMongoConn()
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

func testMongoUtils2() {
	fmt.Println("开始测试")

	ret := []mode.T_user{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_USER).Find(bson.M{}).All(&ret)
	})
	fmt.Println("得到的结果.size():", len(ret))
	fmt.Println("得到的结果:", ret)
}




