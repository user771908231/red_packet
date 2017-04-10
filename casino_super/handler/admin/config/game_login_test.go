package config

import (
	"testing"
	"reflect"
	"casino_common/common/consts/tableName"
	"casino_common/proto/ddproto"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
	db_init "casino_common/common/db"
	"casino_common/common/sys"
	"errors"
	conf_login "casino_login/conf"
	"casino_common/common/service/configService"
)

func init() {
	ConfigMap = map[string]ConfInfo {
		tableName.DBT_T_USER: ConfInfo{
			Name: "用户表",
			List: &[]ddproto.User{},
			Row: &ddproto.User{},
		},
	}
	db_init.InitMongoDb("192.168.199.200", 27017, "test", "id",[]string{})
	sys.InitRedis("192.168.199.200:6379","test")
}

type ZjhConfig struct {
	Id bson.ObjectId
	Name string `title:"名称" info:"详情"`
}

type ColInfo struct {
	Field string
	Type string
	Title string
	Info  string
	Value interface{}
}

//配置信息
type ConfInfo struct {
	Name string //名称
	List interface{}  //表类型
	Row interface{}   //单条配置类型
}

var ConfigMap map[string]ConfInfo

func TestStructTag(t *testing.T) {
	zjh_conf := ZjhConfig{
		Name: "haha",
		Id: bson.NewObjectId(),
	}

	val := reflect.TypeOf(zjh_conf)

	for i:=0;i<val.NumField();i++ {
		col_info := ColInfo{}
		col := val.Field(i)
		col_info.Field = col.Name
		col_info.Title = col.Tag.Get("title")
		col_info.Info = col.Tag.Get("info")
		t.Log(col_info)
		t.Log(col.Type.Kind().String())
	}

}

//获列字段信息
func GetColInfo(col_struct interface{}) []ColInfo {
	cols := []ColInfo{}
	val := reflect.TypeOf(col_struct)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i:=0;i<val.NumField();i++ {
		col_info := ColInfo{}
		col := val.Field(i)
		col_info.Field = col.Name
		col_info.Type = col.Type.Kind().String()
		col_info.Title = col.Tag.Get("title")
		col_info.Info = col.Tag.Get("info")
		val_val := reflect.ValueOf(col_struct)
		if val_val.Kind() == reflect.Ptr {
			val_val = val_val.Elem()
		}
		col_info.Value = val_val.FieldByName(col.Name).Interface()
		cols = append(cols, col_info)
	}

	return cols
}

//从
func TestGetCol(t *testing.T) {
	conf_info := ConfigMap[tableName.DBT_T_USER]

	list := conf_info.List
	row := conf_info.Row

	err := db.C(tableName.DBT_T_USER).FindAll(bson.M{"id": 10082}, list)
	db.C(tableName.DBT_T_USER).Find(bson.M{"id": 10082}, row)

	col_info := GetColInfo(row)

	t.Log(err, list)
	t.Log(row)
	row_elem := row.(*ddproto.User)
	t.Log(row_elem.GetNickName())
	t.Log(col_info)

}

//获取配置
func GetConfig(table_name string) (err error) {

	if conf_info,ok := ConfigMap[table_name]; ok {
		err = db.C(table_name).FindAll(bson.M{}, conf_info.List)
		return
	}else {
		return errors.New("错误！")
	}

	return nil
}

//赋值到指针
func TestRegist(t *testing.T) {
	serv := []conf_login.ConfStruct{}
	err := configService.Regist(tableName.DBT_GAME_CONFIG_LOGIN, &serv)
	t.Log(err)
	t.Log(serv)
}

//pull config

//update config key

//getconfig
