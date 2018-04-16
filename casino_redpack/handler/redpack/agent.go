package redpack

import (
	"casino_redpack/modules"

	"github.com/skip2/go-qrcode"
	"fmt"
)

//代理中心
func AgentHandler(ctx *modules.Context) {
	//user_info := ctx.IsLogin()
	//str := fmt.Sprintf("%s/home/sign_up?extension_id=%s",ctx.Req.Host,ctx.IsLogin().Id)
	FilePath := fmt.Sprintf("public/qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	UserQrcode := fmt.Sprintf("qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	fmt.Println(FilePath)
	err := qrcode.WriteFile("https://www.baidu.com", qrcode.Medium, 256, FilePath)

	if err !=nil {
		UserQrcode = ""
		fmt.Println("error")
	}
	ctx.Data["user_qrcode"] = UserQrcode
	ctx.HTML(200, "redpack/agent/index")
}

