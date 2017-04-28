package configModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
)

type GoodsInfo struct {
	Id 		bson.ObjectId 	`bson:"_id"`
	Goodsid 	int64		`bson:"goodsid"`
	Name 		string		`bson:"name"`
	Category 	int64		`bson:"category"`
	Pricetype 	int64		`bson:"pricetype"`
	Price 		float64		`bson:"price"`
	Goodstype 	int64		`bson:"goodstype"`
	Amount 		float64		`bson:"amount"`
	Discount 	string		`bson:"discount"`
	Image 		string		`bson:"image"`
	Isshow 		bool		`bson:"isshow"`
	Sort 		float64		`bson:"sort"`

}

//编辑--列表
func GoodsInfoOne(Id bson.ObjectId) *GoodsInfo{
	goods_info := new(GoodsInfo)
	db.C(tableName.DBT_T_GOODS_INFO).Find(bson.M{
		"_id": Id,
	}, goods_info)
	return goods_info
}

//编辑--提交
func GoodsEditUpdate(Id bson.ObjectId,Goodsid int64,Name string,Category int64,Pricetype int64,Price float64,Goodstype int64,Amount float64,Discount string,Image string,Isshow bool,Sort float64) (goods_info *GoodsInfo) {
	err := db.C(tableName.DBT_T_GOODS_INFO).Update(bson.M{"_id": Id},bson.M{
		"$set" : bson.M{
			"goodsid" : Goodsid,
			"name" : Name,
			"category" : Category,
			"pricetype" : Pricetype,
			"price" : Price,
			"goodstype" : Goodstype,
			"amount" : Amount,
			"discount" : Discount,
			"image" : Image,
			"isshow" : Isshow,
			"sort" : Sort,
		},})
	if err != nil {
		return nil
	}
	return goods_info
}
