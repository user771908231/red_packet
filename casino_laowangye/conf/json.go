package conf

import (
	"encoding/json"
	"fmt"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"casino_common/common/sys"
	"casino_common/proto/ddproto"
	"os"
	"time"
	"casino_laowangye/conf/config"
	lconf "github.com/name5566/leaf/conf"
	//"casino_common/common/service/pushService"
)

var Server struct {
	LogLevel   string
	LogPath    string
	LogFileSize int32
	LogFileCount int32
	WSAddr     string
	TCPAddr    string
	MaxConnNum int

	MongoIp   string //mongo数据库相关的配置
	MongoLogIp string
	MongoPort int
	RedisAddr string //redis配置
	RedisPwd string

	CurVersion      int32  //curVersion的配置
	BaseDownloadUrl string //默认下载地

	HallTcpAddr string  //大厅tcp端口

	//rpc addr
	NiuRpcAddr string
	HallRpcAddr string
	ConsolePort int

	DBName string
}

func init() {

}

func LoadConfig() {
	data, err := ioutil.ReadFile("../conf/laowangye.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)

	lconf.LogLevel = Server.LogLevel
	lconf.LogPath = Server.LogPath
	lconf.ConsolePort = Server.ConsolePort

	//初始化系统
	err = sys.SysInit(
		int32(ddproto.COMMON_ENUM_RELEASETAG_R_PRO),
		true,
		Server.RedisAddr,
		Server.DBName,
		Server.RedisPwd,
		Server.LogPath,
		"laowangye",
		Server.LogFileSize,
		Server.LogFileCount,
		Server.MongoIp,
		Server.MongoLogIp,
		Server.DBName,
		[]string{config.DBT_T_TH_GAMENUMBER_SEQ,config.DBT_LAOWANGYE_DESK})
	//判断初始化是否成功
	if err != nil {
		os.Exit(-1) //推出系统
	}

	//载入push服务配置
	//pushService.PoolInit(Server.HallTcpAddr)

	time.Sleep(time.Second * 3) //初始化3秒之后启动程序
}
