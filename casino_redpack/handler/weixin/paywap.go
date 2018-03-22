package weixin

import (
	"casino_redpack/modules"
	"time"
	"fmt"
	"casino_common/common/service/wxservice"
	"strings"
	"strconv"
	"casino_common/common/log"
	"crypto/md5"
	"encoding/hex"
	"casino_common/common/model/goodsRowDao"
	"casino_common/proto/ddproto"
	"casino_common/utils/numUtils"
	"casino_redpack/handler/redpack"
	"net/url"
	"net/http"
)

//旺实富支付接口相关方法

const (
	PAYWAP_USERCODE    = "5010206923"                       //旺实富分配的商户号
	PAYWAP_COMPKEY     = "BBF056CFF745452292E3A2C9DEDBCD6B" //旺实富分配的密钥
	PAYWAP_OFFICIALIP1 = "59.110.175.55"                    //旺实富官方ip
	PAYWAP_OFFICIALIP2 = "59.110.159.71"                    //旺实富官方ip
	PAYWAP_URL_FORMPAY = "http://pay.paywap.cn/form/pay"    //旺实富支付提交url

	PAYWAP_URL_PAY    = "/weixin/paywap/pay"         //调起旺实富支付跳转
	PAYWAP_URL_RETURN = "/weixin/paywap/return_page" //旺实富支付结果展示页面
	PAYWAP_URL_NOTIFY = "/weixin/paywap/notify"      //旺实富支付结果回调页面 发货以此为准

)

//type PayWapKey struct {
//	p1_usercode				string		//旺实富分配的商户号 必填
//	p2_order				string		//(29)用户订单号，建议商户号+14 位时 间 yyyymmddhhmmss+5 位流水 号，中间用“-”分隔。例如: 12345678-20150728132430-12 345。 (只是建议，商户的订单号也可以 不采用这种格式) 必填
//	p3_money				string		//订单金额，精确到分。例如 99.99 必填
//	p4_returnurl			string		//(190) 用户明文跳转地址，用于告知付款 人支付结果。必须包含 http://或 https://。 必填
//	p5_notifyurl			string		//(190)服务器异步通知地址，用于通知商 户系统处理业务(数据库更新等)。 必须包含 http://或 https://。 必填
//	p6_ordertime			string		//商户订单时间，格式 yyyymmddhhmmss 必填
//	p7_sign					string		//(256)商户传递参数加密值，约定 p1_ usercode + "&" +p2_ order + 必填"&" +p3_ money + "&" +p4_returnurl + "&" +p5_notifyurl+ "&" +p6_ordertime +CompKey 连 接起来进行 MD5 加密后 32 位大 写字符串，(参数之间必须添加& 符号，最后 p6_ordertime 和 CompKey 之间不加&符号。 CompKey 为商户的秘钥)目前只 限定 md5 加密。
//	p8_signtype				string		//(5)签名验证方式:1、MD5，传固定 值 1。如果用户传递参数为空，则 默认 MD5 验证。 可空
//	p9_paymethod			string		//(5)商户支付方式:固定值 3。如果用 户传递参数为空，则默认网银支付。 可空
//	p10_paychannelnum		string		//(12)支付通道编码。可空
//	p11_cardtype			string		//(5)为空 可空
//	p12_channel				string		//(5)为空 可空
//	p13_orderfailertime		string		//(14)订单失效时间，格式为 14 位时间格 式，精确到秒 yyyymmddhhmmss。超时则此订 单失效。可空
//	p14_customname			string		//(128)客户、或者玩家所在平台账号。请务必填写真实信息，否则将影响后续查单结果。必填
//	p15_customcontacttype	string		//(10)客户联系方式类型:1、email，2、 phone，3、地址	可空
//	p16_customcontact		string		//(200)	客户联系方式 可空
//	p17_customip			string		//(128) 客户 ip 地址，规定以 192_168_0_253 格式，如果以 “192.168.0.253”可能会发生签名 错误。 必填
//	p18_product				string		//(256)\商品名称 可空
//	p19_productcat			string		//(200)商品种类可空
//	p20_productnum			string		//(10)商品数量，不传递参数默认 0	可空
//	p21_pdesc				string		//(200)商品描述	可空
//	p22_version				string		//(5)接口版本，目前默认 2.0可空
//	p23_charset				string		//(5)提交的编码格式，1、UTF-8，2、 GBK/GB2312，默认 UTF-8可空
//	p24_remark				string		//(256)备注。此参数我们会在下行过程中原样返回。您可以在此参数中记录一些数据，方便在下行过程中直接读取。可空
//	p25_terminal			string		//(5终端设备固定值 2。 必填
//	p26_iswappay			string		//(5)支付场景固定值 3。必填
//}

