/**
 * Created by admin on 2016/9/14.
 */
// r用户cookie前缀
var r_prefix = ".IPAY.R.";
function flagRUser(userid, isRUser) {
	if (userid) {
		if (isRUser) {
			ipay.utils.cookie.set(r_prefix + userid, '1', '', '/');
		} else {
			ipay.utils.cookie.set(r_prefix + userid, '0', '', '/');
		}
	}
}

function isRUser(userid) {
	if (userid) {
		return ipay.utils.cookie.get(r_prefix + userid) == '1';
	}
}

// 可以有小数点,保留两位
function clearPointNoNum(obj) {
	obj.value = obj.value.replace(/^[^\d]*/, "");
	obj.value = obj.value.replace(/^[0]+/, "0");
	obj.value = obj.value.replace(/^0(\d+.*)/, "$1");
	obj.value = obj.value.replace(/[^\d.]/, "");
	obj.value = obj.value.replace(/(\d+)(\.\d{0,2}).*/, "$1$2");
}

// 此方法只可以输入数字，且首位可以为0.
function clearNoNum_zero(obj) {
	// 先把非数字的都替换掉，除了数字
	obj.value = obj.value.replace(/[^0-9]/g, "");
}

// 此方法只可以输入数字，且首位不能为0.
function clearNoNum(obj) {
	// 先把非数字的都替换掉，除了数字
	obj.value = obj.value.replace(/^0/g, "");
	obj.value = obj.value.replace(/[^0-9]/g, "");
}

$(function() {
	// 解决移动端不支持active
	document.body.addEventListener('touchstart', function() {
	}, false);
	// fastclick
	window.addEventListener('load', function() {
		if (typeof (FastClick) == "undefined")
			return;
		FastClick.attach(document.body);
	}, false);
	// 替换虚拟币
	$(".exchangeBalance-vc").each(function() {
		var formater = new Formater($(this).html());
		var money = $(this).attr("data-balance");
		$(this).html(formater.exec({
			balance_vc : fenToVC(money),
			balance_rmb : fenToYuan(money)
		}));
	});
});
(function(doc, win) {
	var docEl = doc.documentElement, resizeEvt = 'orientationchange' in window ? 'orientationchange' : 'resize', recalc = function() {
		var clientWidth = docEl.clientWidth;
		if (!clientWidth)
			return;
		docEl.style.fontSize = 100 * (clientWidth / 750) + 'px';
	};
	if (!doc.addEventListener)
		return;
	win.addEventListener(resizeEvt, recalc, false);
	doc.addEventListener('DOMContentLoaded', recalc, false);
})(document, window);

var passwdPubKey = null;
function getPasswdPubKey() {
	if (!passwdPubKey) {
		setMaxDigits(129);
		passwdPubKey = new RSAKeyPair("10001", "",
				"00d03bbb783cfbaef7ce20279e4f9b5d41877fb8882f92fd0a992665e10af44656af48616cc3e0c0b02a11e7fd0771f62db3e58a8f36a7865cb1ec5d48ab919ed80de648b4e4adf341e3a8e91872ea61527a57e6442439aba3e4edfca014d53b83cee548da9020ee19d7a6f43a8c501050343d1a0eb7efbaa5ab5da8f4e5bd2081",
				1024);
	}
	return passwdPubKey;
}
function encPwd(pwd) {
	getPasswdPubKey();
	var encv = encryptedString(passwdPubKey, pwd, RSAAPP.PKCS1Padding, RSAAPP.RawEncoding);
	return "!!0001" + window.btoa(encv);
}

