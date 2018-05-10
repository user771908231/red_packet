package redpack

import (
	"casino_redpack/modules"

	"github.com/skip2/go-qrcode"
	"fmt"
	"image/color"
	"casino_redpack/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/userModel"
)

//代理中心
func AgentHandler(ctx *modules.Context) {
	id := ctx.IsLogin().Id
	//二维码的存储路径
	FilePath := fmt.Sprintf("public/qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	//二维码所在路径
	UserQrcode := fmt.Sprintf("/qrcode/%dqrcodr.png",ctx.IsLogin().Id)
	//二维码中的URL信息
	URl := fmt.Sprintf("http://%s/home/sign_up?extension_id=%d",ctx.Req.Host,ctx.IsLogin().Id)
	//生成二维码
	err := qrcode.WriteColorFile(URl, qrcode.Medium, 200,color.Transparent,color.Black,FilePath)

	if err !=nil {
		UserQrcode = ""
		fmt.Println("error")
	}
	ctx.Data["user_qrcode"] = UserQrcode
	ctx.Data["user"] = Shouxi(id)
	ctx.HTML(200, "redpack/agent/index")
}

func Shouxi(id uint32) []bson.M {
	log := agentModel.GetAgentRebateLogByIdList(id,agentModel.TimeObject())
	list := []bson.M{}
	for _,item := range log {
		row := bson.M{
			"name":IsName(item.RebateId),
			"money":fmt.Sprintf("%.2f",item.RebateMoeny),
		}
		list = append(list,row)
	}
	return list
}

func IsName(id uint32) string {
	user := userModel.GetUserById(id)
	if user.AccountNumber == "" {
		return user.NickName
	}
	return user.AccountNumber
}