//充值订单
type RechargeOrder struct {
	Id				int64		//订单ID
	UserId			uint32		//充值用户ID
	OrderNumber		string		//订单号
	OrderMoney		float64		//订单价格
	OrderTime		time.Time	//订单生成时间
	OrderStatus		int64		//订单状态
	OrderType		int64		//订单类型
	OrderGoods		int			//订单物品
	OrderDeleteStatus	int64
}

func PayWapMoney(numerical_value float64,ctx *modules.Context) bool{
	val := Paywap(numerical_value,redpack.User_info(ctx).Id,ctx)
	err := Post(val)
	if err == nil {
		return true
	}
	return false
}

func Paywap(numerical float64,userId uint32,ctx *modules.Context) url.Values{
	data := make(url.Values)
	data["p1_usercode"] = []string{PAYWAP_USERCODE}
	data["p2_order"] = []string{service.GetWxpayTradeNo(1, uint32(userId), int32(numerical), time.Now())}
	data["p3_money"] = []string{fmt.Sprintf("%.0f", numerical)}
	data["p4_returnurl"] = []string{"http://" + HOST_IP + PAYWAP_URL_RETURN}
	log.T("旺实富支付结果展示页面地址：%d",data["p4_returnurl"])
	data["p5_notifyurl"] = []string{"http://" + HOST_IP + PAYWAP_URL_NOTIFY}
	log.T("旺实富支付回调地址：%d",data["p5_notifyurl"])
	data["p6_ordertime"] = []string{getOrderTime()}
	mixSignString := fmt.Sprintf("%s&%s&%s&%s&%s&%s%s",data["p1_usercode"],data["p2_order"],data["p3_money"],data["p4_returnurl"],data["p5_notifyurl"],data["p6_ordertime"], PAYWAP_COMPKEY)
	data["p7_sign"] = []string{MD5M(mixSignString)}
	data["p8_signtype"] = []string{"1"}
	data["p9_paymethod"] = []string{"3"}
	data["p10_paychannelnum"] = []string{""}
	data["p11_cardtype"] = []string{""}
	data["p12_channel"] = []string{""}
	data["p13_orderfailertime"] = []string{PAYWAP_ORDER_INVALID_TIME}
	data["p14_customname"] = []string{string(userId)}
	data["p15_customcontacttype"] = []string{""}
	data["p16_customcontact"] = []string{""}
	data["p17_customip"] = []string{"192_168_0_253"}
	data["p18_product"] = []string{""}
	data["p19_productcat"] = []string{"1"}
	data["p20_productnum"] = []string{""}
	data["p21_pdesc"] = []string{"充值"}
	data["p22_version"] =[]string{""}
	data["p23_charset"] =[]string{"UTF-8"}
	data["p24_remark"] = []string{""}
	data["p25_terminal"] = []string{"2"}
	data["p26_iswappay"] =[]string{"3"}
	return data
}

