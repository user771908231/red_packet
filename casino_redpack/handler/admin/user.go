package admin

import (
	"casino_redpack/modules"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"casino_redpack/model/userModel"
	"github.com/go-macaron/captcha"
	"fmt"
)

//需要达到指定等级
func NeedLogin(level int32) macaron.Handler {
	return func(ctx *modules.Context) {
		user := ctx.IsLogin()
		if user == nil {
			ctx.Redirect("/admin/login", 302)
			return
		}
		user = userModel.GetUserById(user.Id)
		if user == nil {
			ctx.Redirect("/admin/login", 302)
			return
		}
		ctx.Data["User"] = user

		if user.Level < level {
			const err_msg = "权限不足，操作失败！"
			if ctx.Req.Method == "POST" {
				ctx.Ajax(-1, err_msg, nil)
			}else {
				ctx.Error(err_msg,"", 0)
			}
		}
	}
}
//判断用户登陆 没有登陆直接到home/loginss
func UserNeedLogin(ctx *modules.Context){
		user := ctx.IsLogin()
		if user == nil {
			ctx.Redirect("/home/login", 302)
			return
		}
}
//显示面板数据
func ShowPanel(ctx *modules.Context) {
	user := ctx.IsLogin()
	ctx.Data["User"] = user
}
//注册
func SignUpHandler(ctx *modules.Context) {
	//ctx.Write([]byte("用户注册"))
	ctx.HTML(200, "admin/sign")
}
//验证码
func NeedCaptcha(ctx *modules.Context, cpt *captcha.Captcha) {
	if !cpt.VerifyReq(ctx.Req) {
		ctx.Error("验证码错误！", "", 1)
	}
}

//登录
func LoginHandler(ctx *modules.Context) {
	ctx.HTML(200, "admin/user/login")
}

//登录
func UserLoginHandler(ctx *modules.Context) {
	ctx.HTML(200, "admin/user/login_2")
}
//登出
func LoginOutHandler(ctx *modules.Context) {
	ctx.Session.Set("user", nil)
	ctx.Success("登出成功！", "/admin", 1)
}
//登录验证
type LoginForm struct {
	Name string `binding:"Required;MinSize(3);MaxSize(12)"`
	Passwd string `binding:"Required;MinSize(4);MaxSize(24)"`
	captcha_id string `binding:"Required;Size(15)"`
	Captcha string `binding:"Required;Size(4)"`
}

func (form LoginForm)Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs)>0 {
		my_ctx := modules.Context{Context:ctx}
		my_ctx.Error("登录失败！请检查用户名或密码。", "", 5)
	}
}
//登录POST
func LoginPostHandler(form LoginForm, ctx *modules.Context , VerificationCode *captcha.Captcha) {
	//if !VerificationCode.VerifyReq(ctx.Req) {
	//	ctx.Error("验证码错误！", "", 1)
	//	return
	//}
	user := userModel.Login(form.Name, form.Passwd)
	if user != nil {
		ctx.Session.Set("user", *user)
		ctx.Success("登录成功！", "/admin", 1)
	}else {
		ctx.Success("登录失败！", "/admin/login", 1)
	}
}

//用户登录POST
func UserLoginPostHandler(form LoginForm, ctx *modules.Context , VerificationCode *captcha.Captcha) {
	if !VerificationCode.VerifyReq(ctx.Req) {
		ctx.Error("验证码错误！", "", 1)
		return
	}
	user := userModel.Login(form.Name, form.Passwd)
	if user != nil {
		ctx.Session.Set("user", *user)
		ctx.Success("登录成功！", "/home", 1)
	}else {
		ctx.Success("登录失败！", "/admin/login", 1)
	}
}

//注册验证
type SiginUpTable struct {
	Name string `binding:"Required;MinSize(6);MaxSize(12)"`			//帐户名
	PasswdOne string `binding:"Required;MinSize(6);MaxSize(24)"`	//密码
	PasswdTwo string `binding:"Required;MinSize(6);MaxSize(24)"`	//重复密码
	captchaId string `binding:"Required;Size(15)"`				//验证id
	Captcha string `binding:"Required;Size(4)"`						//验证码
}

//func (sign SiginTable)Error(ctx *macaron.Context,errs binding.Errors) {
//	if len(errs) > 0 {
//		my_ctx := modules.Context{Context:ctx}
//		my_ctx.Error("注册失败！请检查用户名或密码。","",50)
//	}
//}

func SignUpTableValuesHandler(sign SiginUpTable,ctx *modules.Context, VerificationCode *captcha.Captcha) {

	fmt.Println(sign.Captcha)
	if !VerificationCode.VerifyReq(ctx.Req) {
		//ctx.Error("验证码错误！", "", 1)
		ctx.Ajax(500,"验证码错误！",nil)
		return
	}
	err,msg := userModel.TableValues(sign.Name,sign.PasswdOne,sign.PasswdTwo)
	if err == nil && msg == ""{
		user := userModel.Login(sign.Name,sign.PasswdOne)
		ctx.Session.Set("user",*user)
		//ctx.Redirect("/home",302)
		ctx.Ajax(304,"注册成功！",nil)
		//ctx.Success("注册成功！", "/home", 5)
	}else {
		//ctx.Success(msg, "/admin/sign_up", 10)
		ctx.Ajax(500,msg,nil)
	}

}

