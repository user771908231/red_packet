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