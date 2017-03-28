package majiang

import (
	"casino_common/common/log"
)

/**
	只有判断别人打的牌的时候，需要过的时候才会请求这个协议，自己摸牌 需不需要过的时候不需要请求这个协议...
*/
func (desk *MjDesk) ActGuo(userId uint32) error {
	err := desk.ActGuoChangSha(userId)
	if err != nil {
		log.E("过牌出错:err ", err)
	}
	return err
}
