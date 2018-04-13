package linksModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
	"time"
	"new_links/model/keysModel"
	"fmt"
	"math"
	"math/rand"
	"new_links/conf"
)

type Links struct {
	ObjId 	bson.ObjectId		`bson:"_id"`
	GruopId	bson.ObjectId
	Url	string
	Id 		uint32
	Push 	float64
	Remarks string
	Weight  int
	//Visit	[]*DayVisit	//访问次数
	Time time.Time
	Status int
	Quota int
	ExcessId uint32
}
//每日推送记录
type DayVisit struct {
	ObjId 	bson.ObjectId		`bson:"_id"` //DayVisit.ObjId == Links.ObjId
	TimeUinx time.Time
	Visit	int
}

func (DV DayVisit) Insert() error {
	DV.TimeUinx = conf.Server.CurrentTime
	err := db.C(tableName.DB_LINKS_VISIT).Insert(DV)
	return err
}

func (DV DayVisit) Update() error {
	err := db.C(tableName.DB_LINKS_VISIT).Update(bson.M{"_id":DV.ObjId},bson.M{"$inc":bson.M{"visit":1}})
	return err
}

func GetDayVisit(ObjId bson.ObjectId) int {
	V := DayVisit{}
	err := db.C(tableName.DB_LINKS_VISIT).Find(bson.M{"_id":ObjId,"timeuinx":conf.Server.CurrentTime},&V)
	if err == nil {
		return V.Visit
	}else {
		fmt.Println("推送次数错误：",err)
		return 0
	}
}

func (L Links) Insert() error {
	L.ObjId = bson.NewObjectId()
	L.Time = time.Now()
	L.Status = 1
	V := new(DayVisit)
	V.ObjId = L.ObjId
	V.Insert()
	err := db.C(tableName.DB_LINKS_LISTS).Insert(L)
	return err
}

func (L Links) Update() error{
	err :=db.C(tableName.DB_LINKS_LISTS).Update(bson.M{"_id":L.ObjId},L)
	return err
}

func LinksIdDel(Id bson.ObjectId) error {
	L := GetLinkObjId(Id)
	defer func() {
			LInskPush(L.GruopId.Hex())
	}()
	err := db.C(tableName.DB_LINKS_LISTS).Remove(bson.M{"_id":Id})
	return err
}

func GetLinksAll(query bson.M,page int,number int) (int,[]*Links){
	list := []*Links{}
	_,count := db.C(tableName.DB_LINKS_LISTS).Page(query, &list, "-_id", page, number)
	return count,list
}

func GetKeysStatus() []bson.M {
	query := bson.M{"status":1}
	Keys := keysModel.GetKeysLists(query)
	list := []bson.M{}
	for _,item := range Keys{
		row := bson.M{
			"id":item.ObjId.Hex(),
			"keys":item.Keys,
		}
		list = append(list,row)
	}
	return list
}

type PostForm struct {
	Group	string `form:"group" binding:"Required"`
	Id    uint32  `form:"id" binding:"Required"`
	Url   string   `form:"url" binding:"Required"`
	Push 	int `form:"push" binding:"Required"`
	Remarks string  `form:"remarks"`
	Quota int `form:"quota" binding:Required`
	ExcessId uint32 `form:"excess_id" `
}

func Createlink(f PostForm) error{

	defer func() {
		LInskPush(f.Group)
	}()
	L := Links{
		GruopId:bson.ObjectIdHex(f.Group),
		Url:f.Url,
		Id:f.Id,
		Weight:f.Push,
		Remarks:f.Remarks,
		Quota :f.Quota,
		ExcessId:f.ExcessId,
	}
	err := L.Insert()
	return err
}



//stringObjId 获取
func GetLinkId(string string) *Links {
	L := new(Links)
	err := db.C(tableName.DB_LINKS_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},L)
	if err !=nil {
		return nil
	}

	return L
}
//ObjId 获取
func GetLinkObjId(string bson.ObjectId) *Links {
	L := new(Links)
	err := db.C(tableName.DB_LINKS_LISTS).Find(bson.M{"_id":string},L)
	if err !=nil {
		return nil
	}

	return L
}

