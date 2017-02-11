package modules

import "casino_super/model/userModel"

//是否登录
func (ctx *Context) IsLogin() *userModel.User {
	user := ctx.Session.Get("user")
	if user == nil {
		return nil
	}else {
		user_info := user.(userModel.User)
		return &user_info
	}
}
