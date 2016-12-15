package logHandler

import (
	"gopkg.in/macaron.v1"
	"log"
	"casino_common/utils/db"
	"casino_super/model/logDao"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type SearchParams struct {
	UserId string
	DeskId string
	Level string
	Data string
}

func init() {
	db.Oninit("127.0.0.1", 51668, "test", "id")
}

func Post(reqLog logDao.ReqLog, ctx *macaron.Context, logger *log.Logger) {
	logger.Print(fmt.Sprintf("开始保存 d.id[%v] u.id[%v] level[%v] data[%v]", reqLog.DeskId, reqLog.UserId, reqLog.Level, reqLog.Data))
	error := logDao.SaveLog2Mgo(reqLog)
	if error != nil {
		logger.Print("保存失败 错误信息[%v]", error)
		return
	}
}

func Get(ctx *macaron.Context, logger *log.Logger) {
	m := bson.M{}
	userId := ctx.Query("userId")
	if userId != "" {
		m["userid"] = userId
	}

	deskId := ctx.Query("deskId")
	if deskId != "" {
		m["deskid"] = deskId
	}

	data := ctx.Query("data")
	if data != "" {
		m["data"] = data
	}

	level := ctx.Query("level")
	if level != "" {
		m["level"] = level
	}

	searchParams := SearchParams{
		UserId:userId,
		DeskId:deskId,
		Data:data,
		Level:level,
	}
	ctx.Data["searchParams"] = searchParams
	//logger.Print("u.id[%v] d.id[%v] level[%v] data[%v]", searchParams.UserId, searchParams.DeskId, searchParams.Level, searchParams.Data)
	logs := logDao.FindLogsByMap(m)
	ctx.Data["logs"] = logs
	ctx.HTML(200, "log/logs") // 200 为响应码
}