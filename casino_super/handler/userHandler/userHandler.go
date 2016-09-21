package userHandler

import (
	"gopkg.in/macaron.v1"
	"casino_server/mode/dao/TUserDao"
	"log"
)

//展示用户的列表
func Users(ctx *macaron.Context) {
	//得到数据
	users := TUserDao.FindUsers(100)
	//渲染数据
	ctx.Data["users"] = users
	//输出到模板
	ctx.HTML(200, "users") // 200 为响应码
}

//充值
func Recharge(ctx *macaron.Context, logger *log.Logger) {
	userId := uint32(ctx.QueryInt("userId"))
	//logger.Println("通过userId[%v]查询数据库", userId)
	user := TUserDao.FindUserByUserId(userId)
	ctx.Data["user"] = user
	ctx.HTML(200, "user/recharge") // 200 为响应码
}

//提交充值
func RechargePost(ctx *macaron.Context){

}
