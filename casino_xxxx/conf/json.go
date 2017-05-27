package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"fmt"
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
	ProdMode  bool   //开发模式
}

func init() {
	//配置路径
	configPath := "../conf/json.json"
	if !Server.ProdMode {
		configPath = "conf/json.json"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)
}
