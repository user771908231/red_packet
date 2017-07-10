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
	LogFileSize int32
	LogFileCount int32
	WSAddr     string
	TCPAddr    string
	HttpIp     string
	HttpPort   int
	MaxConnNum int

	//mongo数据库相关的配置
	MongoIp   string
	MongoLogIp   string
	MongoPort int

	//redis配置
	RedisAddr string

	//curVersion的配置
	CurVersion int32

	ProdMod bool

	//外网Ip
	OutIp       string
	HallTcpAddr string

	//返利
	RebateMan  int64 //满1000
	RebateSong int64 //送10
}

func init() {
	data, err := ioutil.ReadFile("../conf/super.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)
}
