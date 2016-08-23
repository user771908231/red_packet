//author saplmm

package db

import (
	"casino_server/conf/casinoConf"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/common/log"
)



//活的链接
func GetMongoConn() (*mongodb.DialContext, error) {
	return mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
}

//保存数据
func InsertMgoData(dbt string, data interface{}) error {
	log.T("insert数据到数据库[%v]",data)
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

//更新数据通过_id来更新
func UpdateMgoData(dbt string, data mode.BaseMode) error {
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		return err
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	error := s.DB(casinoConf.DB_NAME).C(dbt).Update(bson.M{"_id": data.GetMid()}, data)
	return error
}


//得到序列号
func GetNextSeq(dbt string) (int32, error) {
	//连接数据库
	c, err := GetMongoConn()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)
	id, _ := c.NextSeq(casinoConf.DB_NAME, dbt, casinoConf.DB_ENSURECOUNTER_KEY)
	return int32(id), nil
}


//查询一个list
func Query(f func(*mgo.Database)) {
	c, err := GetMongoConn()
	if err != nil {
	}
	defer c.Close()

	//获取session
	s := c.Ref()
	defer c.UnRef(s)
	f(s.DB(casinoConf.DB_NAME))

}