func LinksStatus(id string,value int) error {
	L := GetLinkId(id)
	L.Status = value
	err := L.Update()
	return err
}

func GetGruopIdGroup(GruopId string) []*Links {
	L := []*Links{}
	err := db.C(tableName.DB_LINKS_LISTS).FindAll(bson.M{"gruopid":bson.ObjectIdHex(GruopId)},&L)
	if err != nil {
		return nil
	}
	return L
}

func LInskPush(string string) {
	val := GetGruopIdGroup(string)
	var B int
	var F float64
	for _,item := range val {
		B = B + item.Weight
	}

	for _,item := range val {
		L := GetLinkObjId(item.ObjId)
		F = float64(L.Weight) / float64(B)
		L.Push = FloatValue(F,2)
		err := L.Update()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func FloatValue(f float64,n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func LinksWeight(L []*Links) *Links {
	count := 0
	for _, goods := range L {
		count += int(goods.Push * 100)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand_num := r.Intn(count)

	min, max := 0, 0
	for _, goods := range L {
		max = min + int(goods.Push*100)
		if rand_num >= min && rand_num < max {
			return goods
		}
		min += int(goods.Push * 100)
	}
	return nil
}

func RandLink(GruopId bson.ObjectId) *Links {
	LL := GetGruopIdGroup(GruopId.Hex())
	index := RemoveIndex(LL)
	if index != nil {
		val := LinksWeight(index)
		return val
	}

	return nil
}

func RemoveIndex(L []*Links) []*Links {
	list := []*Links{}
	for i,item := range L {
		V := DayVisit{}
		err := db.C(tableName.DB_LINKS_VISIT).Find(bson.M{"_id":item.ObjId,"timeuinx":conf.Server.CurrentTime},&V)
		if err != nil {
			V := new(DayVisit)
			V.ObjId = item.ObjId
			V.Visit = 0
			V.Insert()
			db.C(tableName.DB_LINKS_VISIT).Find(bson.M{"_id":item.ObjId,"timeuinx":conf.Server.CurrentTime},&V)
		}
			//判断超过限额
		if V.Visit < item.Quota {
			list = append(list,L[i])
		}

	}
	return list
}

func (L *Links) Visssts() error {
	err := db.C(tableName.DB_LINKS_VISIT).Update(bson.M{"_id":L.ObjId,"timeuinx":conf.Server.CurrentTime},bson.M{"$inc":bson.M{"visit":1}})
	if err != nil {
		V := new(DayVisit)
		V.ObjId = L.ObjId
		V.Visit = 1
		err := V.Insert()
		return err
	}
	return err
}

//根据ID号查找
func GetIDSelcet(id uint32) *Links {
	L := new(Links)
	err := db.C(tableName.DB_LINKS_VISIT).Find(bson.M{"id":id},L)
	if err == nil {
		return L
	}
	return nil
}


//判断是否超过限额
func IsExcess(L *Links) *Links{
	fmt.Println(L.Id)
	V := DayVisit{}
	err := db.C(tableName.DB_LINKS_VISIT).Find(bson.M{"_id":L.ObjId,"timeuinx":conf.Server.CurrentTime},&V)
	if err == nil {
		if V.Visit >= L.Quota{
			if L.ExcessId != 0 {
				list_row := GetIDSelcet(L.ExcessId)
				IsExcess(list_row)
			}
			return nil
		}
		return L
	}
	return nil
}


type PostUpload struct {
	ObjId	string `form:"obj_id" binding:"Required"`
	Group	string `form:"group" binding:"Required"`
	Id    uint32  `form:"id" binding:"Required"`
	Url   string   `form:"url" binding:"Required"`
	Push 	int `form:"push" binding:"Required"`
	Remarks string  `form:"remarks"`
	Quota	int `form:"quota" binding:"Required"`
	ExcessId uint32 `form:"excess_id"`
}