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
)

type Links struct {
	ObjId 	bson.ObjectId		`bson:"_id"`
	GruopId	bson.ObjectId
	Url	string
	Id 		uint32
	KeysId  bson.ObjectId
	Push 	float64
	LinkName	string			//链接
	Remarks string
	Weight  int
	Visit	int	//访问次数
	Time time.Time
	Status int
}

func (L Links) Insert() error {
	L.ObjId = bson.NewObjectId()
	L.Time = time.Now()
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
	Keys  string `form:"keys" binding:"Required"`
	Push 	int `form:"push" binding:"Required"`
	Remarks string  `form:"remarks"`
}

func Createlink(f PostForm) error{

	link := fmt.Sprintf("%sfrom=%dn/s?word=%s",f.Url,f.Id,keysModel.GetkeysId(f.Keys).Keys)

	defer func() {
		LInskPush(f.Group)
	}()
	L := Links{
		GruopId:bson.ObjectIdHex(f.Group),
		Url:f.Url,
		Id:f.Id,
		KeysId:bson.ObjectIdHex(f.Keys),
		Weight:f.Push,
		LinkName:link,		//链接
		Remarks:f.Remarks,
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
	return math.Trunc((f/pow10_n)*pow10_n) / pow10_n
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
	val := LinksWeight(LL)
	return val
}

func (L *Links) Visssts() error {
	err := db.C(tableName.DB_LINKS_LISTS).Update(bson.M{"_id":L.ObjId},bson.M{"$inc":bson.M{"visit":1}})
	return err
}


type PostUpload struct {
	ObjId	string `form:"obj_id" binding:"Required"`
	Group	string `form:"group" binding:"Required"`
	Id    uint32  `form:"id" binding:"Required"`
	Url   string   `form:"url" binding:"Required"`
	Keys  string `form:"keys" binding:"Required"`
	Push 	int `form:"push" binding:"Required"`
	Remarks string  `form:"remarks"`
}