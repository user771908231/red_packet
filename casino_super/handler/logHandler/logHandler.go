package logHandler

import (
	"gopkg.in/macaron.v1"
	"casino_common/utils/db"
	"casino_super/model/logDao"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
	"strconv"
	"math"
	"casino_common/common/log"
)

type SearchParams struct {
	UserId    string
	DeskId    string
	Level     string
	Data      string
	CreatedAt string
}

type CodeValidate struct {
	Code string `json:"code" binding:"Required;Include(casino)"` //验证码casino
}

type CallBack struct {
	Successed bool
	Msg       string
}

var limiter = time.Tick(time.Millisecond * 200) //post请求的频率控制 200毫秒

func init() {
	db.Oninit("127.0.0.1", 51668, "test", "id")
}

func Post(reqLog logDao.ReqLog, ctx *macaron.Context) {
	log.T("收到 来自[%v]的日志上传post请求", ctx.Req.RemoteAddr)
	<-limiter
	log.T("开始处理 来自[%v]的日志上传数据[%v]", ctx.Req.RemoteAddr, reqLog)
	if reqLog.Logs == nil || len(reqLog.Logs) <= 0 {
		log.T("[%v]请求数据为空, 已返回", ctx.Req.RemoteAddr)
		return
	}

	successedCount := 0
	for i, logData := range reqLog.Logs {
		if logData.Data == "" || logData.DeskId == "" || logData.Level == "" || logData.UserId == "" {
			log.T(fmt.Sprintf("第[%v]条数据为空, 已跳过", i))
			continue
		}
		log.T(fmt.Sprintf("开始保存 d.id[%v] u.id[%v] level[%v] data[%v]", logData.DeskId, logData.UserId, logData.Level, logData.Data))
		error := logDao.SaveLog2Mgo(logData)
		if error != nil {
			log.T("保存失败 错误信息[%v]", error)
		}
		successedCount++
	}
	log.T("处理完成 [%v]请求上传共[%v]条数据, 已成功插入[%v]条数据", ctx.Req.RemoteAddr, len(reqLog.Logs), successedCount)
	return
}

func Get(ctx *macaron.Context) {
	log.T("收到 来自[%v]的查询日志的get请求", ctx.Req.RemoteAddr)
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

	createdAt := ctx.Query("createdAt")
	if createdAt != "" {
		m["createdAt"] = createdAt
	}

	searchParams := SearchParams{
		UserId:userId,
		DeskId:deskId,
		Data:data,
		Level:level,
		CreatedAt:createdAt,
	}
	ctx.Data["searchParams"] = searchParams
	log.T("查询条件 userId[%v] deskId[%v] level[%v] data[%v] createAt[%v]", searchParams.UserId, searchParams.DeskId, searchParams.Level, searchParams.Data, searchParams.CreatedAt)


	//分页控件
	page := ctx.Params("page")
	limit := ctx.Query("limit")
	if limit == "" {
		limit = "500" //默认每页100条数据
	}
	limitInt64, err := strconv.ParseInt(limit, 10, 64)

	skip := int64(0)
	pageInt64, err := strconv.ParseInt(page, 10, 64)
	if page != "" {
		if err != nil || (pageInt64 - 1) < 0 {
			pageInt64 = 1
		}
		skip = limitInt64 * (pageInt64 - 1)
	}


	//if userId == "" && deskId == "" && data == "" && level == "" {
	//	ctx.Data["logs"] = []logDao.LogData{}
	//} else {}

	log.T("开始查找, 查找条件userId[%v] deskId[%v] data[%v] level[%v] createAt[%v]", userId, deskId, data, level, createdAt)
	logs := logDao.FindLogsByMap(m, int(skip), int(limitInt64))
	ctx.Data["logs"] = logs

	count := len(logs) //总数
	log.T(fmt.Sprintf("已找到[%v]条记录", count))

	paginator := Paginator(int(pageInt64), int(limitInt64), int64(count))

	ctx.Data["recordCount"] = count
	//单页不显示分页控件
	if count <= int(limitInt64) {
		ctx.Data["paginator"] = nil
	}else {
		ctx.Data["paginator"] = paginator
	}

	ctx.Data["params"] = "?userId=" + userId + "&deskId=" + deskId + "&level=" + level + "&data=" + data + "&createdAt=" + createdAt + "&limit=" + limit

	ctx.HTML(200, "log/logs") // 200 为响应码
}

func Delete(code CodeValidate, ctx *macaron.Context) {
	log.T(fmt.Sprintf("[%v]请求清空数据库", ctx.RemoteAddr()))
	logDao.DeleteAllLogs2Mgo()
	callBack := CallBack{
		Successed : true,
		Msg : "",
	}
	ctx.JSON(200, callBack)
}




//分页方法，根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
func Paginator(page, prepage int, nums int64) map[string]interface{} {

	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		firstpage = page - 1
		lastpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i, _ := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		firstpage = page - 3
		for i, _ := range pages {
			pages[i] = start + i
		}
		firstpage = page - 1
		lastpage = page + 1
	default:
		pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page + 1
	//fmt.Println(pages)
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["firstpage"] = firstpage
	paginatorMap["lastpage"] = lastpage
	paginatorMap["currpage"] = page
	return paginatorMap
}
