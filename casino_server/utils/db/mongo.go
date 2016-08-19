package db

import (
	"casino_server/conf/casinoConf"
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
)


//保存数据
func  SaveMgoData(dbt string,data interface{}) error{
	//得到连接
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	s.DB(casinoConf.DB_NAME).C(dbt).Insert(data)
	return nil

}


//得到序列号
func GetNextSeq(dbt string) int32{
	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		return 0
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()
	id,_ :=  c.NextSeq(casinoConf.DB_NAME, dbt, casinoConf.DB_ENSURECOUNTER_KEY)
	return int32(id)
}