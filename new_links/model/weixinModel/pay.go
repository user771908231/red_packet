package weixinModel

import (
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/chanxuehong/wechat.v2/mch/core"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"casino_common/common/model/agentModel"
	"casino_redpack/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"errors"
	"encoding/json"
	"casino_redpack/model/userModel"
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	payHandler core.Handler
	payServer  *core.Server
)

//支付初始化
func WxPayInit() {
	oauth2Endpoint = mpoauth2.NewEndpoint(WX_APP_ID, WX_APP_SECRET)
	payHandler = core.HandlerFunc(notifyHanlder)
	payServer = core.NewServer(WX_APP_ID, WX_MCH_ID, WX_API_KEY, payHandler, nil)
}

const (
	//回调地址
	WX_NOTIFY_URL string = "/mp/pay/callback"  //  http://wx.tondeen.com/mp/pay/callback
)

//网页微信支付的交易数据
type TradeData struct {
	AppId string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NoneStr string `json:"nonceStr"`
	PackageStr string `json:"package"`
	SignType string `json:"signType"`
	Sign string `json:"paySign"`
}

//统一下单
func GetUnifiedOrderResponse(tradeNo string, totalFee float64, detail string, user_ip string, user_openid string, host string) (*pay.UnifiedOrderResponse, error) {
	client := core.NewClient(WX_APP_ID, WX_MCH_ID, WX_API_KEY, nil)
	time_now := time.Now()
	non_str := bson.NewObjectId().Hex()
	request := pay.UnifiedOrderRequest{
		DeviceInfo: "WEB", // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
		NonceStr: non_str, // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
		Body: detail, // 商品或支付单简要描述
		Detail: "", // 商品名称明细列表
		Attach: "", // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
		OutTradeNo: tradeNo, // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
		FeeType: "CNY", // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
		TotalFee: int64(totalFee*100),  // 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: user_ip, // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TimeStart: core.FormatTime(time_now), // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
		TimeExpire: core.FormatTime(time_now.Add(600 * time.Minute)), // 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
		GoodsTag: "", // 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
		NotifyURL: "http://"+ host + WX_NOTIFY_URL, // 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
		TradeType: "JSAPI", // 取值如下：JSAPI，NATIVE，APP，详细说明见参数规定
		ProductId: "", // trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
		LimitPay: "", // no_credit--指定不能使用信用卡支付
		OpenId: user_openid, // rade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。
	}

	res,err := pay.UnifiedOrder2(client, &request)
	return res, err
}

//发起支付需要的数据
func GetTradeData(prepay_id string) *TradeData {
	time_str := fmt.Sprint(time.Now().Unix())
	package_str := "prepay_id=" + prepay_id
	non_str := bson.NewObjectId().Hex()
	sign := core.JsapiSign(WX_APP_ID, time_str, non_str, package_str, "MD5", WX_API_KEY)
	tradeData := TradeData{
		AppId: WX_APP_ID,
		TimeStamp: time_str,
		NoneStr: non_str,
		PackageStr: package_str,
		SignType: "MD5",
		Sign: sign,
	}
	return &tradeData
}


//异步通知处理
func notifyHanlder(ctx *core.Context) {
	fmt.Printf("收到文本消息:\n%s\n", ctx.Msg)
	//收到的信息经过包装，验证sign 已经在框架做过了,这里只需要处理业务就行了
	out_trade_no := ctx.Msg["out_trade_no"] //商户订单号
	//开始处理回调
	err := agentModel.SetRechargeDone(out_trade_no)
	if err != nil {
		//log.E("微信异步回调出错..err :%v", err)
	} else {
		//正常情况下...回复微信...已经接收到了.
		msg := make(map[string]string)
		msg["return_code"] = "SUCCESS"
		ctx.Response(msg)
	}
}

//异步通知
func WxNotifyHandler(r *http.Request, w http.ResponseWriter) {
	payServer.ServeHTTP(w, r, nil)
}
//提现
type Withdrawals struct {
	ObjId 			bson.ObjectId `bson:"_id"`
	UserId 			uint32	//user id
	Number			float64 //提现数量
	Time 			time.Time	//	时间
	Status  		int			//状态 0 未受理 1 受理 2 拒绝
	AcceptanceID 	uint32		//处理人ID
	DeleteStatus 	int			//删除状态 0 未删除 1 删除
}

func (W *Withdrawals) Insert() error{
	//id,_ := db.GetNextIncrementID(config.WITHDRAWALS_KEY_ID,consts.RKEY_WITHSRAWALS_ID_KEY)
	W.ObjId = bson.NewObjectId()
	W.Time = time.Now()
	W.Status = int(0)
	W.DeleteStatus = 0
	err := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Insert(W)
	if err != nil {
		return errors.New("插入一条记录失败！")
	}
	return  nil

}

func WithdrawalsHandler(ctx *modules.Context)  {
	res := bson.M{
		"code": 0,
		"message": "fail",
		"request": bson.M{},
		"msg" : "申请提现失败！请联系客服！",
	}

	val := ctx.QueryFloat64("totalFee")
	Data := Withdrawals{
		UserId:ctx.IsLogin().Id,
		Number:val,
	}
	err := Data.Insert()
	if err == nil {
		res["code"] = 1
		res["message"] = "success"
		data,_ := json.Marshal(res)
		ctx.Write([]byte(data))
	}else{
		data,_ := json.Marshal(res)
		ctx.Write([]byte(data))
	}

}

func (Withdrawals *Withdrawals) UpdateStatus(status int,AcceptanceID uint32) error{
	Withdrawals.Status = status
	Withdrawals.AcceptanceID = AcceptanceID
	err := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Update(bson.M{"_id": Withdrawals.ObjId},Withdrawals)
	return err
}
func (Withdrawals *Withdrawals) Delete(status int,AcceptanceID 	uint32) error{
	Withdrawals.DeleteStatus = status
	Withdrawals.AcceptanceID = AcceptanceID
	err := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Update(bson.M{"_id": Withdrawals.ObjId},Withdrawals)
	return err
}

func GetWithdrawalsId(id bson.ObjectId) *Withdrawals {
	Withdrawals := new(Withdrawals)
	err := db.C(tableName.TABLE_WITHDRAWALS_LISTS).Find(bson.M{"_id":id},Withdrawals)
	if err != nil {
		return nil
	}
	return Withdrawals
}
var withdrawals *Withdrawals
func GetReady(w *Withdrawals) {
	withdrawals = w
}
func GetOver()  {
	withdrawals = nil
}

func Implement() error{

	if (withdrawals.AcceptanceID != uint32(0)) && (withdrawals.Status == 1) {
		User := userModel.GetUserById(withdrawals.UserId)
		if User != nil {
			err := User.CapitalUplete("-",withdrawals.Number)
			if err != nil {
				GetOver()
				return err
			}
			GetOver()
			return nil
		}
		GetOver()
		return nil
	}
	GetOver()
	return nil
}