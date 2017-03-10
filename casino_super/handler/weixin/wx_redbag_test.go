package weixin

import (
	"testing"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers"
	"gopkg.in/mgo.v2/bson"
	"casino_super/model/weixinModel"
)

//微信红包测试
func TestRedBag(t *testing.T) {
	https_client, err := core.NewTLSHttpClient("D:/MyGo/src/conf/apiclient_cert.pem", "D:/MyGo/src/conf/apiclient_key.pem")
	if err != nil {
		t.Log("https err:", err)
		return
	}
	req_map := map[string]string{
		"nonce_str": bson.NewObjectId().Hex(),  //随机数
		"mch_billno": bson.NewObjectId().Hex(),  //订单号
		"mch_id": weixinModel.WX_MCH_ID,
		"wxappid": weixinModel.WX_APP_ID,
		"send_name": "神经棋牌",   //商户名称
		"re_openid": "oG9kZwrv8oFF9ja6WRlHxxMoJZoU",  //用户openid
		"total_amount": "100",  //付款金额，单位（分）
		"total_num": "1",  //红包发放总人数
		"wishing": "这是祝福语。",  //红包祝福语
		"client_ip": "171.221.137.148",  //调用接口的机器Ip地址
		"act_name": "领红包活动",  //活动名称
		"remark": "红包兑换",  //备注信息
	}
	req_map["sign"] = core.Sign(req_map, weixinModel.WX_API_KEY, nil)
	t.Log(req_map)
	client := core.NewClient(weixinModel.WX_APP_ID,weixinModel.WX_MCH_ID, weixinModel.WX_API_KEY, https_client)
	resp,err := mmpaymkttransfers.SendRedPack(client, req_map)
	if err != nil {
		t.Log("resp err:", err)
		return
	}
	t.Log("resp:", resp)
}

func TestSendText(t *testing.T) {
	t.Log(weixinModel.SendText("oG9kZwrv8oFF9ja6WRlHxxMoJZoU", "哈哈哈，测试！"))
}
