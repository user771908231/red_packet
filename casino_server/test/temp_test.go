package mongodb

import (
	"testing"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"casino_server/utils/db"
	"casino_server/conf/casinoConf"
	"fmt"
)

func TestTemp(t *testing.T) {
	userData := &mode.T_user{}
	userData.Mid = bson.ObjectIdHex("57a20f492d9fea18c0a02f45")
	userData.NickName = "dongbing3"
	//err := db.InsertMgoData(casinoConf.DBT_T_USER,userData)
	err := db.UpdateMgoData(casinoConf.DBT_T_USER,userData)
	fmt.Println("是否出错:",err)


}
