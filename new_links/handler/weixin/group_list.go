package weixin

import (
	"casino_redpack/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/proto/ddproto"
	"encoding/json"
	"casino_common/common/log"
	"casino_common/common/service/groupService"
	"github.com/golang/protobuf/proto"
	"casino_common/common/model/agentModel"
	"casino_common/common/userService"
	"casino_common/common/service/rpcService"
	"golang.org/x/net/context"
)

//俱乐部列表
func GroupListHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	agent_info := userService.GetUserById(agent_id)

	list := []*ddproto.GroupInfo{}
	db.C(tableName.DBT_GROUP_INFO).FindAll(bson.M{
		"owner": agent_info.GetId(),
	}, &list)

	//ctx.Dump(list)

	ctx.Data["list"] = list
	ctx.HTML(200, "weixin/agent/group_list")
}

//俱乐部编辑
func GroupEditHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	agent_info := userService.GetUserById(agent_id)

	types := ctx.Query("type")
	group_id := ctx.QueryInt("id")

	group_info := new(ddproto.GroupInfo)
	if types == "edit" {
		if group_id <= 0 {
			ctx.Error("未找到该俱乐部！", "", 0)
			return
		}
		err := db.C(tableName.DBT_GROUP_INFO).Find(bson.M{
			"id": group_id,
			"owner": agent_info.GetId(),
		}, group_info)

		if err != nil {
			group_info = nil
			ctx.Error("未找到该俱乐部！", "", 0)
			return
		}
	}

	ctx.Data["group_info"] = group_info
	ctx.Data["type"] = types
	ctx.HTML(200, "weixin/agent/group_edit")
}

//俱乐部编辑表单
func GroupEditPostHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	agent_info := userService.GetUserById(agent_id)

	types := ctx.Query("type")
	data := ctx.Query("data")

	post_group_info := new(ddproto.GroupInfo)
	err := json.Unmarshal([]byte(data), post_group_info)


	if err != nil {
		log.T("err %v",err)
		ctx.Ajax(-1, err.Error(), nil)
		return
	}
	switch types {
	case "edit":
		ex_group_info := groupService.GetGroupInfoById(post_group_info.GetId())
		if ex_group_info == nil {
			ctx.Ajax(-2, "该俱乐部不存在！", nil)
			return
		}

		err := db.C(tableName.DBT_GROUP_INFO).Update(bson.M{"id": ex_group_info.GetId()}, bson.M{
			"$set": bson.M{
				"info": post_group_info.GetInfo(),
				"name": post_group_info.GetName(),
				"gameopts": post_group_info.GameOpts,
				"syncid": ex_group_info.GetSyncId()+1,
				},
		})

		if err != nil {
			ctx.Ajax(-3, err.Error(), nil)
			return
		}
		//大厅刷新数据
		rpcService.GetHall().RefreshGroupInfo(context.Background(), &ddproto.HallRpcGroupRefresh{GroupId:proto.Int32(post_group_info.GetId())})
		ctx.Ajax(1, "编辑成功！", nil)
		return
	case "add":
		if post_group_info.GetId() < 10000000 || post_group_info.GetId() > 99999999 {
			ctx.Ajax(-1, "请输入8位数字邀请码！", nil)
			return
		}
		if ex_count,_ := db.C(tableName.DBT_GROUP_INFO).Count(bson.M{"id": post_group_info.GetId()});ex_count > 0 {
			ctx.Ajax(-1, "该邀请码已存在，请重新输入。", nil)
			return
		}

		post_group_info.Info = proto.String("")
		post_group_info.Owner = proto.Uint32(agent_info.GetId())
		post_group_info.Members = []*ddproto.GroupMemberInfo{
			&ddproto.GroupMemberInfo{
				Uid: proto.Uint32(agent_info.GetId()),
				NickName: proto.String(agent_info.GetNickName()),
				Remark: proto.String(""),
				HeadImg: proto.String(agent_info.GetHeadUrl()),
				OpenId: proto.String(agent_info.GetOpenId()),
			},
		}

		post_group_info.SyncId = proto.Int32(1)

		//初始化配置id
		for i,opt := range post_group_info.GameOpts {
			opt.Id = proto.Int32(int32(i)+1)
			opt.Remark = proto.String("")
		}

		err := db.C(tableName.DBT_GROUP_INFO).Upsert(bson.M{"id": post_group_info.GetId()}, post_group_info)
		if err != nil {
			ctx.Ajax(-2, "创建俱乐部失败！请联系管理员", nil)
			return
		}
		//大厅刷新数据
		rpcService.GetHall().RefreshGroupInfo(context.Background(), &ddproto.HallRpcGroupRefresh{GroupId:proto.Int32(post_group_info.GetId())})
		ctx.Ajax(1, "创建俱乐部成功", nil)
		return
	}

}
