package conf

import (
	"encoding/json"
	"fmt"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"time"
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
	MongoLogIp string
	MongoPort int

	//redis配置
	RedisAddr string
	RedisPwd  string

	//curVersion的配置
	CurVersion int32

	ProdMod bool

	//外网Ip
	OutIp       string
	HallTcpAddr string

	//返利
	RebateMan  int64 //满1000
	RebateSong int64 //送10

	//微信相关
	WxAppId string
	WxMchId string
	WxApiKey string
	WxAppSecret string

	//rpc配置
	HallRpcAddr string

	//站点相关
	SiteName string

	//当前年月日
	CurrentTime time.Time
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
	t := time.Now()
	y,m,d := t.Date()
	Server.CurrentTime = time.Date(y, m,d,0,0,0,0,t.Location())
	fmt.Println("当前年月日:",Server.CurrentTime)
	fmt.Println("读取到的配置文件的信息:", Server)
}
