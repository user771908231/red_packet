package mjDeskHandler

import (
	"gopkg.in/macaron.v1"
	"casino_majiang/service/majiang"
	"casino_common/utils/numUtils"
)

func Get(ctx *macaron.Context) {

	//取得数据
	code := ctx.Query("code")
	if code != "casino" {
		ctx.Redirect("/")
	}
	//渲染数据
	ctx.Data["desks"] = majiang.FMJRoomIns.Desks
	//输出到模板
	ctx.HTML(200, "mjdesk/desks") // 200 为响应码
}

func GetUsers(ctx *macaron.Context) {
	//取得数据
	id := ctx.Params("id")
	//room.GetDeskByIdAndMatchId()
	//渲染数据
	deskId := int32(numUtils.String2Int(id))
	if deskId == 0 {
		ctx.Redirect("/")
	}
	desk := majiang.GetFMJRoom().GetDeskByDeskId(deskId)

	if desk == nil || desk.Users == nil {
		ctx.Redirect("/")
	}
	renderUser := []*majiang.MjUser{}
	for _, user := range desk.GetUsers() {
		if user != nil {
			renderUser = append(renderUser, user)
		}
	}
	ctx.Data["users"] = renderUser
	ctx.Data["desk"] = desk
	//输出到模板
	ctx.HTML(200, "mjdesk/users") // 200 为响应码
}

