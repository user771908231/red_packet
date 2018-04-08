package weixin

import (
	"casino_redpack/modules"
	"github.com/go-macaron/captcha"
	"casino_common/common/model/agentModel"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"casino_common/common/model/userDao"
)

//代理-账号密码登陆
func LoginPostHandler(ctx *modules.Context, cpt *captcha.Captcha) {
	if !cpt.VerifyReq(ctx.Req) {
		ctx.Error("验证码错误", "", 0)
		return
	}
	userid := ctx.QueryInt("name")
	pwd := ctx.Query("passwd")

	agent_info,err := agentModel.LoginByPwd(uint32(userid), pwd)
	if err != nil {
		ctx.Error(err.Error(), "", 0)
	}else {
		//保存用户信息至session
		user_info := userDao.FindUserById(agent_info.UserId)
		ctx.Session.Set("wx_user", mpoauth2.UserInfo{
			OpenId:user_info.GetOpenId(),
			Nickname: user_info.GetNickName(),
			Sex: int(user_info.GetSex()),
			City: user_info.GetCity(),
			Province: "",
			Country: user_info.GetLocation(),
			HeadImageURL: user_info.GetHeadUrl(),
			UnionId: user_info.GetUnionId(),
			})
		//跳转至发起微信登录时的页面
		if url := ctx.Session.Get("redirect"); url != nil{
			ctx.Session.Delete("redirect")
			ctx.Redirect(url.(string), 302)
			return
		}
		//ctx.JSON(200, userinfo)
		ctx.Redirect("/weixin/agent", 302)
	}
}
