package manage

import (
	"casino_redpack/modules"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"math"
	"github.com/go-macaron/binding"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
	"log"
	"casino_common/common/service/pushService"
	"time"
)

type User struct {
	RoomCard            int64   `protobuf:"varint,1,opt,name=RoomCard" json:"RoomCard,omitempty"`
	Pwd                 string  `protobuf:"bytes,2,opt,name=pwd" json:"pwd,omitempty"`
	Coin                int64   `protobuf:"varint,3,opt,name=coin" json:"coin,omitempty"`
	Id                  uint32  `protobuf:"varint,4,opt,name=id" json:"id,omitempty"`
	NickName            string  `protobuf:"bytes,5,opt,name=nickName" json:"nickName,omitempty"`
	Scores              int32   `protobuf:"varint,6,opt,name=scores" json:"scores,omitempty"`
	LastDrawLotteryTime string  `protobuf:"bytes,7,opt,name=lastDrawLotteryTime" json:"lastDrawLotteryTime,omitempty"`
	LastSignTime        string  `protobuf:"bytes,8,opt,name=lastSignTime" json:"lastSignTime,omitempty"`
	SignTotalDays       int32   `protobuf:"varint,9,opt,name=signTotalDays" json:"signTotalDays,omitempty"`
	SignContinuousDays  int32   `protobuf:"varint,10,opt,name=signContinuousDays" json:"signContinuousDays,omitempty"`
	Diamond             int64   `protobuf:"varint,11,opt,name=Diamond" json:"Diamond,omitempty"`
	Diamond2            int64   `protobuf:"varint,12,opt,name=Diamond2" json:"Diamond2,omitempty"`
	OpenId              string  `protobuf:"bytes,13,opt,name=openId" json:"openId,omitempty"`
	UnionId             string  `protobuf:"bytes,14,opt,name=UnionId" json:"UnionId,omitempty"`
	HeadUrl             string  `protobuf:"bytes,15,opt,name=headUrl" json:"headUrl,omitempty"`
	City                string  `protobuf:"bytes,16,opt,name=city" json:"city,omitempty"`
	Sex                 int32   `protobuf:"varint,17,opt,name=sex" json:"sex,omitempty"`
	RobotType           int32   `protobuf:"varint,18,opt,name=robotType" json:"robotType,omitempty"`
	Ticket              int32   `protobuf:"varint,19,opt,name=ticket" json:"ticket,omitempty"`
	Bonus               float64 `protobuf:"fixed64,20,opt,name=bonus" json:"bonus,omitempty"`
	RegTime             time.Time   `protobuf:"varint,21,opt,name=regTime" json:"regTime,omitempty"`
	RegChannel          string  `protobuf:"bytes,22,opt,name=regChannel" json:"regChannel,omitempty"`
	AgentId             uint32  `protobuf:"varint,23,opt,name=agentId" json:"agentId,omitempty"`
	LastIp              string  `protobuf:"bytes,24,opt,name=lastIp" json:"lastIp,omitempty"`
	LastTime            time.Time   `protobuf:"varint,25,opt,name=lastTime" json:"lastTime,omitempty"`
	// 用户兑换信息
	RealName         string `protobuf:"bytes,26,opt,name=realName" json:"realName,omitempty"`
	PhoneNumber      string `protobuf:"bytes,27,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
	WxNumber         string `protobuf:"bytes,28,opt,name=wxNumber" json:"wxNumber,omitempty"`
	RealAddress      string `protobuf:"bytes,29,opt,name=realAddress" json:"realAddress,omitempty"`
	ChannelId        int32  `protobuf:"varint,30,opt,name=channelId" json:"channelId,omitempty"`
	NewUserAward     bool   `protobuf:"varint,31,opt,name=newUserAward" json:"newUserAward,omitempty"`
	RegIp            string `protobuf:"bytes,32,opt,name=regIp" json:"regIp,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}
func UserIndexHnadler(ctx *modules.Context)  {
	sort := ctx.Query("sort")
	if sort == "" {
		sort = "id"
	}
	ctx.Data["sort"] = sort
	ctx.HTML(200, "admin/manage/user/index")
}

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
	sort := "id"
	if sort = ctx.Query("sort"); sort == "" {
		sort = "id"
	}
	list := new([]*ddproto.User)
	err, count := db.C(tableName.DBT_T_USER).Page(query, list, sort, page, 10)

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

//注册统计列表
func UserRegAllHandler(ctx *modules.Context)  {
	info := []*User{}
	db.C(tableName.DBT_T_USER).FindAll(bson.M{},&info)

	ctx.Data["info"] = info


	ctx.HTML(200,"admin/data/user_reg")
	//return info
}
//活跃统计列表
func UserActiveAllHandler(ctx *modules.Context)  {
	info := []*User{}
	db.C(tableName.DBT_T_USER).FindAll(bson.M{},&info)

	ctx.Data["info"] = info


	ctx.HTML(200,"admin/data/user_active")
	//return info
}

//注册统计查询
func UserRegOneHandler(ctx *modules.Context) {
	ChannelId := ctx.QueryFloat64("ChannelId")
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)

	info := []*User{}
	if(ChannelId == 0){
		db.C(tableName.DBT_T_USER).FindAll(bson.M{},&info)
	}
	//长沙
	if(ChannelId == 1){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 31},bson.M{"channelid" : 32},bson.M{"channelid" : 33}},
			},&info)
		}else {
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"regtime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 31},bson.M{"channelid" : 32},bson.M{"channelid" : 33}},
			},&info)
		}

	}
	//岳阳
	if(ChannelId == 2){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 34},bson.M{"channelid" : 35}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"regtime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 34},bson.M{"channelid" : 35}},
			},&info)
		}
	}
	//四川
	if(ChannelId == 3){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 1},bson.M{"channelid" : 2},bson.M{"channelid" : 3},bson.M{"channelid" : 11},bson.M{"channelid" : 12},bson.M{"channelid" : 41},bson.M{"channelid" : 21},bson.M{"channelid" : 22}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"regtime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 1},bson.M{"channelid" : 2},bson.M{"channelid" : 3},bson.M{"channelid" : 11},bson.M{"channelid" : 12},bson.M{"channelid" : 41},bson.M{"channelid" : 21},bson.M{"channelid" : 22}},
			},&info)
		}
	}
	//白山
	if(ChannelId == 4){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 61},bson.M{"channelid" : 62}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"regtime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 61},bson.M{"channelid" : 62}},
			},&info)
		}
	}
	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/user_reg")
}
//活跃统计查询
func UserActiveOneHandler(ctx *modules.Context) {
	ChannelId := ctx.QueryFloat64("ChannelId")
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)

	info := []*User{}
	if(ChannelId == 0){
		db.C(tableName.DBT_T_USER).FindAll(bson.M{},&info)
	}
	//长沙
	if(ChannelId == 1){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 31},bson.M{"channelid" : 32},bson.M{"channelid" : 33}},
			},&info)
		}else {
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"lasttime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 31},bson.M{"channelid" : 32},bson.M{"channelid" : 33}},
			},&info)
		}

	}
	//岳阳
	if(ChannelId == 2){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 34},bson.M{"channelid" : 35}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"lasttime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 34},bson.M{"channelid" : 35}},
			},&info)
		}
	}
	//四川
	if(ChannelId == 3){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 1},bson.M{"channelid" : 2},bson.M{"channelid" : 3},bson.M{"channelid" : 11},bson.M{"channelid" : 12},bson.M{"channelid" : 41},bson.M{"channelid" : 21},bson.M{"channelid" : 22}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"lasttime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 1},bson.M{"channelid" : 2},bson.M{"channelid" : 3},bson.M{"channelid" : 11},bson.M{"channelid" : 12},bson.M{"channelid" : 41},bson.M{"channelid" : 21},bson.M{"channelid" : 22}},
			},&info)
		}
	}
	//白山
	if(ChannelId == 4){
		if(date1 == date2){
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"channelid" : 61},bson.M{"channelid" : 62}},
			},&info)
		}else{
			db.C(tableName.DBT_T_USER).FindAll(bson.M{
				"lasttime": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"channelid" : 61},bson.M{"channelid" : 62}},
			},&info)
		}
	}
	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/user_active")
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
		_, err = userService.INCRUserRoomcard(form.Id, form.RoomCard, int32(ddproto.CommonEnumGame_GID_SRC), "管理后台给用户充值房卡")
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
	err = pushService.PushUserData(form.Id)
	log.Println("push userData:", err)
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
