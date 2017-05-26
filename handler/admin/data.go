package admin

import (
	"casino_admin/modules"
	"casino_admin/model/dataModel"
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"fmt"
	"casino_common/common/service/statisticsService"
	"casino_common/common/consts/tableName"
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
//在房间玩家
func AtHome(ctx *modules.Context) {
	result := dataModel.AtHome()
	if result !=nil{
		ctx.Data["user"] = result
		ctx.HTML(200,"admin/data/atHome")
	}

}


//条件查询
func AtHomeList(ctx *modules.Context) {
	GameID := ctx.Query("gameID")
	result :=dataModel.AtHomeList(GameID)
	ctx.Data["user"] = result
	ctx.HTML(200,"admin/data/atHome")
}

//在线统计
func OnlineStatic(ctx *modules.Context) {
	a :=statisticsService.OnlineCountAll()
	fmt.Println(a)
	ctx.HTML(200,"admin/data/onlineStatic")
}

//在线统计列表
func OnlineStaticList(ctx *modules.Context) {
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)

	ctx.Data["date_start"] =date_start
	ctx.Data["date_end"] =date_end

	//一天之前
	d, _ := time.ParseDuration("-24h")
	date := date2.Add(d)
	date4 := date.Format("2006-01-02")
	date3,_ := time.Parse("2006-01-02",date4)
	err :=dataModel.OnlineStaticList(date1.Unix(),date2.Unix())
	fmt.Println(date1.Unix())
	fmt.Println(date2.AddDate(0,0,+1).Unix())
	fmt.Println(err)
	ctx.Data["static"] =err
	if(date1 != date3){
		err :=dataModel.OnlineStaticDay(date1.Unix(),date2.Unix())
		ctx.Data["static"] =err
	}else {

	}


	ctx.HTML(200,"admin/data/onlineStatic")
}
//用户使用房卡统计
const T_STATISTICS_ROOMCARD string= "t_statistics_roomcard"

type T_statistics_roomcard struct {
	Id        		bson.ObjectId     	`bson:"_id"` 		//用户ID
	UserId			uint32          	`bson:"userid"`  	 //用户ID
	Gid              	float64          	`bson:"gid"`    	 //游戏ID
	Memo              	string          	`bson:"remo"`    	 //说明
	RoomCardCount           int64          		`bson:"roomcardcount"`    	 //房卡消耗数
	Time           		time.Time          	`bson:"time"`    	 //房卡消耗时间
}

func RoomCard(ctx *modules.Context) {
	info := []*T_statistics_roomcard{}
	db.C(T_STATISTICS_ROOMCARD).FindAll(bson.M{},&info)
	count,_ := db.C(T_STATISTICS_ROOMCARD).Count(bson.M{})
	ctx.Data["count"] = count
	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/roomCard")
}

//房卡每日消耗统计
const T_STATISTICS_ROOMCARD_DAY_DETAILS string= "t_statistics_roomcard_day_details"

type T_statistics_roomcard_day_details struct {
	Id        		bson.ObjectId     	`bson:"_id"` 		//用户ID
	Gid              	float64          	`bson:"gid"`    	 //游戏ID
	RoomCardCount           int64          		`bson:"roomcardcount"`    	 //房卡消耗数
	Time           		time.Time          	`bson:"time"`    	 //房卡消耗时间
}

func RoomCardDay(ctx *modules.Context) {
	info := []*T_statistics_roomcard_day_details{}

	Time := time.Now().Format("2006-01-02")
	date1,_ := time.Parse("2006-01-02",Time)

	db.C(T_STATISTICS_ROOMCARD_DAY_DETAILS).FindAll(bson.M{
		"time": bson.M{"$gte": date1},
	},&info)
	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/roomCardDay")
}

//用户使用房卡统计查询
func RoomCardOne(ctx *modules.Context) {
	Gid := ctx.QueryFloat64("Gid")
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)
	date2 =date2.AddDate(0,0,+1)
	info := []*T_statistics_roomcard{}
	if(Gid == 0){
		db.C(T_STATISTICS_ROOMCARD).FindAll(bson.M{},&info)
	}else{
		if(date1.String() == "0001-01-01 00:00:00 +0000 UTC"){
			db.C(T_STATISTICS_ROOMCARD).FindAll(bson.M{
				"$or" :[]bson.M{bson.M{"gid" : Gid}},
			},&info)
		}else {
			db.C(T_STATISTICS_ROOMCARD).FindAll(bson.M{
				"time": bson.M{"$gte": date1,"$lte": date2},
				"$or" :[]bson.M{bson.M{"gid" : Gid}},
			},&info)
		}
	}


	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/roomCard")
}

//用户使用房卡统计查询
func RoomCardDayOne(ctx *modules.Context) {
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)
	date2 =date2.AddDate(0,0,+1)

	info := []*T_statistics_roomcard_day_details{}

	if(date1.String() == "0001-01-01 00:00:00 +0000 UTC"){
		db.C(T_STATISTICS_ROOMCARD_DAY_DETAILS).FindAll(bson.M{},&info)
	}else {
		db.C(T_STATISTICS_ROOMCARD_DAY_DETAILS).FindAll(bson.M{
			"time": bson.M{"$gte": date1,"$lte": date2},
		},&info)
	}



	ctx.Data["info"] = info
	ctx.HTML(200,"admin/data/roomCardDay")
}