func Post(data url.Values) error{
	//把post表单发送给目标服务器
	res, err := http.PostForm(PAYWAP_URL_FORMPAY, data)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func MD5M(mixSignString string) string {
	h := md5.New()
	h.Write([]byte(mixSignString)) // 需要加密的字符串为 123456
	sign := hex.EncodeToString(h.Sum(nil))

	log.T("转大写前的签名字符串:[%v]", sign)
	return strings.ToUpper(sign)
}

//处理客户端提交支付金额 返回选择支付方式的页面
func PayWapPaymethodHandler(ctx *modules.Context) {
	comboid := ctx.Query("comboid")
	//todo 根据comboid套餐信息得到money
	userid := ctx.Query("userid")
	var err error
	var comboId int = 0
	var userId int = 0
	comboId, err = strconv.Atoi(comboid)
	userId, err = strconv.Atoi(userid)
	if err != nil {
		log.E("请求旺实富支付方式页面 参数错误 comboid[%v] userid[%v]", comboid, userid)
		ctx.Error("参数错误 code:-1", "", 0)
		return
	}

	//这里生成订单号 传给选择页面展示
	order := service.GetWxpayTradeNo(1, uint32(userId), int32(comboId), time.Now())
	/**************************请求参数**************************/
	ctx.Data["p2_order"] = order

	//todo money 只保留小数点后两位 若没有小数 也要显示 如 50.00
	goods_info := goodsRowDao.GetGoodsInfo(int32(comboId))
	if goods_info == nil {
		log.E("商品id(%d)不存在!", comboid)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	//生成充值的明细，此数据是要保存到数据库的
	_, err = service.NewAndSavePayDetails(uint32(userId), int32(comboId), 1, order, int64(goods_info.Amount))
	if err != nil {
		log.E("订单插入数据库失败！err:%v", err)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	ctx.Data["p3_money"] = fmt.Sprintf("%.0f", goods_info.Price)
	ctx.Data["p6_ordertime"] = getOrderTime()
	ctx.Data["p14_customname"] = userid //终端客户

	//todo 根据套餐信息得到账单名 金币、钻石 and so on
	bill_name := "钻石"
	switch goods_info.GoodsType {
	case ddproto.HallEnumTradeType_TRADE_COIN:
		bill_name = "金币"
	case ddproto.HallEnumTradeType_TRADE_DIAMOND:
		bill_name = "钻石"
	case ddproto.HallEnumTradeType_PROPS_FANGKA:
		bill_name = "房卡"
	default:
		bill_name = "其他"
	}
	ctx.Data["bill_name"] = bill_name
	ctx.Data["pay_url"] = "http://" + ctx.Req.Host + PAYWAP_URL_PAY

	ctx.HTML(200, "weixin/paywap/paymethod")
}

//客户端提交支付请求页面
func PayWapPayHandler(ctx *modules.Context) {
	p2_order := ctx.Query("p2_order")
	p3_money := ctx.Query("p3_money")
	p6_ordertime := ctx.Query("p6_ordertime")
	p9_paymethod := ctx.Query("p9_paymethod")
	p14_customname := ctx.Query("p14_customname")

	if p2_order == "" {
		log.E("请求旺实富支付页面 order参数错误 ")
		ctx.Error("参数错误 code:-1", "", 0)
		return
	}
	if p3_money == "" {
		log.E("请求旺实富支付页面 money参数错误 ")
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}
	if p6_ordertime == "" {
		log.E("请求旺实富支付页面 ordertime参数错误 ")
		ctx.Error("参数错误 code:-3", "", 0)
		return
	}
	if p9_paymethod == "" {
		log.E("请求旺实富支付页面 paymethod参数错误 ")
		ctx.Error("参数错误 code:-4", "", 0)
		return
	}
	if p14_customname == "" {
		log.E("请求旺实富支付页面 customname参数错误 ")
		ctx.Error("参数错误 code:-5", "", 0)
		return
	}

	/**************************请求参数**************************/
	ctx.Data["p1_usercode"] = PAYWAP_USERCODE //旺实富分配的商户号
	ctx.Data["compkey"] = PAYWAP_COMPKEY      //旺实富分配的密钥
	ctx.Data["p2_order"] = p2_order
	ctx.Data["p3_money"] = p3_money
	ctx.Data["p4_returnurl"] = "http://" + ctx.Req.Host + PAYWAP_URL_RETURN
	ctx.Data["p5_notifyurl"] = "http://" + ctx.Req.Host + PAYWAP_URL_NOTIFY
	ctx.Data["p6_ordertime"] = getOrderTime()
	ctx.Data["p14_customname"] = p14_customname //终端客户
	ctx.Data["p17_customip"] = "192_168_0_253"  //提交ip 需改为自动获取
	ctx.Data["p25_terminal"] = "2"
	ctx.Data["p26_iswappay"] = "3"
	mixSignString := fmt.Sprintf("%s&%s&%s&%s&%s&%s%s", ctx.Data["p1_usercode"], ctx.Data["p2_order"], ctx.Data["p3_money"], ctx.Data["p4_returnurl"], ctx.Data["p5_notifyurl"], ctx.Data["p6_ordertime"], PAYWAP_COMPKEY)
	log.T("md5前的签名字符串:[%v]", mixSignString)

	h := md5.New()
	h.Write([]byte(mixSignString)) // 需要加密的字符串为 123456
	sign := hex.EncodeToString(h.Sum(nil))

	log.T("转大写前的签名字符串:[%v]", sign)
	ctx.Data["p7_sign"] = strings.ToUpper(sign)
	log.T("md5后的签名字符串:[%v]", ctx.Data["p7_sign"])
	ctx.Data["p9_paymethod"] = p9_paymethod //支付方式 3微信 4支付宝
	ctx.Data["formpay_url"] = PAYWAP_URL_FORMPAY

	ctx.HTML(200, "weixin/paywap/pay")
}

//旺实富支付结果返回页面
//下行第一步:旺实富微信平台将支付结果传递给 p4_returnurl(用户在上行过程中提
//交的参数),此部分用于在付款人浏览器中显示支付结果,传递方式为 get。
func PayWapReturnPageHandler(ctx *modules.Context) {
	paywapIp := ctx.RemoteAddr()

	if paywapIp != PAYWAP_OFFICIALIP1 && paywapIp != PAYWAP_OFFICIALIP2 {
		log.E("PayWapReturnPageHandler ip地址错误 未经验证的请求ip[%v] 官方ip1[%v] 官方ip2[%v]", paywapIp, PAYWAP_OFFICIALIP1, PAYWAP_OFFICIALIP2)
		ctx.Error("参数错误 code:-1", "", 0)
		return
	}
	p1_usercode := ctx.Query("p1_usercode")
	compkey := ctx.Query("compkey")
	p2_order := ctx.Query("p2_order")
	p3_money := ctx.Query("p3_money")
	p4_returnurl := ctx.Query("p4_returnurl")
	p5_notifyurl := ctx.Query("p5_notifyurl")
	p6_ordertime := ctx.Query("p6_ordertime")
	p7_sign := ctx.Query("p7_sign")
	p9_paymethod := ctx.Query("p9_paymethod")
	p14_customname := ctx.Query("p14_customname")
	p17_customip := ctx.Query("p17_customip")
	p25_terminal := ctx.Query("p25_terminal")
	p26_iswappay := ctx.Query("p26_iswappay")

	log.T("支付返回的结果 "+
		"p1_usercode[%s] compkey[%v] \n"+
		"p2_order[%v] p3_money[%v] \n"+
		"p4_returnurl[%v] p5_notifyurl[%v] \n"+
		"p6_ordertime[%v] p7_sign[%v] \n"+
		"p9_paymethod[%v] p14_customname[%v] \n"+
		"p17_customip[%v] p25_terminal[%v] p26_iswappay[%v]",
		p1_usercode, compkey,
		p2_order, p3_money,
		p4_returnurl, p5_notifyurl,
		p6_ordertime, p7_sign,
		p9_paymethod, p14_customname,
		p17_customip, p25_terminal, p26_iswappay)

	//校验
	if p1_usercode != PAYWAP_USERCODE {
		log.E("商户号错误 请求的商户号[%v] 设置的商户号[%v]", p1_usercode, PAYWAP_USERCODE)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	if compkey != PAYWAP_COMPKEY {
		log.E("商户秘钥错误 请求的商户秘钥[%v] 设置的商户秘钥[%v]", compkey, PAYWAP_COMPKEY)
		ctx.Error("参数错误 code:-3", "", 0)
		return
	}

	//根据参数计算出md5签名 并与请求的签名做匹配
	mixSignString := fmt.Sprintf("%s&%s&%s&%s&%s&%s%s", p1_usercode, p2_order, p3_money, p4_returnurl, p5_notifyurl, p6_ordertime, compkey)
	log.T("md5前的签名字符串:[%v]", mixSignString)

	h := md5.New()
	h.Write([]byte(mixSignString)) // 需要加密的字符串为 123456
	sign := hex.EncodeToString(h.Sum(nil))

	log.T("转大写前的签名字符串:[%v]", sign)
	sign = strings.ToUpper(sign)
	log.T("md5后的签名字符串:[%v]", sign)

	if sign != p7_sign {
		log.E("参数签名错误 请求参数的签名[%v] 计算出的签名[%v]", p7_sign, sign)
		ctx.Error("参数错误 code:-4", "", 0)
		return
	}

	ctx.HTML(200, "weixin/paywap/return_page")
}

//旺实富支付结果异步回调
//下行第二步:旺实富微信支付平台将支付结果传递给 p5_notifyurl(用户在上行过程 中提交的参数),此部分用于通知商户的系统处理业务(包括数据库更新,在系统
//中为付款人增加虚拟货币等),传递方式为 post。
func PayWapNotifyHandler(ctx *modules.Context) {
	paywapIp := ctx.RemoteAddr()

	if paywapIp != PAYWAP_OFFICIALIP1 && paywapIp != PAYWAP_OFFICIALIP2 {
		log.E("PayWapNotifyHandler ip地址错误 未经验证的请求ip[%v] 官方ip1[%v] 官方ip2[%v]", paywapIp, PAYWAP_OFFICIALIP1, PAYWAP_OFFICIALIP2)
		ctx.Error("参数错误 code:-1", "", 0)
		return
	}
	p1_usercode := ctx.Query("p1_usercode")
	compkey := ctx.Query("CompKey")
	p2_order := ctx.Query("p2_order")
	p3_money := ctx.Query("p3_money")
	p4_status := ctx.Query("p4_status")
	p5_payorder := ctx.Query("p5_payorder")
	p6_paymethod := ctx.Query("p6_paymethod")
	p7_sign := ctx.Query("p7_sign")
	p8_charset := ctx.Query("p8_charset")
	p9_signtype := ctx.Query("p9_signtype")
	CompKey := ctx.Query("CompKey")
	p25_terminal := ctx.Query("p25_terminal")
	p26_iswappay := ctx.Query("p26_iswappay")
	p10_sign := ctx.Query("p10_sign")

	log.T("支付异步回调 "+
		"p1_usercode[%s] compkey[%v] \n"+
		"p2_order[%v] p3_money[%v] \n"+
		"p4_status[%v] p5_payorder[%v] \n"+
		"p6_ordertime[%v] p7_sign[%v] \n"+
		"p8_charset[%v] p9_signtype[%v] \n"+
		"CompKey[%v] p25_terminal[%v] p26_iswappay[%v] p10_sign[%v]",
		p1_usercode, compkey,
		p2_order, p3_money,
		p4_status, p5_payorder,
		p6_paymethod, p7_sign,
		p8_charset, p9_signtype,
		CompKey, p25_terminal, p26_iswappay, p10_sign)

	//校验
	if p1_usercode != PAYWAP_USERCODE {
		log.E("商户号错误 请求的商户号[%v] 设置的商户号[%v]", p1_usercode, PAYWAP_USERCODE)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	//if compkey != PAYWAP_COMPKEY {
	//	log.E("商户秘钥错误 请求的商户秘钥[%v] 设置的商户秘钥[%v]", compkey, PAYWAP_COMPKEY)
	//	ctx.Error("参数错误 code:-3", "", 0)
	//	return
	//}

	//根据参数计算出md5签名 并与请求的签名做匹配
	mixSignString := fmt.Sprintf("%s&%s&%s&%s&%s&%s&&%s&%s&%s", p1_usercode, p2_order, p3_money, p4_status, p5_payorder, p6_paymethod, p8_charset, p9_signtype, PAYWAP_COMPKEY)
	log.T("md5前的签名字符串:[%v]", mixSignString)

	h := md5.New()
	h.Write([]byte(mixSignString)) // 需要加密的字符串为 123456
	sign := hex.EncodeToString(h.Sum(nil))

	log.T("转大写前的签名字符串:[%v]", sign)
	sign = strings.ToUpper(sign)
	log.T("md5后的签名字符串:[%v]", sign)

	if sign != p10_sign {
		log.E("参数签名错误 请求参数的签名[%v] 计算出的签名[%v]", p10_sign, sign)
		ctx.Error("参数错误 code:-4", "", 0)
		return
	}

	//异步回调不需要返回页面
	//todo 增加货币
	err := service.DoAsynCb(p2_order, numUtils.String2Float64(p3_money))
	if err == nil {
		log.T("支付回调成功[%v:%v]！", p2_order, p3_money)
		ctx.Write([]byte("success"))
	}else {
		log.E("支付回调失败[%v:%v] err:%v！", p2_order, p3_money, err)
		ctx.Error("回调 code:-5", "", 0)
	}
	return
}

func getOrderTime() string {
	now := time.Now()
	year, month, day := now.Date()
	return fmt.Sprintf("%d%d%d%d%d%d", year, int(month), day, now.Hour(), now.Minute(), now.Second())
}
