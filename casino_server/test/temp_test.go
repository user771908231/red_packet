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
	userData.Mid = bson.ObjectIdHex("57b6c2911ba69d41e7384791")
	userData.NickName = "dongbing2"
	err := db.InsertMgoData(casinoConf.DBT_T_USER,userData)
	fmt.Println(err)
}
