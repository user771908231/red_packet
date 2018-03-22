package redpack

import (
	"casino_redpack/modules"
)

//扫雷接龙
func SaoleiJLHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/saoleijl")
}

//扫雷接龙发红包
func SaoleiJLAddHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/saoleijl_add")
}

//扫雷红包领取记录
func SaoleiRedOpenRecordHandler(ctx *modules.Context) {
	ctx.Data["redId"] = ctx.QueryInt("redId")
	ctx.HTML(200, "redpack/home/saoleijl_open_record")
}

//炸弹接龙-红包列表
func SaoleiPackListHandler(ctx *modules.Context) {
	list := `{
	"code": 1,
	"message": "success",
	"request": [{
		"id": 2332,
		"type": 5,
		"money": "2000.00",
		"moneyfa": "2000.00",
		"member_id": 1407,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 2,
		"nickname": "andy",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/Q0j4TwGTfTKy1HSTdt9MBa3X9WMV4WXo1bhqll3oMOH1bcGnAXiaoArXm1PT5nk8Z7EewiceaPRDgDYajYYm42yA\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4989,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 6,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4992,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 6,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4993,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 6,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4994,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 7,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4996,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 6,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4997,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 8,
		"tail_number": 6,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4998,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 9,
		"tail_number": 7,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}, {
		"id": 4999,
		"type": 5,
		"money": "10.00",
		"moneyfa": "10.00",
		"member_id": 1483,
		"all_membey": 10,
		"has_member": 7,
		"tail_number": 7,
		"nickname": "\u9648\u5e86\u5982",
		"headimgurl": "http:\/\/wx.qlogo.cn\/mmopen\/vi_32\/j77c2Dah8pqOnyfNp5wRBKNZRnwagGyC71l1XbPdzP1btsvefw3dYQgk7msnVSDSTKTH1MN0ZbpWMBSLKP1ETw\/132",
		"xianshi": 0,
		"jiaru": 0
	}]
}`

	ctx.Write([]byte(list))
}

//接龙列表
func SaoleiPackLqListHandler(ctx *modules.Context) {
	list := `{
	"code": 1,
	"message": "success",
	"request": "0"
}`

	ctx.Write([]byte(list))
}