$(function() {
	$.mask = function(options) {
		this.id = $.nexTagID('mask');
		var maskOptions = {
			tip : ''
		};
		var settings = $.extend({}, maskOptions, options);
		$.maskUI({
			id : this.id,
			message : '<span style="width:33px; height:33px; background:url(' + ipay['cfg']['icon_loading']
					+ ') no-repeat; background-size:100%; position:relative;  display:block;margin: 20px auto 10px auto;">&nbsp;</span><span style="text-align: center; color:#eee; font-size:14px;">' + settings.tip + '</span>',
			overlayCSS : {
				backgroundColor : 'rgba(0,0,0,0)'
			},
			css : {
				border : 'none',
				'text-align' : 'center',
				width : '95px',
				'line-height' : '16px',
				position : 'fixed',
				left : '50%',
				'margin-left' : '-48px',
				right : '0',
				top : '50%',
				'margin-top' : '-16px',
				background : '#000',
				opacity : '0.4',
				'border-radius' : '8px',
				'padding-bottom' : '10px'
			}
		})
	};
	$.unmask = function() {
		$.unmaskUI({
			id : this.id
		});
	};
	$.abAlert = function(options, config) {
		options = options || {};
		var defaultOptions = {
			id : $.nexTagID('cp'),
			title : '',
			msg : '',
			okKnow : '',
			okTxt : '确定',
			cancTxt : '取消',
			textAlign : 'center',
			showTitle : true,
			showKnow : true,
			showCloseBtn : false,
			// 是否显示关闭按钮
			showOk : true,
			// 是否显示确定按钮
			showCancel : true,
			// 是否显示取消按钮
			showFooter : false,
			// 是否显示最底下文字
			butTxt : '',
			footerTxt : '',
			acceptLink : 'javascript:void(0)',
			onKnow : function() {
			},
			onAccept : function() {
			},
			onCancel : function() {
			},
			onClose : function() {
			},
			onButton : function() {
			},
			onFooter : function() {
			}
		};
		options = $.extend({}, defaultOptions, options);
		$.maskUI({
			id : options.id,
			message : '<div data-id="' + options.id + '"  class="pop ' + config + '"><p class="title">' + options.title + '</p><a href="javascript:void(0)" class="close ac_Close">X</a>' + '<div class="body" ' + options.msg
					+ '</div><div class="op"><a href="javascript:void(0);" class="ac_popCancel">' + options.cancTxt + '</a><a href="' + options.acceptLink + '" class="ac_popAccept">' + options.okTxt + '</a><a href="' + options.acceptLink + '" class="ac_popKnow">'
					+ options.okKnow + '</a></div><a href="javascript:void(0);" class="ac_popBut">' + options.butTxt + '</a><p class="ac_popFooter"> ' + options.footerTxt + '</p></div>',
			overlayCSS : {
				backgroundColor : 'rgba(0,0,0,0.4)'
			},
			css : {
				width : 'auto',
				position : 'fixed',
				top : '40%',
				left : '50px',
				right : '50px',
				backgroundColor : '#fff',
				'border-radius' : '8px',
				border : 'none'
			}
		});
		var $pop_alert = $('.pop');
		if (options.showTitle) {
			$pop_alert.find('.title').show();
		} else {
			$pop_alert.find('.title').hide();
		}
		if (options.showCloseBtn) {
			$pop_alert.find('.ac_Close').show();
		} else {
			$pop_alert.find('.ac_Close').hide();
		}
		if (options.showKnow) {
			$pop_alert.find('.ac_popKnow').show();
		} else {
			$pop_alert.find('.ac_popKnow').hide();
		}
		if (options.showOk) {
			$pop_alert.find('.ac_popAccept').show();
		} else {
			$pop_alert.find('.ac_popAccept').hide();
		}
		if (options.showCancel) {
			$pop_alert.find('.ac_popCancel').show();
		} else {
			$pop_alert.find('.ac_popCancel').hide();
		}
		if (options.showFooter) {
			$pop_alert.find('.ac_popFooter').show();
		} else {
			$pop_alert.find('.ac_popFooter').hide();
		}
		$('.pop .ac_popCancel').click(function() {
			var dataId = $(this).parents('.pop').attr('data-id');
			$.unmaskUI({
				id : dataId
			});
			options.onCancel();
		});
		$('.pop .ac_popKnow').click(function() {
			var dataId = $(this).parents('.pop').attr('data-id');
			$.unmaskUI({
				id : dataId
			});
			options.onKnow();
		});
		$('.pop .ac_popAccept').click(function() {
			var dataId = $(this).parents('.pop').attr('data-id');
			$.unmaskUI({
				id : dataId
			});
			options.onAccept();
		});
		$('.pop .ac_popBut').click(function() {
			var dataId = $(this).parents('.pop' + options.id + '').attr('data-id');
			$.unmaskUI({
				id : dataId
			});
			options.onButton();
		});
		$('.pop .ac_Close').click(function() {
			var dataId = $(this).parents('.pop').attr('data-id');
			$.unmaskUI({
				id : dataId
			});
			options.onClose();
		});
		$('.pop .callWX').click(function() {
			var dataId = $(this).parents('.pop').attr('data-id');
			options.onFooter();
		});
	};
	$.webToast = function(options) {
		this.id = $.nexTagID('toast');
		var dcfg = {
			msg : "提示信息",
			time : 1500,// 展示时间
			onCallback : function() {
			}
		};
		var settings = $.extend({}, dcfg, options);
		$.maskUI({
			id : this.id,
			message : '<span style="text-align: center; color:#eee; font-size:14px; line-height: 20px;">' + settings.msg + '</span>',
			overlayCSS : {
				backgroundColor : 'rgba(0,0,0,0)'
			},
			css : {
				border : 'none',
				'line-height' : '16px',
				float : 'left',
				left : '33%',
				width : 'auto',
				right : '33%',
				'padding' : '10px',
				top : '50%',
				background : '#000',
				opacity : '0.4',
				'border-radius' : '4px',
				'font-size' : '14px',
				'text-align' : 'center'
			}
		});
		var toast = this;
		this.close = function() {
			$.unmaskUI({
				id : toast.id
			});
			if (settings.onCallback && typeof settings.onCallback == 'function') {
				settings.onCallback();
			}
		};
		window.setTimeout(this.close, settings.time);
	};
});

