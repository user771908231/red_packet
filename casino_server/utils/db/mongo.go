package db

import (
	"casino_server/conf/casinoConf"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
)



//活的链接
func GetMongoConn()(*mongodb.DialContext, error) {
	return  mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
}

//保存数据
func InsertMgoData(dbt string,data interface{}) error{
	//得到连接
	c, err := GetMongoConn()
	if err != nil {
		return err
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	error := s.DB(casinoConf.DB_NAME).C(dbt).Insert(data)
	return error

}

//更新数据
func UpdateMgoData(dbt string,data mode.BaseMode) error{
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		return err
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	error := s.DB(casinoConf.DB_NAME).C(dbt).Update(bson.M{"_id": data.GetMid()},data)
	return error
}


//得到序列号
func GetNextSeq(dbt string) int32{
	//连接数据库
	c,err := GetMongoConn()
	if err != nil{
		return 0
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)
	id,_ :=  c.NextSeq(casinoConf.DB_NAME, dbt, casinoConf.DB_ENSURECOUNTER_KEY)
	return int32(id)
}

//查询
func GetList(des interface{},dbt string,query interface{},sort string,limit int) interface{}{
	c,err := GetMongoConn()
	if err != nil{
		return 0
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)

	ret := []mode.BaseMode{}
	s.DB(casinoConf.DB_NAME).C(dbt).Find(query).Sort(sort).Limit(limit).All(&ret)
	return ret
}

func GetList2(f func(*mongodb.Session)){
	c,err := GetMongoConn()
	if err != nil{
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)

	f(s)
}