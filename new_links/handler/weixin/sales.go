package weixin

import (
	"casino_redpack/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/proto/ddproto"
	"casino_common/common/model/agentModel"
	"casino_common/common/userService"
	"github.com/go-macaron/binding"
	"github.com/golang/protobuf/proto"
	"math"
	"time"
	"casino_common/common/log"
)

//售卖首页
func SalesIndexHandler(ctx *modules.Context) {
	uid := ctx.QueryInt("uid")
	var user *ddproto.User = nil
	db.C(tableName.DBT_T_USER).Find(bson.M{"id": uint32(uid)}, &user)

	if user != nil {
		user.RoomCard = proto.Int64(userService.GetUserRoomCard(user.GetId()))
	}

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
		log.E("参数非法")
		ctx.Ajax(-1, "表单参数非法！请重新填写。",nil)
		return
	}
	wx_info := ctx.IsWxLogin()
	if wx_info == nil {
		log.E("为该用户添加房卡失败！请重新登录！")
		ctx.Ajax(-2, "为该用户添加房卡失败！请重新登录！",nil)
		return
	}
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	if agent_id == 0 {
		log.E("为该用户添加房卡失败！请重新登录！")
		ctx.Ajax(-3, "为该用户添加房卡失败！请重新登录！",nil)
		return
	}
	log.T("代理%d开始为用户%d充值房卡%d张", agent_id, form.Uid, form.Num)
	//代理之间无法跨级转房卡，除非该用户为总代理
	if agentModel.IsAgent(form.Uid) {
		my_agent_info := agentModel.GetAgentInfoById(agent_id)
		agent_info := agentModel.GetAgentInfoById(form.Uid)
		if my_agent_info.UserId == agent_info.Pid {
			//为上下级关系，直接充值
			//pass
		} else {
			if my_agent_info.Level == 1 && agent_info.RootId == my_agent_info.UserId {
				//为我的下级代理的代理
				//pass
			}else {
				//否则
				log.E("代理%d为代理%d充值失败，因为该代理不属于你的子代理。", my_agent_info.UserId, agent_info.UserId)
				ctx.Ajax(-7, "为该代理商添加房卡失败！因为该代理不属于你的子代理。",nil)
				return
			}
		}
	}

	//防止填负数，导致bug
	if form.Num < 0 {
		form.Num = 0
	}

	roomCardNum := userService.GetUserRoomCard(agent_id)
	if roomCardNum < form.Num {
		log.E("%d为用户%d充值失败,房卡不足（before:%d dec:%d）", agent_id, form.Uid, roomCardNum, form.Num)
		ctx.Ajax(-4, "为该用户添加房卡失败！您的房卡数不足！",nil)
		return
	}
	_,err := userService.DECRUserRoomcard(agent_id, form.Num, int32(ddproto.CommonEnumGame_GID_SRC), "代理商为用户充值房卡")
	if err != nil {
		log.E("代理%d为用户%d充值失败，扣除房卡%d失败", agent_id, form.Uid, form.Num)
		ctx.Ajax(-5, "为该用户添加房卡失败，扣除房卡失败！",nil)
		return
	}else {
		_, err = userService.INCRUserRoomcard(form.Uid, form.Num, int32(ddproto.CommonEnumGame_GID_SRC), "代理商为用户充值房卡")
		if err != nil {
			log.E("为用户%d充值房卡%d失败，err:%v", form.Uid, form.Num, err)
			ctx.Ajax(-6, "为该用户添加房卡失败!",nil)
			return
		}
	}

	err = agentModel.AddNewSalesLog(agent_id, form.Uid, agentModel.RoomCard, form.Num, form.Money, form.Remark)
	if err != nil {
		ctx.Ajax(-7, "为该用户添加房卡成功！但生成充值记录失败。",nil)
		return
	}
	log.T("代理%d为用户%d充值房卡%d成功", agent_id, form.Uid, form.Num)

	//如果该用户第一次被代理充值，则设置该用户为该代理的返利客户
	user := userService.GetUserById(form.Uid)
	if user.GetAgentId() == 0 {
		user.AgentId = proto.Uint32(agent_id)
		userService.UpdateUser2Mgo(user)
	}

	//尝试领取返利
	//for {
	//	rebate_err := agentModel.DoGetRebateRoomCard(agent_id)
	//	if rebate_err != nil {
	//		break
	//	}
	//}

	ctx.Ajax(1, "为该用户添加房卡成功！",nil)
}

//出售记录
func SalesLogHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	agent := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(agent.UnionId)
	list := []agentModel.SalesLog{}
	query := bson.M{
		"$and": []bson.M{
			bson.M{"agentid": bson.M{"$eq": agent_id}},
		},
	}
	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"addtime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"addtime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	_, count := db.C(tableName.DBT_AGENT_SALES_LOG).Page(query, &list, "-addtime", page, 10)

	data := bson.M{
		"list": list,
		"page": bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		},
	}

	ctx.Data["Logs"] = data
	ctx.Data["start_time"] = start_time
	ctx.Data["end_time"] = end_time
	ctx.HTML(200, "weixin/agent/sales_log")
}
