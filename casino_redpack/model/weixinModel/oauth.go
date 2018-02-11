package weixinModel

import (
	"casino_redpack/modules"
	"github.com/chanxuehong/rand"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"log"
	"strings"
)

const (
	oauth2RedirectURI = "/weixin/oauth/callback" // 填上自己的参数 http://wx.tondeen.com/weixin/oauth/callback
	oauth2Scope       = "snsapi_userinfo"                 // 填上自己的参数
)

var (
	oauth2Endpoint oauth2.Endpoint
)

//oauth验证第一步发起认证
func OauthLogin(ctx *modules.Context) {
	//微信浏览器
	if strings.Contains(ctx.Req.Header.Get("User-Agent"), "MicroMessenger"){
		state := string(rand.NewHex())
		ctx.Session.Set("state", state)
		AuthCodeURL := mpoauth2.AuthCodeURL(WX_APP_ID, "http://" + ctx.Req.Host + oauth2RedirectURI, oauth2Scope, state)
		log.Println("AuthCodeURL:", AuthCodeURL)

		ctx.Redirect(AuthCodeURL, 302)
	}else {
		//普通浏览器
		ctx.HTML(200, "weixin/agent/login")
	}
}

//oauth验证第二步回调，获取用户信息并存入session
func OauthCallBack(ctx *modules.Context) {
	code := ctx.Query("code")

	if ctx.Query("state") != ctx.Session.Get("state").(string) {
		ctx.Error("state错误！", "", 0)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}

	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		ctx.Error("根据code获取token失败！"+err.Error(),"",0)
		return
	}

	wx_info, err := mpoauth2.GetUserInfo(token.AccessToken, token.UnionId, "", nil)
	if err != nil {
		ctx.Error("根据token获取userinfo失败！"+err.Error(),"",0)
		return
	}

	//保存用户信息至session
	ctx.Session.Set("wx_user", *wx_info)
	//跳转至发起微信登录时的页面
	if url := ctx.Session.Get("redirect"); url != nil{
		ctx.Session.Delete("redirect")
		ctx.Redirect(url.(string), 302)
		return
	}
	//ctx.JSON(200, userinfo)
	ctx.Redirect("/weixin/agent", 302)
}
