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
	"casino_common/proto/ddproto"
	"casino_common/utils/numUtils"
	"casino_common/utils/db"
	"errors"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/userModel"
	"casino_common/common/Error"
	"github.com/golang/protobuf/proto"
	"casino_common/common/model/wxpayDao"
	"casino_redpack/model/googsModel"
)

////旺实富支付接口相关方法
//	PAYWAP_USERCODE    = "5010206923"                       //旺实富分配的商户号
//	PAYWAP_COMPKEY     = "BBF056CFF745452292E3A2C9DEDBCD6B" //旺实富分配的密钥
//	PAYWAP_OFFICIALIP1 = "59.110.175.55"                    //旺实富官方ip
//	PAYWAP_OFFICIALIP2 = "59.110.159.71"                    //旺实富官方ip
//	PAYWAP_URL_FORMPAY = "http://pay.paywap.cn/form/pay"    //旺实富支付提交url
//
//	PAYWAP_URL_PAY    = "/weixin/paywap/pay"         //调起旺实富支付跳转
//	PAYWAP_URL_RETURN = "/weixin/paywap/return_page" //旺实富支付结果展示页面
//	PAYWAP_URL_NOTIFY = "/weixin/paywap/notify"      //旺实富支付结果回调页面 发货以此为准

//充值订单
type RechargeOrder struct {
	ObjId 			bson.ObjectId `bson:"_id"`		//订单ID
	UserId			uint32		//充值用户ID
	OrderNumber		string		//订单号
	OrderMoney		float64		//订单价格
	OrderTime		time.Time	//订单生成时间
	OrderStatus		int64		//订单状态	0 未支付 1 支付
	OrderType		int64		//订单类型	1 充值
	OrderGoods		string		//订单物品
	GoodsNunber		int64		//物品数量
	OrderDeleteStatus	int64	//是否删除 0 删除 1 未删除
}

//插入一个新订单
func (Order *RechargeOrder) Insert() error{
	Order.ObjId= bson.NewObjectId()
	Order.OrderTime = time.Now()
	Order.OrderStatus = int64(0)
	Order.OrderDeleteStatus = int64(1)
	err := db.C(tableName.TABLE_ORDER_LISTS).Insert(Order)
	return err
}

//更新订单信息
func (Order *RechargeOrder) Update() error{
	err := db.C(tableName.TABLE_ORDER_LISTS).Update(bson.M{"id": Order.ObjId}, Order)
	return err
}
//订单删除
func (Order *RechargeOrder) Delete()  error{
	Order.OrderDeleteStatus = int64(0)
	err := db.C(tableName.TABLE_ORDER_LISTS).Update(bson.M{"id": Order.ObjId}, Order)
	return err
}


