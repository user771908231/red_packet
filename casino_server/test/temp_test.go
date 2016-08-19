package mongodb

import (
	"testing"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"casino_server/conf/casinoConf"
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
)

func TestTemp(t *testing.T) {
	userData := &mode.T_user{}
	userData.Mid = bson.ObjectIdHex("57b6d0ba1ba69d3b3472e625")
	userData.NickName = "dongbing333"
	//err := db.InsertMgoData(casinoConf.DBT_T_USER,userData)
	//err := db.UpdateMgoData(casinoConf.DBT_T_USER,userData)
	//fmt.Println("是否出错:",err)

	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Update(bson.M{"_id": userData.GetMid()},userData)

}
