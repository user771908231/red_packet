package redpack

import "casino_redpack/modules"

//绑定银行卡
func BindBankHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/bind_bank")
}

//银行卡列表
func BankListHandler(ctx *modules.Context) {
	list := `{
	"code": 1,
	"message": "success",
	"request": []
}`
	ctx.Write([]byte(list))
}

//银行卡添加
func BankAddHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/bank_add")
}
