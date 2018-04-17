package weixin

import (
	"casino_redpack/modules"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"casino_common/common/model/agentModel"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
	"casino_common/utils/numUtils"
	"fmt"
)

//新增子代理
func AgentAddHandler(ctx *modules.Context) {

	ctx.HTML(200, "weixin/agent/add_agent")
}

type AddAgentForm struct {
	AgentId uint32 `form:"agent_id" binding:"Required"`
	UserId uint32 `form:"userid" binding:"Required"`
	RealName string `form:"real_name" binding:"Required"`
	Phone string `form:"phone" binding:"Required"`
	Passwd string `form:"passwd" binding:"Required"`
	PasswdConfirm string `form:"passwd_confirm" binding:"Required"`
}

func (form AddAgentForm) Error(macaron_ctx *macaron.Context, errors binding.Errors) {
	ctx := &modules.Context{Context: macaron_ctx}
	if len(errors) >0 {
		ctx.Error("表单参数异常！", "", 3)
	}
}

//新增子代理
func AgentAddPostHandler(ctx *modules.Context, form AddAgentForm) {
	ex_agent_info := agentModel.GetAgentInfoById(form.UserId)
	if ex_agent_info != nil {
		ctx.Ajax(-1, "该玩家已经是代理了！新增代理失败！", nil)
		return
	}

	user_info := userService.GetUserById(form.UserId)
	if user_info == nil {
		ctx.Ajax(-2, "该玩家id不存在！新增代理失败！", nil)
		return
	}

	new_agent_info := agentModel.AgentInfo{
		UserId: form.UserId,
		NickName: user_info.GetNickName(),
		RealName: form.RealName,
		Phone: form.Phone,
		OpenId: user_info.GetOpenId(),
		UnionId: user_info.GetUnionId(),
		RootId: 0,
		Pid: form.AgentId,
		Level: 1,
		Type: agentModel.AGENT_TYPE_2,
	}

	switch form.AgentId {
	case 1:
		new_agent_info.Pid = 0
		new_agent_info.Level = 1
		new_agent_info.RootId = 0
		new_agent_info.Type = agentModel.AGENT_TYPE_1
	case 2:
		new_agent_info.Pid = 0
		new_agent_info.Level = 1
		new_agent_info.RootId = 0
		new_agent_info.Type = agentModel.AGENT_TYPE_2
	default:
		parent_info := agentModel.GetAgentInfoById(form.AgentId)
		if parent_info == nil {
			ctx.Error( "父级代理不存在！","",0)
			return
		}
		new_agent_info.Pid = parent_info.UserId
		new_agent_info.RootId = parent_info.RootId
		new_agent_info.Level = parent_info.Level + 1
		new_agent_info.Type = agentModel.AGENT_TYPE_3
	}
	err := new_agent_info.Insert()

	if err != nil {
		ctx.Error("新增代理失败，错误:"+err.Error(),"",0)
		return
	}

	//设置新密码
	user_info.Pwd = proto.String(form.Passwd)
	userService.UpdateUser2Mgo(user_info)

	ctx.Success("新增代理成功！", "",0)
}

//编辑代理
func EditAgentHandler(ctx *modules.Context) {
	agent_id := uint32(ctx.QueryInt("agent"))
	edit_type := ctx.Query("type")
	value := ctx.Query("value")

	switch edit_type {
	case "invcode":
		new_code := numUtils.String2Uint32(value)
		if agentModel.GetAgentInfoById(new_code) == nil {
			ctx.Ajax(-1, "邀请码错误，修改失败！", nil)
			return
		}

		user_info := userService.GetUserById(agent_id)
		if user_info == nil {
			ctx.Ajax(-2, "该用户不存在，修改失败！", nil)
			return
		}

		user_info.InvCode = proto.String(fmt.Sprintf("%d", new_code))
		err := userService.UpdateUser2Mgo(user_info)

		if err == nil {
			ctx.Ajax(1, "修改成功！", nil)
		}else {
			ctx.Ajax(3, "修改失败！错误:"+err.Error(), nil)
		}
		return
	}

	ctx.Ajax(-8, "参数异常！", nil)

	ctx.HTML(200, "weixin/agent/edit_agent")
}
