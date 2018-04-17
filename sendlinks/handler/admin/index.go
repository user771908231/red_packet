package admin

import (
	"sendlinks/modules"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

const (
	UserName string = "admin"
	PassWord string = "admin"
)
//admin 主页面
func IndexHandler(ctx *modules.Context)  {
	ctx.HTML(200,"x-admin/index")
}
//admin 登陆页面
func LoginHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/login")
}
//admin 登陆验证
func IsLoginHandler(ctx *modules.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	res := bson.M{
		"code":0,
		"msg":"账号和密码不正确，请重新输入！",
	}
	defer func() {
		res_json,_ := json.Marshal(res)
		ctx.Write(res_json)
	}()

	if (username == UserName) && (password == PassWord) {
		res["code"] = 1
		res["msg"] = "登陆成功！"
		res["url"] = "/login"
		res_json,_ := json.Marshal(res)
		ctx.Write(res_json)
	}
	return
}

func ListsHandler(ctx *modules.Context) {

}

func AddHandler(ctx *modules.Context)  {

}

func PowerHandler(ctx *modules.Context)  {

}

func ClassiFicationHandler(ctx *modules.Context) {

}
