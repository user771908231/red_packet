package agentModel

import (
	"testing"
)

func TestGetUnifiedOrderResponse(t *testing.T) {
	//统一下单
	resp, err := GetUnifiedOrderResponse("58aa56dd6fec9cba7bbab634", 0.01, "神经游戏-代理充值-10房卡", "171.221.140.74", "oG9kZwrv8oFF9ja6WRlHxxMoJZoU")
	t.Log(err)
	t.Log(resp)
	//发起支付的交易数据
	trade_data := GetTradeData(resp.PrepayId)
	t.Log(trade_data)
}
