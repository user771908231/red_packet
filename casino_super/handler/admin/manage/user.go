package manage

import (
	"casino_super/modules"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"math"
	"github.com/go-macaron/binding"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
	"log"
)

//所有用户
func UserListHandler(ctx *modules.Context) {
	query := bson.M{
		"$and": []bson.M{
			bson.M{"id": bson.M{"$gt": 0}},
		},
	}
	if keyword := ctx.Query("keyword"); keyword != "" {
		query["$and"] = append(query["$and"].([]bson.M), bson.M{"$or": []bson.M{
			bson.M{"id": ctx.QueryInt("keyword")},
			bson.M{"nickname": bson.RegEx{keyword, "$i"}},
		}})
	}
	page := 1
	if page = ctx.QueryInt("page"); page <= 0 {
		page = 1
	}

	list := new([]*ddproto.User)
	err, count := db.C(tableName.DBT_T_USER).Page(query, list, "id", page, 10)

	for _,u := range *list {
		u.RoomCard = proto.Int64(userService.GetUserRoomCard(u.GetId()))
		u.Coin = proto.Int64(userService.GetUserCoin(u.GetId()))
		u.Diamond = proto.Int64(userService.GetUserDiamond(u.GetId()))
	}

	data := bson.M{
		"list": list,
		"page": bson.M{
			"count":      count,
			"list_count": len(*list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		},
	}

	if err != nil {
		ctx.Ajax(-1, err.Error(), nil)
	} else {
		ctx.Ajax(1, "成功返回列表成功!", data)
	}

	//ctx.Data["List"] = list
	//ctx.HTML(200, "admin/manage/user/all")
}

//用户更新表单
type UserUpdateForm struct {
	RoomCard  int64  `form:"RoomCard"`
	Coin      int64  `form:"coin"`
	Id        uint32 `form:"id"`
	NickName  string `form:"nickName"`
	Diamond   int64  `form:"Diamond"`
	Ticket    int32  `form:"ticket"`
	Bonus     float64 `form:"bonus"`
	OpenId    string `form:"openId"`
	UnionId   string `form:"UnionId"`
	HeadUrl   string `form:"headUrl"`
	Sex       int32  `form:"sex"`
	RobotType int32  `form:"robotType"`
}

//更新一个用户的数据
func UpdateUserHandler(ctx *modules.Context, form UserUpdateForm) {
	err := db.C(tableName.DBT_T_USER).Update(bson.M{"id": form.Id}, bson.M{"$set": form})
	if err != nil {
		ctx.Ajax(-1, err.Error(), form)
		return
	}
	ctx.Ajax(1, "更新成功！", form)
}

//充值表单
type RechargeForm struct {
	Id       uint32 `form:"Id" binding:"Required"`
	Coin     int64 `form:"Coin" binding:""`
	Diamond  int64 `form:"Diamond" binding:""`
	RoomCard int64 `form:"RoomCard" binding:""`
	Bonus    float64 `form:"bonus"`
	Ticket   int32 `form:"ticket"`
}

//用户充值
func RechargeHandler(ctx *modules.Context, form RechargeForm, errs binding.Errors) {
	log.Println(form)
	if len(errs) > 0 {
		ctx.Ajax(-1, "表单参数错误！", errs)
		return
	}
	var err error = nil
	if form.Coin != 0 {
		_, err = userService.INCRUserCOIN(form.Id, form.Coin)
		if err != nil {
			ctx.Ajax(-2, "充值金币失败！", nil)
			return
		}
		//userService.UpdateUser2MgoById()
	}
	if form.Diamond != 0 {
		_, err = userService.INCRUserDiamond(form.Id, form.Diamond)
		if err != nil {
			ctx.Ajax(-3, "充值钻石失败！", nil)
			return
		}
	}
	if form.RoomCard != 0 {
		_, err = userService.INCRUserRoomcard(form.Id, form.RoomCard)
		if err != nil {
			ctx.Ajax(-4, "充值房卡失败！", nil)
			return
		}
	}
	if form.Bonus != 0 {
		_, err = userService.INCRUserBonus(form.Id, form.Bonus)
		if err != nil {
			ctx.Ajax(-4, "充值红包失败！", nil)
			return
		}
	}
	if form.Ticket != 0 {
		_, err = userService.INCRUserTicket(form.Id, form.Ticket)
		if err != nil {
			ctx.Ajax(-4, "充值奖券失败！", nil)
			return
		}
	}

	ctx.Ajax(1, "充值成功！", form)
	//userService.SyncMgoUserMoney(form.Id) //同步user的数据到mgo
}

//删除单个用户
func DelUserHandler(ctx *modules.Context) {
	err := db.C(tableName.DBT_T_USER).Remove(bson.M{"id": ctx.ParamsInt("id")})
	if err != nil {
		//ctx.Error("删除失败！"+err.Error(), "/admin/manage/user/all", 3)
		ctx.Ajax(-1, "删除失败！", nil)
		return
	}
	ctx.Ajax(1, "删除成功！", nil)
}
