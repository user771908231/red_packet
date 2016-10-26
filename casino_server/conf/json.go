package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"fmt"
	"runtime/debug"
)

var Server struct {
	LogLevel   string
	LogPath    string
	WSAddr     string
	TCPAddr    string
	MaxConnNum int

	//mongo数据库相关的配置
	MongoIp    string
	MongoPort  int

	//redis配置
	RedisAddr  string

	//curVersion的配置
	CurVersion int32
}

func init() {
	data, err := ioutil.ReadFile("../conf/server.json")
	if err != nil {
		debug.PrintStack()
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Print("当前的配置文件信息:", Server)

}
