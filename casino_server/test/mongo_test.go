package test

import (
	"testing"
	"fmt"
	"casino_server/mode/dao/TUserDao"
)

func TestMongoUtils(t *testing.T) {
	users := TUserDao.FindUserByUserId(10189)
	fmt.Println(users)
}
