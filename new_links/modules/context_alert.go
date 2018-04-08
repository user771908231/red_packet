package modules

func (this *Context)Success(msg string, jump_url string, time int) {
	this.alert(msg, "", jump_url, time)
}

func (this *Context)Error(msg string, jump_url string, time int) {
	this.alert("", msg, jump_url, time)
}

func (ctx *Context)alert(msg string, err string, jump_url string, time int) {
	ctx.Data["Msg"] = msg
	ctx.Data["Error"] = err
	ctx.Data["JumpUrl"] = jump_url
	ctx.Data["Time"] = time

	ctx.HTML(200, "common/alert")
	ctx.Resp.Write([]byte{})
}
