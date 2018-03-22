package pay

import (
	"casino_redpack/modules"
	"fmt"
	"casino_redpack/handler/weixin"
)
//金币充值
func GoldRechargeHandler(ctx *modules.Context)  {
	val := ctx.QueryInt("totalFee")
	switch val {
	case 20:
		value := weixin.PayWapMoney(float64(0.01),ctx)
		fmt.Println(value)
	case 50:
		fmt.Println("50")
	case 100:
		fmt.Println("100")
	case 200:
		fmt.Println("200")
	case 500:
		fmt.Println("500")
	case 1000:
		fmt.Println("1000")
	case 2000:
		fmt.Println("2000")
	default:
		list := `{
			"code": 1,
			"message": "success",
			"request": []
			}`
		ctx.Write([]byte(list))
	}
}