// 分转爱贝币
function fenToVC(fen) {
	if (!fen) {
		return 0.0;
	}
	return (fen / 10).toFixed(1);
}
// 分转元
function fenToYuan(fen) {
	if (!fen) {
		return 0.00;
	}
	return (fen / 100).toFixed(2);
}

// 元转爱贝币
function yuanToVC(yuan) {
	if (!yuan) {
		return 0.0;
	}
	return (yuan * 10).toFixed(1);
}

var js_Obj = {
	js_path : js_path, // 基础路径
	pub_path : pub_path, // 外网路径
	pay_path : js_path + "pay.json", // 下单支付
	qr_path : js_path + "qr.json", // 查询支付结果
	qrp_path : js_path + "qrp", // 查询页面路径
	pgc_path : js_path + "ptype/gc", // 游戏点卡页面
	ptc_path : js_path + "ptype/tc", // 话费支付页面
	qrc_path : js_path + 'qrc', // 微信扫码
	qrcRechr_path : js_path + 'rechr/qrc',
	rechrHome_path : js_path + "rechr/home",// 充值首页
	person_path : js_path + "u/user",// 个人中心首页
	rgc_path : js_path + "rechr/ptype/gc", // 游戏点卡页面
	rtc_path : js_path + "rechr/ptype/tc", // 话费支付页面
	rechrPay_path : js_path + "rechr/pay.json", // 充值支付
	rechrqr_path : js_path + "rechr/qr.json",// 充值查询结果
	rechrQrp_path : js_path + "rechr/qrp", // 充值查询页面路径
	trans_path : js_path + "u/transRec.json",// 消费记录
	rechr_path : js_path + "u/rechr.json",// 充值记录
	regSms_path : js_path + "u/reg_sms.json",// 发送短信验证码:注册
	fbSms_path : js_path + "u/fb_sms.json",// 发送短信验证码:找回密码
	fbPwd_path : js_path + "u/fb_pwd.json",// 找回密码接口
	isRegisted_path : js_path + "u/isRegisted.json",// 判断手机号是否已注册：用于注册
	findBackCheck_path : js_path + "u/fb_check.json",// 判断手机号是否已注册:用于找回
	verifyRegSms_path : js_path + "u/verifyRegSms.json",// 校验注册短信验证码
	verifyFBSms_path : js_path + "u/verifyFBSms.json",// 校验注册短信验证码
	reg_path : js_path + "u/reg.json",// 注册接口
	login_path : js_path + "u/login.json",// 登录接口
	verifyPayPwd_path : js_path + "u/verifyPayPwd.json",// 校验支付密码接口
	verifyLoginPwd_path : js_path + "u/verifyLoginPwd.json",// 校验登录密码接口
	changePayPwd_path : js_path + "u/changePayPwd.json",// 修改支付密码接口
	changeLoginPwd_path : js_path + "u/changeLoginPwd.json",// 修改登录密码接口
	msg_path : js_path + "notify.json",// 消息栏（跑马灯）
	path_proxy_wx_ali : pub_path + '/h5/proxy/wx-ali',
	path_proxy_wxjsapi : pub_path + '/h5/proxy/wxjsapi',
	path_proxy_wxh5 : pub_path + '/h5/proxy/wxh5',
	path_proxy_appwap : pub_path + '/h5/proxy/appwap',//代理页面，可以处理wx,qq钱包等
	checkCardNo_path : js_path + 'fastpay/checkCardNo.json', // 检查银行卡号信息
	checkCardInfo_path : js_path + 'fastpay/checkCardInfo.json',// 检查银行卡信息
	fastpayResult_path : js_path + 'fastpay/queryResult.json',// 查询结果
	sendVerifyCode_path : js_path + 'fastpay/sendVerifyCode.json',// 发送验证码
	verifyCode_path : js_path + 'fastpay/verifyCode.json',// 验证验证码
	portal_path : js_path + 'fastpay/portal',// 一键支付支付入口
	portal_rechr_path : js_path + 'fastpay/rechr',// 一键支付充值入口
	path_proxy_wxjsapi : pub_path + '/h5/proxy/wxjsapi',
	qrc_order : js_path + 'order.json'
};
