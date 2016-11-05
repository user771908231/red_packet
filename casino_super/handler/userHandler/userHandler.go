package userHandler

import (
	"gopkg.in/macaron.v1"
	//"casino_server/mode/dao/TUserDao"
	"log"
	"casino_server/mode"
)

//展示用户的列表
func get(ctx *macaron.Context, logger *log.Logger) {
	//test data
	users := []mode.T_user{
		{
			Id	: 1,
			NickName: "user1",
			Coin	: 10001,
			Diamond	: 10,
		},
		{
			Id	: 2,
			NickName: "user2",
			Coin	: 10002,
			Diamond	: 20,
		},
		{
			Id	: 3,
			NickName: "user3",
			Coin	: 10003,
			Diamond	: 30,
		},
		{
			Id	: 3,
			NickName: "user3",
			Coin	: 10003,
			Diamond	: 40,
		},
	}
	//得到数据
	//users := TUserDao.FindUsers(100)
	//渲染数据
	ctx.Data["users"] = users
	//输出到模板
	ctx.HTML(200, "users") // 200 为响应码
}

func post(ctx *macaron.Context, logger *log.logger) {

}

func put(ctx *macaron.Context, logger *log.logger) {

}

func delete(ctx *macaron.Context, logger *log.logger) {

}
//
////充值
//func Recharge(ctx *macaron.Context, logger *log.Logger) {
//	userId := uint32(ctx.QueryInt("userId"))
//	//logger.Println("通过userId[%v]查询数据库", userId)
//	user := TUserDao.FindUserByUserId(userId)
//	ctx.Data["user"] = user
//	ctx.HTML(200, "user/recharge") // 200 为响应码
//}
//
////提交充值
//func RechargePost(ctx *macaron.Context){
//
//}
