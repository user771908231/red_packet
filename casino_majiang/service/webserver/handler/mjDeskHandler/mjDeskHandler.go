package mjDeskHandler

import (
	"gopkg.in/macaron.v1"
	"casino_majiang/service/majiang"
	"casino_server/utils/numUtils"
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

func GetDeskMJInfo(desk *majiang.MjDesk) string {
	if desk == nil || desk.AllMJPai == nil {
		return "暂时没有初始化麻将"
	}
	s := ""
	for i, p := range desk.AllMJPai {
		is, _ := numUtils.Int2String(int32(i))

		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + " (" + is + "-" + ii + "-" + p.LogDes() + ")"
	}
	return s
}

func getUserPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	} else {
		return user.GameData.HandPai.GetDes()
	}
}

func getUserPengPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.PengPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func getUserGnagPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.GangPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func getUserHuPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.HuPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}
func getUserInPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return user.GameData.HandPai.InPai.LogDes()

}
