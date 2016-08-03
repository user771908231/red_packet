package mongodb

import (
	"testing"
	"casino_server/service/userService"
	"fmt"
)

func TestUser(t *testing.T){
	initSys()
	fmt.Println("----------------------------分割线----------------------------")

	//addUserCoin(10010,87000)
	//addUserCoin(10011,989000)
	//getUser(10084)
	getUser(10087)
	getUser(10084)
}


func  addUserCoin(userId uint32,coin int64){
	user := userService.GetUserById(userId)
	fmt.Println("user:",user)
	*user.Coin += coin
	userService.SaveUser2RedisAndMongo(user)
}

func getUser(userId uint32){
	user := userService.GetUserById(userId)
	fmt.Println("user.NickName",*user.NickName)
	fmt.Println("user.Name",*user.Name)
	fmt.Println("user.id",*user.Id)
	fmt.Println("user.coin",*user.Coin)
}
