package weixin

import (
	"testing"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers"
	"casino_super/model/agentModel"
)

//微信红包测试
func TestRedBag(t *testing.T) {
	https_client, err := core.NewTLSHttpClient("../conf/apiclient_cert.pem", "")
	if err != nil {
		t.Log("https err:", err)
		return
	}
	req_map := map[string]string{

	}
	client := core.NewClient(agentModel.WX_APP_ID,agentModel.WX_MCH_ID, agentModel.WX_API_KEY, https_client)
	resp,err := mmpaymkttransfers.SendRedPack(client, req_map)
	if err != nil {
		t.Log("resp err:", err)
		return
	}
	t.Log("resp:", resp)
}
