package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"math/rand"
	"sync"
	"casino_redpack/model/userModel"
	"casino_common/utils/db"
	"encoding/json"

	"casino_common/common/consts/tableName"
	"fmt"
	"errors"
	"math"
	"strconv"
	"casino_common/common/log"
)

//红包详情表
var TABLE_NAME_REDPACK_INFO = tableName.TABLE_NAME_REDPACK_INFO
var TABLE_NAME_OPEN_PACKET_LISTS = tableName.TABLE_NAME_OPEN_PACKET_LISTS

//单个红包
type Redpack struct {
	ObjId       bson.ObjectId `bson:"_id"`
	CreatorUser uint32            //发包的人id
	CreatorName string            //发包人昵称
	CreatorHead string            //发包人头像
	Id          int32             //红包id
	Money       float64           //红包大小
	Lost        float64           //剩余大小
	Piece       int             //分成几份
	Type        RoomType          //红包类型
	TailNumber  int             //尾数 0-9, 默认-1
	RoomId      int32             //房间id
	OpenRecord  []*OpenRecordItem //开红包记录
	Time        time.Time         //发红包时间
}

//开红包记录
type OpenRecordItem struct {
	UserId uint32  //领红包的人id
	NickName string  //领红包的人的昵称
	Head string  //领红包的人的头像
	Money float64  //领了多少钱
	Time time.Time  //时间
	Is bool
}


//记录
type OpenPacketlist struct {
	RedpackId   int32    //红包id
	UserId uint32 //id
	NickName string	//昵称
	Head string  //领红包的人的头像
	Money float64  //领了多少钱
	Time time.Time  //时间
	GameType RoomType//游戏类型
	CreatorId uint32 //发红包人的ID
	Tail	int //雷号
}
//金币加减记录
type CoinAddSbtract struct {
	ObjId  bson.ObjectId `bson:"_id"`
	UserId uint32 //用户id
	SendOrOpenPacket int //0 发红包 1 开红包
	UserCoin float64 //用户总金币
	AddOrSubtract	float64 //加减的金币
	Msg 	string
	Time time.Time
}
//


func (C *CoinAddSbtract) Isert() error {
	C.ObjId = bson.NewObjectId()
	C.Time = time.Now()
	err := db.C(tableName.TABLE_COIN_ADD_SUBTRACT_NOTES).Insert(C)
	if err != nil {
		return errors.New("插入一条记录失败！")
	}
	return nil
}
//红包锁
var redLock sync.Map

