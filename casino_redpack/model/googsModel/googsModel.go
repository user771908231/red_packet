package googsModel

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"math"
	"errors"
	"casino_common/utils/numUtils"
	"strings"
	rand2 "casino_common/utils/rand"
	"casino_common/common/log"
	wx "github.com/chanxuehong/wechat.v2/mch/core"
)

//物品结构
type GoogsCoin struct {
	ObjId    bson.ObjectId `bson:"_id"`
	Number 	int		//数量
	Price	float64	//价格
	Googstype int //类型 1 金币
	Discount  float64 //折扣
	Status		int //状态 0 未使用 1 使用
	Time  time.Time //时间
}

func (G *GoogsCoin) Insert() error {
	G.ObjId = bson.NewObjectId()
	G.Status = int(0)
	G.Time = time.Now()
	err := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).Insert(G)
	return err
}

//更新订单信息
func (G *GoogsCoin) Update() error{
	err := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).Update(bson.M{"_id": G.ObjId}, G)
	return err
}
//删除一条数据
func (G *GoogsCoin) Delete() error{
	err := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).Remove(bson.M{"_id":G.ObjId})
	return err
}

//获取全部物品

func GetWhole(qruey bson.M,page int) bson.M {
	GoogsCoin := []*GoogsCoin{}
	_,count := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).Page(qruey,&GoogsCoin, "-endtime", page, 10)
	list := bson.M{
		"count":      count,
		"list_count": len(GoogsCoin),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(20)),
		"data":[]bson.M{},
	}
	listData := []bson.M{}
	for i,item := range GoogsCoin {
		if i == count {
			continue
		}

		NewList := bson.M{
			"Id":item.ObjId.Hex(),
			"username":googtype(item.Googstype),
			"Time": item.Time.Format("2006-01-02 15:04:05"),
			"Number":item.Number,
			"Status":item.Status,
			"Discount":item.Discount,
		}
		listData = append(listData,NewList)
	}
	list["data"] = listData
	return list
}

func googtype(val int) string {
	switch val {
	case 1:
		return "金币"
	default:
		return "未知"
	}
}

type GoogsForm struct {
	Number int `binding:"Required;MinSize(3);MaxSize(12)"`
	Price float64 `binding:"Required;MinSize(4);MaxSize(24)"`
	Googstype int `binding:"Required;Size(15)"`
	Discount float64 `binding:"Required;Size(4)"`
}

func AddOne(F GoogsForm)error  {
	G := new(GoogsCoin)
	G.Number = F.Number
	G.Discount = F.Discount
	G.Googstype = F.Googstype
	G.Price = F.Price
	err := G.Insert()
	if err != nil {
		return err
	}
	return nil
}

func Operation(id string,types string) error {
	G := GetGoog(bson.ObjectIdHex(id))
	if types == "ok" {
		G.Status = 1
		err := G.Update()
		if err == nil {
			return nil
		}
		return err
	}else if types == "no" {
		G.Status = 0
		err := G.Update()
		if err == nil {
			return nil
		}
		return err
	}else if types == "del" {
		err := G.Delete()
		if err == nil {
			return nil
		}
		return err
	}else {
		return errors.New("未知的参数types的value",)
	}
}

func GetGoog(id bson.ObjectId) *GoogsCoin {
	G := new(GoogsCoin)
	err := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).Find(bson.M{"_id":id},G)
	if err == nil {
		return G
	}
	return nil
}

func GetGoogs() []bson.M {
	G := []*GoogsCoin{}
	err := db.C(tableName.TABLE_GOOGSCOIN_ROW_LISTS).FindAll(bson.M{"status":1},&G)
	if err == nil {
		list := []bson.M{}
		count := len(G)
		for i,item := range G {
			if i == count {
				continue
			}

			new_list := bson.M{
				"id":item.ObjId.Hex(),
				"munber":item.Number,
			}
			list = append(list,new_list)
		}

		return list
	}
	return nil
}

//生成支付订单号
//格式是: time+peymodel+userid+mealid+randstr，不合适的话需要和@kory商量
func GetWxpayTradeNo(payModelId int32, userId uint32, mealId string, tnow time.Time) string {
	//支付方式
	payModelIdStr, _ := numUtils.Int2String(payModelId)

	//userId 字符串
	userIdStr, _ := numUtils.Uint2String(userId)


	//时间字符串

	timeStr := wx.FormatTime(tnow)

	//随机数
	rand := rand2.Rand(100, 999)
	randStr, _ := numUtils.Int2String(rand)

	result := strings.Join([]string{timeStr, payModelIdStr, userIdStr, mealId, randStr}, "")
	log.T("玩家[%v]通过paymodel[%v]购买套餐[%v]时的订单号是:[%v]", userId, payModelId, mealId, result)
	return result
}


