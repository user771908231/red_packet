package errors

import (
	"casino_common/common/Error"
	"casino_common/common/consts"
)

//错误都将放在这里
var ERR_DESK_NOT_FOUND error = Error.NewError(consts.ACK_RESULT_ERROR, "没有找到牌桌")
