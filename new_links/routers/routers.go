package routers

import (
	"gopkg.in/macaron.v1"
	"casino_redpack/handler/admin"
	"github.com/go-macaron/binding"
	"casino_redpack/modules"
	"casino_redpack/handler/weixin"
	"casino_redpack/handler/qrLoginHandler"
	"casino_redpack/model/agentProModel"
	"casino_redpack/handler/agentPro"
	"casino_redpack/model/weixinModel"
	"casino_redpack/handler/redpack"
	"casino_redpack/model/redModel"
	"casino_redpack/handler/pay"
	"casino_redpack/handler/admin/manage"
	"casino_redpack/handler/admin/gooods"
	"casino_redpack/model/googsModel"
)

//注册路由
func Regist(m *macaron.Macaron) {
	m.Group("", func() {
		//微信
		m.Group("/weixin", func() {
			m.Get("/oauth/login", weixinModel.OauthLogin)
			m.Get("/oauth/callback", weixinModel.OauthCallBack)
			//需要微信登录
			m.Group("/agent", func() {
				m.Get("/", weixin.MainHandler)
				//新增代理
				m.Get("/add", weixin.AgentAddHandler)
				m.Post("/add", binding.Bind(weixin.AddAgentForm{}), weixin.AgentAddPostHandler)
				//代理个人信息
				m.Get("/info", weixin.AgentInfoHandler)
				//充值
				m.Group("/recharge", func() {
					m.Get("/", weixin.RechargeListHandler)
					m.Get("/done", weixin.RechargeDoneHandler)
					m.Get("/wx_pay", weixin.RechargeAjaxWxTradeDataHandler)
					m.Get("/log", weixin.RechargeLogHandler)
				}, weixin.NeedIsRootAgent)
				//销售
				m.Group("/sales", func() {
					m.Get("/", weixin.SalesIndexHandler)
					m.Post("/",binding.BindIgnErr(weixin.SalesForm{}), weixin.SalesToUserHandler)
					m.Get("/log", weixin.SalesLogHandler)
				})
				//我的客户
				m.Get("/customers", weixin.CustomersListHandler)
				//我的-金币场客户
				m.Get("/customers_coin", weixin.CustomersCoinListHandler)

				//代理审核
				m.Group("/apply", func() {
					m.Get("/log", weixin.ApplyLogHandler)
					m.Get("/switch", weixin.ApplySwitchState)
				})

				//我的下线代理
				m.Get("/my_agents", weixin.MyAgentsHandler)
				//返利记录
				m.Group("/rebate", func() {
					m.Get("/log", weixin.RebateLogHandler)
					m.Get("/check", weixin.CheckRebateHandler)
				})

				//房费返利
				m.Group("/coin_fee", func() {
					//返利记录
					m.Get("/rebate_log", weixin.CoinFeeRebateLogHandler)
					//提现post
					m.Get("/withdraw_post", weixin.CoinFeeRebateWithdrawPost)
				})

				////登录-登出
				//m.Group("/user", func() {
				//	m.Get("/login", func(ctx *modules.Context) {
				//		ctx.Redirect("/admin/login_2", 200)
				//	})
				//	m.Get("/logout", func(ctx *modules.Context) {
				//		ctx.Session.Delete("wx_user")
				//		ctx.Success("退出成功！", "/weixin/", 3)
				//	})
				//})
				//俱乐部管理
				m.Group("/group", func() {
					//代理俱乐部列表
					m.Get("/list", weixin.GroupListHandler)
					//俱乐部编辑
					m.Get("/edit", weixin.GroupEditHandler)
					//俱乐部编辑
					m.Get("/edit_post", weixin.GroupEditPostHandler)
				})
			}, weixin.NeedWxLogin, weixin.NeedIsAgent)

			//二维码登陆
			m.Group("/game", func() {
				m.Get("/qrlogin", qrLoginHandler.QrLoginHandler)
			}, weixin.NeedWxLogin)

		})

		//管理后台
		m.Group("/admin", func() {
			//主页
			m.Get("/", admin.IndexHandler)
			m.Get("/logout", admin.LoginOutHandler)
			//管理
			m.Group("/manage", func() {
				//红包兑换相关
				m.Group("/exchange", func() {
					m.Get("/", manage.ExchangeListHandler)
					m.Get("/switch", manage.ExchangeSwitchState)
				})
				m.Group("/Withdrawals", func() {
					m.Get("/",manage.WithdrawalsHandle)
					m.Get("/operation",manage.WithdrawalsOperationHandle)
				})

				m.Group("/goods", func() {
					m.Get("/",goods.IndexHandler)
					m.Get("/add",goods.AddHandler)
					m.Post("/add",binding.BindIgnErr(googsModel.GoogsForm{}),goods.AddPost)
					m.Get("/update",goods.UpdateHandler)
					m.Get("/operation",goods.OperationHandler)
				})
				m.Group("/links", func() {
					m.Get("/",)
				})

				m.Get("/postal",manage.PostalHandle)
			})
		}, admin.NeedLogin(2))
		m.Post("/admin/login", admin.NeedCaptcha, binding.Bind(admin.LoginForm{}), admin.LoginPostHandler)
		//管理登录
		m.Get("/admin/login", admin.LoginHandler)
		//红包项目
		m.Group("/home", func() {
			//首页
			m.Get("/", redpack.HomeHandler)
			//退出地址
			m.Get("/outlogin", admin.LoginOutHandler)
			m.Group("/member", func() {
				//充值
				m.Get("/recharge", redpack.RechargeHandler)
				//充值订单
				m.Post("/recharge",pay.PostRechargeHandler)
				//充值确认
				//m.Get("/recharge/confirm",redpack.RechargeConfirmHandler)
				//提现
				m.Get("/outred", redpack.OutredHandler)
				//游戏记录
				m.Get("/record", redpack.RecordHandler)

				//绑定银行卡
				m.Get("/yhk", redpack.BindBankHandler)
				//银行卡列表
				m.Get("/bank_list", redpack.BankListHandler)
				//添加银行卡
				m.Get("/yhkadd", redpack.BankAddHandler)

				//代理中心
				m.Get("/agent", redpack.AgentHandler)


			})

			//五人对战
			m.Group("/index", func() {
				//五人对战->房间列表
				m.Get("/getRedPacketList", redpack.GetRedPacketListHandler)

				m.Any("/getGaameRedPacketList", redpack.GetGaameRedPacketListHandler)   //todo 需要重构
				//发包记录
				m.Get("/getGaameRedPacketjlSendList",redpack.GetGaameRedPacketjlSendListHandler)
				//发红包(五人对战、牛牛、二八杠)
				m.Get("/add_red_packet", redpack.SendWurenRedPacketHandler)
				//红包信息
				m.Get("/get_red_packet_info",redpack.GetRedPacketInfoHandler)
				//加入红包列表
				m.Get("/join_red_packet", redpack.JoinWurenRedPacketHandler)
			})

			//牛牛
			m.Group("/niuniu", func() {
				//牛牛首页
				m.Get("/", redpack.NiuniuHandler)

				//发红包页面
				m.Get("/add", redpack.NiuniuAddHandler)

			})

			//二八杠
			m.Group("/gang", func() {
				//二八杠首页
				m.Get("/", redpack.GangHandler)

				//二八杠 发红包 页面
				m.Get("/add", redpack.GangAddHandler)
			})

			//扫雷接龙
			m.Group("/jlindex", func() {
				//红包列表
				m.Post("/getRedPacketList", redpack.SaoleiPackListHandler)
				//领取记录列表
				m.Get("/getRedPacketlqList", redpack.SaoleiPackLqListHandler)

				//发红包
				m.Get("/add_red_packet", redpack.SendZhadanRedPacketHandler)

				//开红包（按钮）
				m.Get("/join_open_red_packet", redpack.SaoleiJLOpenRedButtonAjaxHandler)

				//开红包记录
				m.Get("/open_red_packet", redpack.SaoleiRedOpenRecordAjaxHandler)
				//开包详情
				m.Get("/open_red_packets",redpack.OpenPacketDetailsHandler)
			})

			//登录
			m.Group("/login", func() {
				//获取当前用户信息
				m.Get("/getMemberInfo", redpack.GetMemberInfo)
			})

			//五人对战
			m.Group("/wurenduizhan", func() {
				m.Get("/", redpack.WurenDZHandler)
				m.Get("/add", redpack.WurenFahongbaoHandler)
				m.Get("/kai_ok", redpack.WurenKaiOkHandler)
			})

			//扫雷接龙
			m.Group("/saoleijl", func() {
				m.Get("/", redpack.SaoleiJLHandler)
				m.Get("/add", redpack.SaoleiJLAddHandler)

				m.Get("/kai_ok", redpack.SaoleiRedOpenRecordHandler)
			})

			m.Group("/pay", func() {
				//充值金币
				//m.Get("/wxPay",pay.GoldRechargeHandler)

				m.Get("/add_bank_log", weixinModel.WithdrawalsHandler)
				m.Get("/weixin")

			})



		}, admin.UserNeedLogin)

			//登陆页面
			m.Get("/home/login",admin.UserLoginHandler)
			//登陆提交地址
			m.Post("/home/login",binding.Bind(admin.LoginForm{}), admin.UserLoginPostHandler)
			//注册页面
			m.Get("/home/sign_up",admin.SignUpHandler)
			//注册提交地址
			m.Post("/home/sign_up",binding.Bind(admin.SiginUpTable{}),admin.SignUpTableValuesHandler)


		//代理申请
		m.Get("/weixin/agent/apply", weixin.NeedWxLogin, weixin.ApplyHandler)
		m.Post("/weixin/agent/apply", weixin.NeedWxLogin, weixin.NeedIsGamer, admin.NeedCaptcha, binding.BindIgnErr(weixin.ApplyForm{}), weixin.ApplyPostHandler)
		//微信充值回调
		m.Any("/mp/pay/callback", weixinModel.WxNotifyHandler)
		//微信领红包
		m.Get("/weixin/get_redpack", weixin.NeedWxLogin, weixin.NeedIsGamer, weixin.GetRedPackHandler)

		//代理申请
		m.Post("/agentpro/apply", binding.BindIgnErr(agentProModel.AgentProRecordRow{}), agentPro.ApplyPostHandler)

		//代理密码登陆
		m.Get("/weixin/agent/login", weixin.LoginPostHandler)

		//旺实富支付页面接口
		m.Get("/weixin/paywap/paymethod", weixin.PayWapPaymethodHandler)
		m.Post("/weixin/paywap/pay", weixin.PayWapPayHandler)
		m.Get("/weixin/paywap/return_page", weixin.PayWapReturnPageHandler)
		m.Post("/weixin/paywap/notify", weixin.PayWapNotifyHandler)

		//首页
		m.Get("/", func(ctx *modules.Context) {
			//ctx.Success("即将跳转至后台！", "/admin", 3)
			//ctx.Redirect("/admin", 302)
			ctx.Redirect("/home", 302)
		})

		//websocket处理
		m.Get("/ws/redpack/room/:id", func(ctx *modules.Context) {
			roomId := ctx.ParamsInt(":id")
			user := ctx.IsLogin()

			keys := map[string]interface{}{
				redModel.TagRoomId: int32(roomId),
				redModel.TagUserId: uint32(0),
			}

			if user != nil {
				keys[redModel.TagUserId] = user.Id
			}

			redModel.WsServer.HandleRequestWithKeys(ctx.Resp, ctx.Req.Request, keys)
		})
	}, weixin.RootMidware)
}
