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
	MongoLogIp string
	MongoPort int

	DBName string

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

	//超级管理员id
	SuperAdminUser uint32

	//旺实富支付配置
	PaywapUserCode string  //旺实富分配的商户号
	PaywapCompKey  string  //旺实富分配的密钥

	//白名单用户名、密码
	UserName string
	PassWord string
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
