package weixin

import (
	"casino_super/modules"
	"github.com/go-macaron/binding"
	"casino_super/model/agentModel"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
	"math"
	"casino_common/common/service/exchangeService"
	"casino_common/common/userService"
)

func ApplyHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	user_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	//验证代理商不得重复申请
	if agentModel.GetAgentInfoById(user_id) != nil {
		ctx.Success("您现在已经是代理商了，请不要重复申请！", "", 0)
		return
	}
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

	//验证代理商不得重复申请
	if agentModel.GetAgentInfoById(user_id) != nil {
		ctx.Success("您现在已经是代理商了，请不要重复申请！", "", 0)
		return
	}

	//验证InvitedId
	if form.InvitedId != 0 {
		invited_user := agentModel.GetAgentInfoById(form.InvitedId)
		if invited_user == nil {
			ctx.Error("验证推荐人Id失败！请检查表单。", "/weixin/agent/apply", 3)
			return
		}
	}
	//判断是否为重复申请
	exist_agent := agentModel.ApplyRecord{}
	err_exists := db.C(tableName.DBT_AGENT_APPLY_LOG).Find(bson.M{
		"userid": user_id,
		"status": exchangeService.PROCESS_ING,
	}, &exist_agent)
	if err_exists == nil {
		ctx.Error("您已经发过申请，请不要重复发送！", "", 3)
		return
	}
	//插入代理申请表
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

//申请审核列表
func ApplyLogHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	//ctx.Dump(agent_id)
	page := ctx.QueryInt("page")
	status := ctx.QueryInt("status")
	if page == 0 {
		page = 1
	}
	if status == 0 {
		status = 1
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{
				"status": status,
				"invitedid": agent_id,
			},
		},
	}

	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	list := []*agentModel.ApplyRecord{}
	_,count := db.C(tableName.DBT_AGENT_APPLY_LOG).Page(query, &list, "-requesttime", page, 10)
	ctx.Data["status"] = status
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "weixin/agent/apply_log")
}

//切换状态
func ApplySwitchState(ctx *modules.Context) {
	id := ctx.Query("id")
	status := ctx.QueryInt("status")
	if len(id) != 24 || status <= 0 || status > 3 {
		ctx.Ajax(-1, "参数错误", nil)
		return
	}
	row := agentModel.GetApplyRecordById(id)
	if row == nil {
		ctx.Ajax(-2, "切换状态失败！", nil)
		return
	}
	my_wx_info := ctx.IsWxLogin()
	my_agent_id := agentModel.GetUserIdByUnionId(my_wx_info.UnionId)
	if row.InvitedId != my_agent_id {
		ctx.Ajax(-2, "切换状态失败！", nil)
		return
	}
	row.Status = exchangeService.ExchangeState(status)
	row.ProcessTime = time.Now()
	err := row.Save()
	if err != nil {
		ctx.Ajax(-4, "您没有权限处理这个申请！", nil)
		return
	}
	//切换为申请成功,则插入代理信息表
	if row.Status == exchangeService.PROCESS_TRUE {
		//申请人信息
		user_info := userService.GetUserById(row.UserId)
		//当前代理商信息
		my_agent_info := agentModel.GetAgentInfoById(row.InvitedId)
		new_agent_info := agentModel.AgentInfo{
			UserId: row.UserId,
			NickName: user_info.GetNickName(),
			RealName: row.Name,
			Phone: row.Phone,
			OpenId: user_info.GetOpenId(),
			UnionId: user_info.GetUnionId(),
			Pid: row.InvitedId,
			Level: my_agent_info.Level+1,
			Type: agentModel.AGENT_TYPE_3,
		}
		if my_agent_info.Level == 1 {
			new_agent_info.RootId = my_agent_id
		}else if my_agent_info.Level > 1 {
			new_agent_info.RootId = my_agent_info.RootId
		}
		new_agent_info.Insert()
	}
	ctx.Ajax(1,"切换状态成功！", nil)
}
