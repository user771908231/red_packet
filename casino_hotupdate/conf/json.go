package conf

import (
	"encoding/json"
	"fmt"
	"github.com/name5566/leaf/log"
	"io/ioutil"
)

var Server struct {
	//http
	HttpIp     string
	HttpPort   int
	//mongo数据库相关的配置
	MongoIp   string
	MongoLogIp   string
	//redis配置
	RedisAddr string

	//其他配置
	UploadPath string

}

func init() {
	data, err := ioutil.ReadFile("../conf/hotupdate.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)
}
