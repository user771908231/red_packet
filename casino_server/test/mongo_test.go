package test

import (
	"testing"
	"fmt"
	"casino_server/conf/casinoConf"
	"casino_server/utils/db"
)

func TestMongoUtils(t *testing.T){
	id := db.GetNextSeq(casinoConf.DBT_T_USER)
	fmt.Println("哈哈哈",id)
}