//处理客户端提交支付金额 返回选择支付方式的页面
func PayWapPaymethodHandler(ctx *modules.Context) {
	comboid := ctx.Query("totalFee")
	//判断错误值 重新赋值
	if len(comboid) < 14 {
		comboid = "5abde8b56ca16d4822991b48"
	}
	//todo 根据comboid套餐信息得到money
	userid := ctx.Query("userid")
	var err error
	var userId int = 0
	userId, err = strconv.Atoi(userid)
	if err != nil {
		log.E("请求旺实富支付方式页面 参数错误 comboid[%v] userid[%v]", comboid, userid)
		ctx.Error("参数错误 code:-1", "", 0)
		return
	}

	//这里生成订单号 传给选择页面展示 新写的
	//order := googsModel.GetWxpayTradeNo(1, uint32(userId), comboid, time.Now())
	order := service.GetWxpayTradeNo(1, uint32(userId), int32(userId), time.Now())
	/**************************请求参数**************************/
	ctx.Data["p2_order"] = order
	//fmt.Println("id:",comboid,"userid:",userid,comboid,userId,order)
	//return
	//todo money 只保留小数点后两位 若没有小数 也要显示 如 50.00
	//获取商品价
	goods_info := googsModel.GetGoog(bson.ObjectIdHex(comboid))
	//goods_info := goodsRowDao.GetGoodsInfo(int32(comboId))
	if goods_info == nil {
		log.E("商品id(%d)不存在!", comboid)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	//生成充值的明细，此数据是要保存到数据库的
	 err = NewAndSavePayDetails(uint32(userId), comboid, int64(1), order, int64(goods_info.Number),goods_info.Price)
	if err != nil {
		log.E("订单插入数据库失败！err:%v", err)
		ctx.Error("参数错误 code:-2", "", 0)
		return
	}

	ctx.Data["p3_money"] = fmt.Sprintf("%.0f", goods_info.Price)
	ctx.Data["p6_ordertime"] = getOrderTime()
	ctx.Data["p14_customname"] = userid //终端客户

	//todo 根据套餐信息得到账单名 金币、钻石 and so on
	bill_name := "金币"
	//switch goods_info.GoodsType {
	//case ddproto.HallEnumTradeType_TRADE_COIN:
	//	bill_name = "金币"
	//case ddproto.HallEnumTradeType_TRADE_DIAMOND:
	//	bill_name = "钻石"
	//case ddproto.HallEnumTradeType_PROPS_FANGKA:
	//	bill_name = "房卡"
	//default:
	//	bill_name = "其他"
	//}
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

	//paywapIp := ctx.RemoteAddr()

	//暂时注释
	//if paywapIp != PAYWAP_OFFICIALIP1 && paywapIp != PAYWAP_OFFICIALIP2 {
		//log.E("PayWapNotifyHandler ip地址错误 未经验证的请求ip[%v] 官方ip1[%v] 官方ip2[%v]", paywapIp, PAYWAP_OFFICIALIP1, PAYWAP_OFFICIALIP2)
		//ctx.Error("参数错误 code:-1", "", 0)
		//return
	//}
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
	//err := service.DoAsynCb(p2_order, numUtils.String2Float64(p3_money))
	err := CheckOrder(p2_order, numUtils.String2Float64(p3_money))
	if err == nil {
		log.T("支付回调成功[%v:%v]！", p2_order, p3_money)
		////根据id获取用户
		//user := userModel.GetUserById(ctx.IsLogin().Id)
		//err := user.CapitalUplete("+",numUtils.String2Float64(p3_money))
		//if err != nil {
		//	//记录充值订单
		//	err := NewOrder(ctx.IsLogin().Id,p2_order, numUtils.String2Float64(p3_money))
		//	if err == nil {
		//
		//	}
		//	ctx.Error("回调 code:-5", "", 0)
		//
		//}
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

func GetOrderId(OrderNumber string) *RechargeOrder{
	var err error = nil
	Order_row := new(RechargeOrder)
	err = db.C(tableName.TABLE_ORDER_LISTS).Find(bson.M{
		"id": OrderNumber,
	}, Order_row)
	if err != nil {
		return Order_row
	}
	return nil
}
//生成一条订单支付记录
func NewOrder(userid uint32,OrderNumber string,numerical float64) error{
	data := RechargeOrder{
		UserId: userid,
		OrderNumber:OrderNumber,
		OrderMoney:numerical,
		OrderType:	int64(1),
		OrderGoods:"金币",
		GoodsNunber:int64(20),
	}
	err := data.Insert()
	if err == nil {
		return nil
	}
	return errors.New("生成订单支付记录失败")
}
//生成订单信息
func  GenerateOtder(OrderNumber string,total_fee float64) error{
	detail := service.GetDetailsByTradeNo(OrderNumber)
	if detail == nil {
		msg := fmt.Sprintf("没有在数据中找到订单号[%v]对应的套餐..", OrderNumber)
		log.E(msg)
		return errors.New(msg)
	}

	//判断是否是重复回调
	if detail.GetStatus() == ddproto.PayEnumTradeStatus_PAY_S_SUCC {
		log.E("tradeNo[%v]重复回调", OrderNumber)
		return nil
	}

	log.T("更新订单[%v]的回调信息，detail[%v]", OrderNumber, detail)

	//找到套餐
	meal := service.GetMealById(detail.GetProductId())
	User := userModel.GetUserById(detail.GetUserId())
	if User == nil {
		msg := fmt.Sprintf("没有在数据中找到用户ID：【%d】..", detail.GetUserId())
		log.E(msg)
		return errors.New(msg)
	}
	User.CapitalUplete("+",float64(meal.Amount),"充值")
	//更新订单状态
	service.UpdateDetailsStatus(OrderNumber, ddproto.PayEnumTradeStatus_PAY_S_SUCC)
	//保存订单到数据库...
	log.T("微信支付成功，为用户%d充值%d金币。", detail.GetUserId(), int64(meal.Amount))

	go func() {
		defer Error.ErrorRecovery("UpdateUserByMeal wxpayDao.UpsertDetail")
		detail.Status = ddproto.PayEnumTradeStatus_PAY_S_SUCC.Enum()
		detail.Money = proto.Float64(total_fee)
		wxpayDao.UpsertDetail(detail) //保存到数据库
		//DelDetails(tradeNo)           //保存到数据库之后删除//	app收到回复之后再删除
	}()
	return nil
}

//得到一个支付明细
func NewAndSavePayDetails(userId uint32, mealId string, payModelId int64, tradeNo string, diamond int64,Price float64)  error{
	R := new(RechargeOrder)
	R.UserId = userId
	R.OrderNumber = tradeNo
	R.OrderMoney = Price
	R.OrderType = payModelId
	R.GoodsNunber = diamond
	R.OrderGoods = mealId
	err := R.Insert()
	if err != nil {
		return err
	}
	return err
}

func CheckOrder(OrderNumber string,total_fee float64) error{
	//找到订单信息
	R := GetOrderId(OrderNumber)
	if R == nil {
		msg := fmt.Sprintf("没有在数据中找到订单号[%v]对应的套餐..", OrderNumber)
		log.E(msg)
		return errors.New(msg)
	}

	if R.OrderMoney != total_fee {
		msg := fmt.Sprintf("订单号[%v]对应的的金额与支付金额不一致[%s]", OrderNumber,total_fee )
		log.E(msg)
		return errors.New(msg)
	}

	//判断是否是重复回调
	if R.OrderStatus == 1 {
		log.E("tradeNo[%v]重复回调", OrderNumber)
		return nil
	}
	log.T("更新订单[%v]的回调信息，detail[%v]", OrderNumber, R)
	//找到套餐
	G := googsModel.GetGoog(bson.ObjectIdHex(R.OrderGoods))
	User := userModel.GetUserById(R.UserId)
	//更新订单状态
	if User == nil {
		msg := fmt.Sprintf("没有在数据中找到用户ID：【%d】..", R.UserId)
		log.E(msg)
		return errors.New(msg)
	}
	User.CapitalUplete("+",float64(G.Number),"充值")
	//更新订单状态
	R.OrderStatus = int64(1)
	err := R.Update()
	if err != nil {
		User.CapitalUplete("-",float64(G.Number),"")
		msg := fmt.Sprintf("更新订单tradeNo[%v]支付状态失败", OrderNumber)
		log.E(msg)
		return errors.New(msg)
	}
	log.T("微信支付成功，为用户%d充值%d金币。", User.Id, int64(G.Number))
	return nil
}