//锁
func Lock(redId int32) {
	lock,_ := redLock.LoadOrStore(redId, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
}

//解锁
func UnLock(redId int32) {
	lock,ok := redLock.Load(redId)
	if ok {
		lock.(*sync.Mutex).Unlock()
	}
}

//保存到mongo数据库
func (redInfo *Redpack) Upsert() error {
	return db.C(TABLE_NAME_REDPACK_INFO).Upsert(bson.M{
		"_id": redInfo.ObjId,
	}, redInfo)
}

func (R *Redpack) Find(Id int32) *Redpack{
	err := db.C(TABLE_NAME_REDPACK_INFO).Find(bson.M{"id":Id},&R)
	if err != nil {
		return nil
	}
	return R
}
func (redInfo *Redpack) IsOpen(user *userModel.User) bool {
	//是否已经拆过红包了
	for _,item := range redInfo.OpenRecord {
		if item.UserId == user.Id {
			return false
		}
	}
	return true
}

//拆红包
func (redInfo *Redpack) Open(user *userModel.User) (bool,float64) {
	//加逻辑锁，保证线程
	Lock(redInfo.Id)
	defer UnLock(redInfo.Id)

	//是否还有余额
	if redInfo.Lost <= 0 || redInfo.Piece <= len(redInfo.OpenRecord) {
		return false,0
	}

	//是否已经拆过红包了
	for _,item := range redInfo.OpenRecord {
		if item.UserId == user.Id {
			return false,item.Money
		}
	}
	Am := int(redInfo.Lost*100)

	//开始拆红包
	open_money := getOpenRedMoney(Am, redInfo.Piece - len(redInfo.OpenRecord),redInfo.CreatorUser,redInfo.TailNumber,user)

	//更新红包余额
	redInfo.Lost = float64(Am - open_money)/100
	//更新开包记录
	redInfo.OpenRecord = append(redInfo.OpenRecord, &OpenRecordItem{
		UserId: user.Id,
		NickName: user.NickName,
		Head: user.HeadUrl,
		Money: float64(open_money)/100,
		Time: time.Now(),
	})
	data := &OpenPacketlist{
		RedpackId:redInfo.Id,
		CreatorId:redInfo.CreatorUser,
		GameType:redInfo.Type,
		UserId:user.Id,
		NickName:user.NickName,
		Head :user.HeadUrl,
		Money :float64(open_money)/100,
		Time:time.Now(),
		Tail:redInfo.TailNumber,
	}
	//更新到mongodb
	db.C(TABLE_NAME_OPEN_PACKET_LISTS).Insert(data)
	redInfo.Upsert()

	return true,float64(open_money)/100
}

//拆红包算法(剩余的钱、剩余的人)
func GetOpenRedMoney(lost_money float64, lost_person int) float64 {
	//参数合法性验证
	if lost_money < 0.01 || lost_person <= 0 {
		return 0
	}

	//只有一个人，则把钱全给他
	if lost_person == 1 {
		return float64(int(lost_money * 100))/100
	}

	//取0.01 - 平均金额 * 2 的值
	lost_score := int(lost_money * 100)
	avg_score := lost_score / lost_person
	res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
	if res_score == 0 {
		res_score = 1
	}

	return float64(res_score)/100
}
//拆红包算法(剩余的钱、剩余的人)
func getOpenRedMoney(lost_money int, lost_person int,id uint32,L int,u *userModel.User) int {
	//参数合法性验证
	if lost_money < 1 || lost_person <= 0 {
		return 0
	}

	//只有一个人，则把钱全给他
	if lost_person == 1 {
		return lost_money
	}

	//取0.01 - 平均金额 * 2 的值
	lost_score := lost_money
	avg_score := lost_score / lost_person
	res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
	weishu := GetWeishu(float64(res_score)/100)
	res_money := res_score
	//
	if u.Id == 10117 {
		log.T("米啊嘛%d",u.Id)
		log.T("尾数:%d",weishu)
		if weishu != L {
			fmt.Println("开红包算法------------------------------------------------------%d",u.Id)
			val := getOpenRedMoney(lost_money, lost_person,id ,L ,u )
			return val
		}else{
			if res_score == 0 {
				res_score = 1
			}
			return res_money
		}
	}else if id == 10117 && u.Id != 10117{
		log.T("米啊嘛%d",u.Id)
		log.T("尾数:%d",weishu)
		if weishu == L {
			fmt.Println("开红包算法------------------------------------------------------%d",u.Id)
			getOpenRedMoney(lost_money, lost_person,id ,L ,u )
		}else{
			if res_score == 0 {
				res_score = 1
			}
			return res_money
		}
	}

	//if res_score == 0 {
	//	res_score = 1
	//}

	//
	//if u.Id != 10117 && weishu  != L{
	//	log.T("%d",u.Id)
	//	getOpenRedMoney(lost_money, lost_person,id ,L ,u )
	//}
	//fmt.Println("开红包算法------------------------------------------------------%d",u.Id)
	//if id == 10117 && weishu  != L && u.Id != uint32(10117){
	//	log.T("米啊嘛%d",u.Id)
	//	//getOpenRedMoney(lost_money,lost_money,id,L,u)
	//	log.T("%d",getOpenRedMoney(lost_money,lost_money,id,L,u))
	//}


	return res_money
}

//根据ID获取用户记录列
func GetRedPacketRecord (Id uint32) []*OpenPacketlist{
	var err error = nil
	Redpack := []*OpenPacketlist{}
	err = db.C(TABLE_NAME_OPEN_PACKET_LISTS).FindAll(bson.M{
		"userid": Id,
	}, &Redpack)
	if err == nil {
		return Redpack
	}
	return nil
}

func GetRedPacketRecordRow (Id uint32,RedpackId int32) *OpenPacketlist{
	var err error = nil
	Redpack := new(OpenPacketlist)
	err = db.C(TABLE_NAME_OPEN_PACKET_LISTS).Find(bson.M{
		"userid": Id,
		"redpackid":RedpackId,
	}, &Redpack)
	if err == nil {
		return Redpack
	}
	return nil
}

func GetUserNameValues(Id uint32) *Redpack{

	Redpack := new(Redpack)
	err := db.C(TABLE_NAME_REDPACK_INFO).Find(bson.M{
		"CreatorUser":Id,
	},Redpack)
	if err != nil {
		fmt.Println("NO")
		return nil
	}
	fmt.Println("YES")
	return Redpack
}

func GetPacketSendRecord(Id uint32,Type int) []byte {
	var err error = nil
	Redpack := []*Redpack{}
	err = db.C(TABLE_NAME_REDPACK_INFO).FindAll(bson.M{
		"creatoruser":Id,
		"type":Type,
	},&Redpack)
	if err == nil {
		data := redpackListJson(Redpack,Id)

		return data
	}
	return nil
}

//解析成客户端需要的红包列表 游戏记录
func GetLists(list []*OpenPacketlist) []byte {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}
	res_list := []bson.M{}
	lenth := len(list)
	for i,item := range list {
		if i == lenth {
			continue
		}
		var str string
		if item.Tail == int(item.Money * 100)%10 {
			str = "中雷"
		}else {
			str = "未中雷"
		}
		new_item := bson.M{
			"id": item.RedpackId,
			"GameType": item.GameType,
			"Money": item.Money,
			"nickname": item.NickName,
			"headimgurl": item.Head,
			"Time":item.Time.Unix(),
			"win_money":str,
		}
		res_list = append(res_list, new_item)
	}

	res["request"] = res_list
	data,_ := json.Marshal(res)
	return data
}
//根据ID获取单个红包信息
func getdata(Id int32) *Redpack{
	Redpack := new(Redpack)
	err := db.C(TABLE_NAME_REDPACK_INFO).Find(bson.M{
		"id":Id,
	},&Redpack)
	if err != nil {
		return nil
	}
	return Redpack

}
//序列化红包信息
func redpackListJson(list []*Redpack, userId uint32) []byte {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}

	res_list := []bson.M{}
	lenth := len(list)
	for i,item := range list {
		if i == lenth {
			continue
		}

		new_item := bson.M{
			"id": item.Id,
			"GameType": item.Type,
			"Money": item.Money,
			"member_id": item.CreatorUser,
			"tail_number": item.TailNumber,
			"nickname": item.CreatorName,
			"headimgurl": item.CreatorHead,
			"Time":item.Time.Unix(),
			"win_money": len(item.OpenRecord),
		}
		res_list = append(res_list, new_item)
	}

	res["request"] = res_list
	data,_ := json.Marshal(res)
	return data
}

