package weixinModel

import (
	"github.com/chanxuehong/wechat.v2/mp/message/custom"
	"github.com/chanxuehong/wechat.v2/mp/core"
)

//发送消息
func SendText(open_id string, content string) error {
	text := custom.NewText(open_id, content, "")
	new_access_token_server := core.NewDefaultAccessTokenServer(WX_APP_ID, WX_APP_SECRET, nil)
	client := core.NewClient(new_access_token_server, nil)
	return custom.Send(client, text)
}
