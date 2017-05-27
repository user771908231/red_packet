package mjDeskHandler

import (
	"gopkg.in/macaron.v1"
	"casino_majiang/service/majiang"
	"casino_common/utils/numUtils"
	"time"
)

func Get(ctx *macaron.Context) {

	//取得数据
	code := ctx.Query("code")
	if code != "casino" {
		ctx.Redirect("/")
	}
	//渲染数据

	desks := []*majiang.MjDesk{}

	//先获取金币场房间
	for _, cRoom := range majiang.MjroomManagerIns.CMJRoomIns {
		desks = append(desks, cRoom.Desks...)
	}
	//再获取朋友桌房间
	desks = append(desks, majiang.MjroomManagerIns.FMJRoomIns.Desks...)

	ctx.Data["renderTime"] = time.Now().Format("2006-01-02 15:04:05") //固定参数
	ctx.Data["desks"] = desks

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

	//先找朋友桌
	desk := majiang.MjroomManagerIns.GetFMJRoom().GetDeskByDeskId(deskId)

	//找不到再找金币场
	if desk == nil {
		for _, cRoom := range majiang.MjroomManagerIns.CMJRoomIns {
			for _, cDesk := range cRoom.Desks {
				if cDesk.GetDeskId() == deskId {
					desk = cDesk
				}
			}
		}
	}

	if desk == nil || desk.Users == nil {
		ctx.Redirect("/")
		return
	}
	renderUser := []*majiang.MjUser{}
	for _, user := range desk.GetUsers() {
		if user != nil {
			renderUser = append(renderUser, user)
		}
	}

	ctx.Data["renderTime"] = time.Now().Format("2006-01-02 15:04:05") //固定参数
	ctx.Data["users"] = renderUser
	ctx.Data["desk"] = desk
	//输出到模板
	ctx.HTML(200, "mjdesk/users") // 200 为响应码
}

func GetBills(ctx *macaron.Context) {
	//取得数据
	id := ctx.Params("id")
	//room.GetDeskByIdAndMatchId()
	//渲染数据
	deskId := int32(numUtils.String2Int(id))
	if deskId == 0 {
		ctx.Redirect("/")
	}
	desk := majiang.MjroomManagerIns.GetFMJRoom().GetDeskByDeskId(deskId)

	if desk == nil || desk.Users == nil {
		ctx.Redirect("/")
	}
	renderUser := []*majiang.MjUser{}
	for _, user := range desk.GetUsers() {
		if user != nil {
			renderUser = append(renderUser, user)
		}
	}

	ctx.Data["renderTime"] = time.Now().Format("2006-01-02 15:04:05") //固定参数
	ctx.Data["users"] = renderUser
	ctx.Data["desk"] = desk
	//输出到模板
	ctx.HTML(200, "mjdesk/bills") // 200 为响应码
}

func GetTest(ctx *macaron.Context) {

	ctx.HTML(200, "mjdesk/test") // 200 为响应码
}