package game

import (
	"casino_common/utils/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/proto/ddproto"
	"reflect"
)

//表操作
func Query(tableName string, obj interface{}, query_obj bson.M, order string, page int, limit int) error {
	//通过id获取用户信息
	var err error = nil
	db.Query(func(d *mgo.Database) {
		bson_obj := bson.M{"$and": bson.M{}}
		bson_obj["$and"]["$or"] = query_obj
		query := d.C(tableName).Find(bson_obj)
		//用户名或id关键字

		if order != "" {
			query = query.Sort(order)
		}
		if limit > 0 {
			query = query.Limit(limit)
			if page > 0 {
				query = query.Skip(page*limit)
			}
		}
		//如果传入的对象是Slice则查询出所有的数据
		if reflect.TypeOf(obj).Kind() == reflect.Slice {
			err = query.All(obj)
		}else {
			err = query.One(obj)
		}
	})

	return err
}

//用户列表
func GetUserList(keyword string, order string, page int, limit int) ([]*ddproto.User, error) {
	bson_obj := bson.M{}
	if keyword != "" {
		bson_obj["$or"] = bson.M{
			"id": bson.M{"$regex":keyword,"$options":"$i"},
			"nickname": bson.M{"$regex":keyword,"$options":"$i"},
		}
	}
	users := []*ddproto.User{}
	err := Query(tableName.DBT_T_USER, users, bson_obj, "", 0, 10)
	return users, err
}
