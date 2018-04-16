package redpack

import (
	"casino_redpack/modules"

	"github.com/skip2/go-qrcode"
	"fmt"
	"image/color"
)

//代理中心
func AgentHandler(ctx *modules.Context) {
	//二维码的存储路径
	FilePath := fmt.Sprintf("public/qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	//二维码所在路径
	UserQrcode := fmt.Sprintf("/qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	//二维码中的URL信息
	URl := fmt.Sprintf("%s/home/sign_up?extension_id=%d",ctx.Req.Host,ctx.IsLogin().Id)
	fmt.Println(FilePath)
	err := qrcode.WriteColorFile(URl, qrcode.Medium, 200,color.Transparent,color.Black,FilePath)

	if err !=nil {
		UserQrcode = ""
		fmt.Println("error")
	}
	ctx.Data["user_qrcode"] = UserQrcode
	ctx.HTML(200, "redpack/agent/index")
}

