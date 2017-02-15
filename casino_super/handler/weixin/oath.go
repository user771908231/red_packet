package weixin

import (
	"casino_super/modules"
	"github.com/chanxuehong/rand"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"log"
)

const (
	wxAppId           = "APPID"                           // 填上自己的参数
	wxAppSecret       = "APPSECRET"                       // 填上自己的参数
	oauth2RedirectURI = "http://192.168.1.129:8080/page2" // 填上自己的参数
	oauth2Scope       = "snsapi_userinfo"                 // 填上自己的参数
)

var (
	oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
)

//oath验证
func Oath(ctx *modules.Context) {
	state := string(rand.NewHex())
	ctx.Session.Set("state", state)
	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, state)
	log.Println("AuthCodeURL:", AuthCodeURL)

	ctx.Redirect(AuthCodeURL, 302)
}

func UserInfo(ctx *modules.Context) {
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

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		ctx.Error("根据token获取userinfo失败！"+err.Error(),"",0)
		return
	}

	ctx.JSON(200, userinfo)
}
