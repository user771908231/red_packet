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

	//mongo数据库相关的配置
	MongoIp   string
	MongoPort int

	//redis配置
	RedisAddr string

	//curVersion的配置
	CurVersion  int32
	ProdMode    bool
	ConsolePort int

	//push service配置
	HallTcpAddr string
}

func init() {
	//自动载入配置
	LoadJsonConfig()
}

//载入配置
func LoadJsonConfig() {
	//port 3802
	filename := ""
	if Server.ProdMode {
		filename = "../conf/csmj.json"
	} else {
		filename = "conf/csmj.json"
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("%v", err)
	}
	fmt.Println("读取到的配置文件的信息:%v", string(data))
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

}
