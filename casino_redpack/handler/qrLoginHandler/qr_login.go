package qrLoginHandler

import (
	"casino_redpack/modules"
	"casino_common/utils/redisUtils"
	"casino_common/common/consts"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
)

//登陆成功
func QrLoginHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	code := ctx.Query("code")
	key := consts.RKEY_QR_CODE+ "_" +code

	redis_wxinfo := &ddproto.WeixinInfo{
		OpenId:   &wx_info.OpenId,
		NickName: &wx_info.Nickname,
		HeadUrl:  &wx_info.HeadImageURL,
		Sex:      proto.Int32(int32(wx_info.Sex)),
		City:     &wx_info.City,
		UnionId:  &wx_info.UnionId,
	}

	err := redisUtils.SetObj(key, redis_wxinfo)

	if err != nil {
		ctx.Error("登陆失败！请联系管理员。", "", 0)
		return
	}
	ctx.Success("恭喜你，登陆成功！", "", 0)
}
