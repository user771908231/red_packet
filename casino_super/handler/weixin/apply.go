package weixin

import (
	"casino_super/modules"
	"github.com/go-macaron/binding"
	"casino_super/model/agentModel"
	"casino_common/common/userService"
)

func ApplyHandler(ctx *modules.Context) {

	ctx.HTML(200, "weixin/agent/apply")
}

//代理申请表单
type ApplyForm struct {
	NickName string `form:"name" binding:"Required"`
	Phone string `form:"phone" binding:"Required"`
	InvitedId uint32 `form:"invited"`
}

//代理申请表单处理
func ApplyPostHandler(ctx *modules.Context, errs binding.Errors, form ApplyForm) {
	if errs.Len() > 0 {
		ctx.Error("表单验证失败！", "/weixin/agent/apply", 3)
		return
	}
	wx_info := ctx.IsWxLogin()
	user_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)

	//验证InvitedId
	if form.InvitedId != 0 && userService.GetUserById(form.InvitedId) == nil {
		ctx.Error("验证推荐人Id失败！请检查表单。", "/weixin/agent/apply", 3)
		return
	}
	//插入表
	new_row := agentModel.ApplyRecord{
		Name: form.NickName,
		Phone: form.Phone,
		InvitedId: form.InvitedId,
		UserId: user_id,
	}
	err := new_row.Insert()
	if err != nil {
		ctx.Error("代理申请发送失败！", "/weixin/agent/apply", 3)
		return
	}
	ctx.Success("代理申请发送成功！我们的工作人员稍后会与您取得联系。", "", 0)
}
