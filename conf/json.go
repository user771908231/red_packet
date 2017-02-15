package conf

import (
	"encoding/json"
	"fmt"
	"github.com/name5566/leaf/log"
	"io/ioutil"
)

var Server struct {
	LogLevel   string
	LogPath    string
	WSAddr     string
	TCPAddr    string
	MaxConnNum int

	MongoIp   string //mongo数据库相关的配置
	MongoPort int
	RedisAddr string //redis配置

	CurVersion      int32  //curVersion的配置
	BaseDownloadUrl string //默认下载地
}

func init() {
	data, err := ioutil.ReadFile("../conf/zjh.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)
}
