package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"fmt"
)

var Server struct {
	LogLevel                string
	LogPath                 string
	WSAddr                  string
	TCPAddr                 string
	MaxConnNum              int

	MongoIp                 string //mongo数据库相关的配置
	MongoPort               int
	RedisAddr               string //redis配置

	CurVersion              int32  //curVersion的配置
	BaseDownloadUrl         string //默认下载地址


				       //-----------------------斗地主的信息
	DDZ_CurVersion          int32  //斗地主的当前版本
	DDZ_IsUpdate            int32  //斗地主是否需要强制升级
	DDZ_IsMaintain          int32  //斗地主是否在维护中
	DDZ_MaintainMsg         string //斗地主的维护信息
	DDZ_ReleaseTag          int32  //斗地主已经发布的版本
	DDZ_DownloadUrl         string //斗地主的下载连接
	DDZ_LatestClientVersion int32  // 斗地主最后的版本号
	DDZ_IP                  string
	DDZ_PORT                int32
	DDZ_STATUS              int32

				       //------------------------麻将的信息
	MJ_CurVersion           int32
	MJ_IsUpdate             int32
	MJ_IsMaintain           int32
	MJ_MaintainMsg          string
	MJ_ReleaseTag           int32
	MJ_DownloadUrl          string
	MJ_LatestClientVersion  int32
	MJ_IP                   string
	MJ_PORT                 int32
	MJ_STATUS               int32
}

func init() {
	data, err := ioutil.ReadFile("../conf/login.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	fmt.Println("读取到的配置文件的信息:", Server)
}
