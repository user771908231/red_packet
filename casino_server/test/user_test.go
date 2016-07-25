package mongodb

import (
	"testing"
	"casino_server/service/userService"
	"fmt"
)

func TestUser(t *testing.T){
	initSys()
	fmt.Println("----------------------------分割线----------------------------")

	addUserCoin(10007,1000)
}


func  addUserCoin(userId uint32,coin int64){
	user := userService.GetUserById(userId)
	fmt.Println("user:",user)
	*user.Coin += coin
	userService.SaveUser2RedisAndMongo(user)
}