func OpenPacketDetails(Id int32,user_id uint32) []byte{
	Details := getdata(Id)
	user_Packet_row := GetRedPacketRecordRow(user_id,Details.Id)
	if Details != nil {
		redInfo := bson.M{
			"has_member":len(Details.OpenRecord),
			"all_membey":Details.Piece,
			"headimgurl":Details.CreatorHead,
			"nickname":Details.CreatorName,
			"Money":Details.Piece,
			"tail_number":Details.TailNumber,
			"moneyfa":user_Packet_row.Money,

		}
		redItemList := []bson.M{}

		for _,item := range Details.OpenRecord {
			List := bson.M{
				"money":item.Money,
				"code":int(item.Money * 100)%10,
				"headimgurl":item.Head,
				"nickname":item.NickName,
				"open_time":item.Time.Unix(),
				"deduct_money":FloatValue(item.Money*0.03,2),
			}
			redItemList = append(redItemList,List)

		}
		res := bson.M{
			"code": 1,
			"message": "success",
			"request": bson.M{
				"redInfo":redInfo,
				"redItemList":redItemList,
			},
		}
		data,_ := json.Marshal(res)
		return data
	}else {
		res := bson.M{
			"code": 1,
			"message": "success",
			"request": bson.M{
			},
		}
		data,_ := json.Marshal(res)
		return data
	}
}

func FloatValue(f float64,n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func GetWeishu(str float64) int {
	var val int
	s :=fmt.Sprintf("%.2f", str)
	by := []byte(s)
	lengt := len(by)
	for i,_ := range by {
		if i == lengt-1{
			val,_ = strconv.Atoi(string(by[i]))
		}

	}
	return val
}

func RandValue(lost_money int,lost_person int) int {
	lost_score := lost_money
	avg_score := lost_score / lost_person
	res_score := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(avg_score*2)
	return  res_score
}




