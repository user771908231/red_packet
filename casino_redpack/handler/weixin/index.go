package weixin

import (
	"casino_redpack/modules"
	"strings"
	"casino_redpack/conf"
)

func MainHandler(ctx *modules.Context) {
	channelid := ctx.Query("channelid")

	show_add_agent, show_apply_agent := false, false

	switch channelid {
	case "61":
			show_add_agent = false
			show_apply_agent = true
	default:
		show_add_agent = true
		show_apply_agent = false
	}

	ctx.Data["show_add_agent"] = show_add_agent
	ctx.Data["show_apply_agent"] = show_apply_agent

	if strings.Contains(conf.Server.SiteName, "来一圈") {
		ctx.Data["site_logo"] = "indexlogo.jpg"
	}else if strings.Contains(conf.Server.SiteName, "里来") {
		ctx.Data["site_logo"] = "indexlogo_lilai.png"
	}

	ctx.HTML(200, "weixin/agent/index")
}
