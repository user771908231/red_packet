package config

import (
	"casino_admin/modules"
	"casino_common/proto/ddproto"
	"strings"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
)

//notice表单
type NoticeForm struct {
	ObjId            string `form:"ObjId"`
	Id               int32 `form:"Id"`
	NoticeType       int32 `form:"NoticeType"`
	NoticeTitle      string `form:"NoticeTitle"`
	NoticeContent    string `form:"NoticeContent"`
	NoticeMemo       string `form:"NoticeMemo"`
	Noticefileds     string `form:"Noticefileds"`
	ChannelId        string `form:"ChannelId"`
}

//查询objectid
type ObjId struct {
	ObjId bson.ObjectId `bson:"_id"`
}

func FindObjectId() {

}

//配置转表单
func Notice2Form(notice ddproto.TNotice) NoticeForm {
	form := NoticeForm{
		Id: notice.GetId(),
		NoticeType: notice.GetNoticeType(),
		NoticeTitle: notice.GetNoticeTitle(),
		NoticeContent: notice.GetNoticeContent(),
		NoticeMemo: notice.GetNoticeMemo(),
		Noticefileds: strings.Join(notice.Noticefileds, "\n"),
		ChannelId: notice.GetChannelId(),
	}

	return form
}

//表单转配置
func Form2Notice(form NoticeForm) ddproto.TNotice {
	notice := ddproto.TNotice{
		Id: &form.Id,
		NoticeType: &form.NoticeType,
		NoticeTitle: &form.NoticeTitle,
		NoticeContent: &form.NoticeContent,
		NoticeMemo: &form.NoticeMemo,
		Noticefileds: strings.Split(form.Noticefileds, "\n"),
		ChannelId: &form.ChannelId,
	}

	return notice
}

//notice list
func NoticeListHandler(ctx *modules.Context) {
	channel_id := ctx.Query("channel_id")

	notice_list := []ddproto.TNotice{}

	query := bson.M{
		"channelid": channel_id,
	}

	db.C(tableName.DBT_T_TH_NOTICE).FindAll(query, &notice_list)

	form_list := []NoticeForm{}

	for _,notice := range notice_list {
		form_list = append(form_list, Notice2Form(notice))
	}

	ids,_ := db.C(tableName.DBT_T_TH_NOTICE).FindAllId(query)

	for i,_ := range form_list{
		form_list[i].ObjId = ids[i].Hex()
	}

	ctx.Data["list"] = form_list
	ctx.Data["channel_id"] = channel_id
	ctx.HTML(200, "admin/config/notice/list")
}

//notice edit
func NoticeEditHandler(form NoticeForm, ctx *modules.Context) {
	notice_row := Form2Notice(form)

	err := db.C(tableName.DBT_T_TH_NOTICE).Update(bson.M{
		"_id": bson.ObjectIdHex(form.ObjId),
	}, notice_row)

	if err != nil {
		ctx.Ajax(-1, "更改失败！", nil)
		return
	}

	ctx.Ajax(1, "更改成功！", notice_row)
}
