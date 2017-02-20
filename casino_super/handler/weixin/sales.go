package weixin

import (
	"casino_super/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/proto/ddproto"
	"casino_super/model/agentModel"
	"casino_common/common/userService"
	"github.com/go-macaron/binding"
)

//售卖首页
func SalesIndexHandler(ctx *modules.Context) {
	uid := ctx.QueryInt("uid")
	var user *ddproto.User = nil
	db.C(tableName.DBT_T_USER).Find(bson.M{"id": uint32(uid)}, &user)

	ctx.Data["Uid"] = uid
	ctx.Data["User"] = user
	ctx.HTML(200, "weixin/agent/sales")
}

//售卖表单
type SalesForm struct {
	Uid uint32 `form:"uid" binding:"Required"`
	Num int64 `form:"num" binding:"Required"`
	Money float64 `form:"money"`
	Remark string `form:"remark"`
}
//售卖给用户
func SalesToUserHandler(ctx *modules.Context, form SalesForm, errs binding.Errors)  {
	if errs.Len() > 0 {
		ctx.Ajax(-1, "表单参数非法！请重新填写。",nil)
		return
	}
	wx_info := ctx.IsWxLogin()
	if wx_info == nil {
		ctx.Ajax(-2, "为该用户添加房卡失败！请重新登录！",nil)
		return
	}
	agent_id := agentModel.GetUserIdByOpenId(wx_info.OpenId)
	if agent_id == 0 {
		ctx.Ajax(-3, "为该用户添加房卡失败！请重新登录！",nil)
		return
	}
	roomCardNum := userService.GetUserRoomCard(agent_id)
	if roomCardNum < form.Num {
		ctx.Ajax(-4, "为该用户添加房卡失败！您的房卡数不足！",nil)
		return
	}
	_,err := userService.DECRUserRoomcard(agent_id, form.Num)
	if err != nil {
		ctx.Ajax(-5, "为该用户添加房卡失败，扣除房卡失败！",nil)
		return
	}else {
		_, err = userService.INCRUserRoomcard(form.Uid, form.Num)
		if err != nil {
			ctx.Ajax(-6, "为该用户添加房卡失败!",nil)
			return
		}
	}

	err = agentModel.AddNewSalesLog(agent_id, form.Uid, agentModel.RoomCard, form.Num, form.Money, form.Remark)
	if err != nil {
		ctx.Ajax(-7, "为该用户添加房卡成功！但生成充值记录失败。",nil)
		return
	}
	ctx.Ajax(1, "为该用户添加房卡成功！",nil)
}
