package weixinModel

import (
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"casino_admin/conf"
)

//发红包
func SendRedPack(open_id string, money float64, trade_no string) error {
	https_client, err := core.NewTLSHttpClient("../conf/apiclient_cert.pem", "../conf/apiclient_key.pem")
	if err != nil {
		return err
	}

	req_map := map[string]string{
		"nonce_str": bson.NewObjectId().Hex(),  //随机数
		"mch_billno": trade_no,  //订单号
		"mch_id": WX_MCH_ID,
		"wxappid": WX_APP_ID,
		"send_name": "神经棋牌",   //商户名称
		"re_openid": open_id,  //用户openid
		"total_amount": fmt.Sprintf("%.0f", money * 100),  //付款金额，单位（分）
		"total_num": "1",  //红包发放总人数
		"wishing": "神经棋牌祝您恭喜发财、大吉大利！",  //红包祝福语
		"client_ip": conf.Server.OutIp,  //调用接口的机器Ip地址
		"act_name": "玩游戏领红包",  //活动名称
		"remark": "红包兑换",  //备注信息
	}

	req_map["sign"] = core.Sign(req_map, WX_API_KEY, nil)

	client := core.NewClient(WX_APP_ID, WX_MCH_ID, WX_API_KEY, https_client)
	_,err = mmpaymkttransfers.SendRedPack(client, req_map)

	return err
}
