package errors

import (
	"casino_common/common/Error"
	"casino_common/common/consts"
)

//所有的错误都会定义在这里

var ERR_DESK_ENTER error = Error.NewError(consts.ACK_RESULT_ERROR, "进入房间失败")

