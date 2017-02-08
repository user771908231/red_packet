package modules

import "gopkg.in/macaron.v1"

type Alert struct {
	Contex *macaron.Context
}

func (a Alert)Success(msg string, jump_url string, time int) {
	alert(a.Contex, msg, "", jump_url, time)
}

func (a Alert)Error(msg string, jump_url string, time int) {
	alert(a.Contex, "", msg, jump_url, time)
}

func alert(ctx *macaron.Context, msg string, err string, jump_url string, time int) {
	ctx.Data["Msg"] = msg
	ctx.Data["Error"] = err
	ctx.Data["JumpUrl"] = jump_url
	ctx.Data["Time"] = time

	ctx.HTML(200, "common/alert")
	ctx.Resp.Write([]byte{})
}
