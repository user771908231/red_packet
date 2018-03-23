package weixin

//旺实富支付接口相关方法
const (
	PAYWAP_USERCODE    = "5010206923"                       //旺实富分配的商户号
	PAYWAP_COMPKEY     = "BBF056CFF745452292E3A2C9DEDBCD6B" //旺实富分配的密钥
	PAYWAP_OFFICIALIP1 = "59.110.175.55"                    //旺实富官方ip
	PAYWAP_OFFICIALIP2 = "59.110.159.71"                    //旺实富官方ip
	PAYWAP_URL_FORMPAY = "http://pay.paywap.cn/form/pay"    //旺实富支付提交url

	PAYWAP_URL_PAY    = "/weixin/paywap/pay"         		//调起旺实富支付跳转
	PAYWAP_URL_RETURN = "/weixin/paywap/return_page" 		//旺实富支付结果展示页面
	PAYWAP_URL_NOTIFY = "/weixin/paywap/notify"      		//旺实富支付结果回调页面 发货以此为准
	PAYWAP_ORDER_INVALID_TIME = "" 							//订单失效时间
	HOST_IP = "182.150.164.207:9091"
	PAYWAP_RETURN_URL = "/home/member/recharge/confirm"		//用户确认返回


)
