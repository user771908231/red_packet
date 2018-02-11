package redpack

import "casino_redpack/modules"

//充值
func RechargeHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/recharge")
}
