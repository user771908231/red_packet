//判断是否为空，默认值
if (!js_path) {
	js_path = "";
}

var _orderParams = {
	auth : null, // json: {type:'pwd', data: ''}
	payInfo : null, // json string
	onerr : function(tag, tt, data, errMethod) {
		if (errMethod)
			errMethod(data.code, data.msg, tag, tt);
	},
	onok : function(tag, tt, data) {

	}
}

// 确认支付
function confirmPay(tag, tt, orderParams) {
	var obj = tag;
	var key = obj.attr("data-key");
	var payMethod = findMethod('_pay', key);
	if (!payMethod) {
		alert('js异常, 无相应支付下单方式');
	}
	opts = $.extend({}, _orderParams, orderParams || {});
	payMethod(obj, tt, opts);
}

// 默认支付下单
function _pay(obj, tt, orderParams) {
	var ttStr = tt;
	var url = js_Obj.pay_path;
	var params = {
		tt : ttStr,
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx")
	};
	if (orderParams && orderParams.payInfo) {
		if (typeof orderParams.payInfo === 'string') {
			params['payInfo'] = orderParams.payInfo;
		} else {
			params['payInfo'] = JSON.stringify(orderParams.payInfo);
		}
	}
	if (orderParams && orderParams.auth) {
		params['auth.' + orderParams.auth.type] = orderParams.auth.data;
	}

	$.ajax4lock({
		url : url,
		cache : false,// 不保存缓存
		data : params,
		type : "get",
		success : function(data) {
			if (data.code != 0) {
				// 异常处理
				var m = findMethod('_onerr', data.code);
				if (orderParams && orderParams.onerr && typeof orderParams.onerr === 'function') {
					orderParams.onerr(obj, ttStr, data, m);
				} else {
					if (m) {
						m(data.code, data.msg, obj, ttStr);
					}
				}
				return;
			}
			var payData = data.data;
			var inv = payData.invoke;
			var paytt = payData.tt;
			var ot = payData.ot;
			var payParm = payData.payParam;
			var channelKey = payData.channelKey;
			var mInvoke = findMethod("_invoke_" + inv, channelKey);
			if (!mInvoke) {
				alert('js异常, 无相应支付调起方式, invoke=' + inv + ", channelKey=" + channelKey);
				return;
			}
			try {
				// false表示为支付，不是充值，true表示为充值 isRechr
				mInvoke(params, ot, channelKey, inv, payParm, false, obj);
				if (orderParams && orderParams.onok && typeof orderParams.onok === 'function') {
					orderParams.onok(obj, paytt, data);
				}
			} catch (e) {
				ipay.utils.log(e);
			}
		}
	});
}

// get方式调起
function _doget(params) {
	window.location.href = params.url;
}

// post方式调起
function _dopost(params) {
	var form = $('<form></form>');
	form.hide();
	// 设置表单属性
	form.attr('action', params.url);
	form.attr('method', 'post');
	form.attr('target', '_self');
	// 创建原始
	for (p in params.params) {
		var t_hidden = $('<input type="hidden"/>');
		t_hidden.attr('name', p);
		t_hidden.attr('value', params.params[p]);
		form.append(t_hidden);
	}
	form.appendTo("body")
	form.submit();
	// form用完之后应该清除。
	form.remove();
	return;
}

function findMethod(baseMethodName, key) {
	if (key) {
		try {
			var m = eval(baseMethodName + "_" + key);
			if (m && (typeof m) == "function") {
				return m;
			}
		} catch (e) {
		}
	}
	try {
		var m = eval(baseMethodName);
		if (m && (typeof m) == "function") {
			return m;
		}
	} catch (e) {
	}
}

// 错误处理定义
function _onerr(code, errmsg, selectObj) {
	var msg = errmsg;
	if (!msg) {
		msg = '错误码：' + code;
	}
	error_msg(msg);
}

// token 过期后，跳转
function _onerr_6201(code, errmsg, selectObj, tt) {
	var url = $('body').attr('data-cpurl');
	if (url) {
		$.webToast({
			msg : errmsg + "<br>5秒后将自动返回到商户",
			time : 5000,
			onCallback : function() {
				window.location.href = url;
			}
		});
	} else {
		_onerr(code, errmsg, selectObj);
	}
}

// 订单支付超时，请返回商户重新下单(5515)
function _onerr_5515(code, errmsg, selectObj, tt) {
	_onerr_6201(code, errmsg, selectObj, tt);
}

// 特殊错误提示逻辑xxx为特殊错误码
// PAY_USER_CHANGED
function _onerr_6003(code, errmsg, selectObj) {
	__onerr_doreload(code, errmsg, selectObj);
}
// PAY_ACCOUNT_LOCKED
function _onerr_6004(code, errmsg, selectObj) {
	__onerr_doreload(code, errmsg, selectObj);
}
// PAY_PASSWD_NEEDED
function _onerr_6006(code, errmsg, selectObj) {
	__onerr_doreload(code, errmsg, selectObj);
}
// 商户订单已经被完成支付
function _onerr_5513(code, errmsg, selectObj, tt) {
	var msg = errmsg;
	if (!msg) {
		msg = '错误码：' + code;
	}
	msg += '。稍后将跳转支付结果页...'
	var url = js_Obj.qrp_path;
	url = ipay.utils.toGetURL(url, {
		"tt" : tt,
	});
	$.webToast({
		msg : msg,
		onCallback : function() {
			window.location.href = url;
		}
	});
}

function __onerr_doreload(code, errmsg, selectObj) {
	var msg = errmsg;
	if (!msg) {
		msg = '错误码：' + code;
	}
	msg += '。稍后将自动刷新页面...'
	$.webToast({
		msg : msg,
		onCallback : function() {
			ipay.utils.reload();
		}
	});
}

function error_msg(msg) {
	$.abAlert({
		okKnow : '确认',
		showTitle : false,
		showOk : false, // 是否显示确定按钮
		showCancel : false,// 是否显示取消按钮
		msg : '<p>' + msg + '</p>'
	}, 'pop_confirm');
}

// 处理fastpay
function _pay_fastpay(obj, tt, orderParams) {
	var url = js_Obj.portal_path;
	var objs = {
		tt : tt,
		// payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx"),
	// payKey : obj.attr("data-key"),
	// payPre : obj.attr("data-pre")
	}
	window.location.href = ipay.utils.toGetURL(url, objs);
}

// 处理fastpay充值
function _rechrPay_fastpay(obj, tt, orderParams) {
	var url = js_Obj.portal_rechr_path;
	var objs = {
		tt : tt,
		// payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx"),
		// payKey : obj.attr("data-key"),
		// payPre : obj.attr("data-pre"),
		money : orderParams.money
	}
	window.location.href = ipay.utils.toGetURL(url, objs);
}

// 特殊支付下单的定义。
// 点卡调起,自定义函数格式为："_pay_" + 付款类型的key。
function _pay_gc(obj, tt, orderParams) {
	var url = js_Obj.pgc_path;
	var objs = {
		'tt' : tt,
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx")
	}
	window.location.href = ipay.utils.toGetURL(url, objs);
}

function _pay_tc(obj, tt, orderParams) {
	var url = js_Obj.ptc_path;
	var objs = {
		'tt' : tt,
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx")
	}
	window.location.href = ipay.utils.toGetURL(url, objs);
}

// 特殊支付调起的定义
function _invoke_NO(params, ot, channelKey, inv, payParam, isRechr) {
	// 支付下单后无需特殊处理，直接查询支付结果
	var url = js_Obj.qrp_path;
	if (isRechr) {
		url = js_Obj.rechrQrp_path;
	}
	url = ipay.utils.toGetURL(url, {
		"tt" : params.tt,
		"ot" : ot
	});
	var qr_url = js_Obj.qr_path;
	if (isRechr) {
		qr_url = js_Obj.rechrqr_path;
	}
	// 鎖屏
	var options = {
		tip : "正在确认支付结果..."
	}
	$.mask(options);
	setTimeout(function() {
		__qr_result(qr_url, url, {
			s : new Date().getTime(),// 开始查询时间点
			i : 3000,// 查询时间间隔
			timeout : 15000,// 查询时长
			onstart : function() {
				$.mask(options);
			},
			onfinish : function() {
				$.unmask(options);
			}
		}, params.tt, ot);
	}, 1000);
}

function _invoke_UDS(params, ot, channelKey, inv, payParam, isRechr) {
}

function _invoke_JS(params, ot, channelKey, inv, payParam, isRechr) {
	alert('系统异常：channel:' + channelKey + ", isRechr=" + isRechr);
}

function _invoke_WEB(params, ot, channelKey, inv, payParam, isRechr, obj) {
	if (!payParam) {
		alert("支付方式暂不可用（商户通道配置异常）。请联系商户或客服处理。");
		return;
	}
	var param = JSON.parse(payParam);
	var url = param.url;
	if (!ipay.utils.url.isHttp(url)) {
		abComfirm({
			'url' : url,
			'tt' : params.tt,
			'ot' : ot,
			'isRechr' : isRechr
		}, getPayTypeInfo(obj));
		if ('qh360' == ipay.utils.browser.getUa() && ipay.utils.browser.isIOS()) {
			_ifr(url);
		} else {
			top.location.href = url;
		}
	} else if (param.type && param.type.toUpperCase() == "GET") {
		_doget(param);
	} else {
		_dopost(param);
	}
}

function _invoke_WEB_wxwap(tt_params, ot, channelKey, inv, payParam, isRechr, obj) {
	if (!payParam) {
		alert("支付方式暂不可用（商户通道配置异常）。请联系商户或客服处理。");
		return;
	}
	var params = JSON.parse(payParam);
	var url = params.url;
	var comfirmParams = {
		'url' : url,
		'tt' : tt_params.tt,
		'ot' : ot,
		'isRechr' : isRechr
	};
	var app = getPayTypeInfo(obj);
	if (!ipay.utils.url.isHttp(url)) {
		abComfirm(comfirmParams, app);
		if ('qh360' == ipay.utils.browser.getUa() && ipay.utils.browser.isIOS()) {
			_ifr(url);
		} else {
			top.location.href = url;
		}
	} else {
		if (!ipay.utils.isInIframe() && /^https:\/\/wx.tenpay.com/i.test(url)) {
			var getParams = {
				payType : tt_params.payType,
				r : url,
				ot : ot,
				burl : js_Obj['js_path']
			};
			if (tt_params.tt) {
				getParams['tt'] = tt_params.tt;
			}
			getParams['openSafari'] = ('safari' == ipay.utils.browser.getUa());
			window.location.href = ipay.utils.toGetURL(js_Obj['path_proxy_wxh5'], getParams);
		} else {
			abComfirm(comfirmParams, app);
			_ifr(url);
		}
	}
}

// 处理qq钱包的情况
function _invoke_WEB_qqwap(tt_params, ot, channelKey, inv, payParam, isRechr, obj) {
	if (!payParam) {
		alert("支付方式暂不可用（商户通道配置异常）。请联系商户或客服处理。");
		return;
	}
	var params = JSON.parse(payParam);
	var url = params.url;
	abComfirm({
		'url' : url,
		'tt' : tt_params.tt,
		'ot' : ot,
		'isRechr' : isRechr
	}, getPayTypeInfo(obj));
	if (ipay.utils.url.isHttp(url)) {
		if (!ipay.utils.isInIframe() && /^https:\/\/myun.tenpay.com/i.test(url)) {
			var getParams = {
				payType : tt_params.payType,// 支付方式类型
				r : url,
				ot : ot,
				burl : js_Obj['js_path']
			};
			if (tt_params.tt) {
				getParams['tt'] = tt_params.tt;
			}
			getParams['openSafari'] = ('safari' == ipay.utils.browser.getUa());
			window.location.href = ipay.utils.toGetURL(js_Obj['path_proxy_appwap'], getParams);
		} else {
			_ifr(url);
		}
	} else {
		top.location.href = url;
	}
}

function _invoke_DAT(params, ot, channelKey, inv, payParam, isRechr) {
	if (!payParam) {
		alert("支付方式暂不可用（商户通道配置异常）。请联系商户或客服处理。");
		return;
	}
	var param = JSON.parse(payParam);
	var data = {
		tt : params.tt,
		ot : ot,
		r : param.url,
		payex : params.payEx,
		payType : params.payType
	}
	var base_url = js_Obj.qrc_path;
	if (isRechr) { // 充值时
		base_url = js_Obj.qrcRechr_path;
	}
	var url = ipay.utils.toGetURL(base_url, data);
	window.location.href = url;
}

function _invoke_WEB_aliwap(tt_params, ot, channelKey, inv, payParam, isRechr, obj) {
	if (!payParam) {
		alert("支付方式暂不可用（商户通道配置异常）。请联系商户或客服处理。");
		return;
	}
	var params = JSON.parse(payParam);
	var url = params.url;
	// 支付宝wap,微信浏览器中走代理
	if (ipay.utils.isWXBrowser()) {
		if (params['type'] == "POST") {
			url = ipay.utils.toGetURL(url, params['params']);
		}
		var getParams = {
			r : url,
			ot : ot,
			burl : js_Obj['js_path']
		};
		if (tt_params.tt) {
			getParams['tt'] = tt_params.tt;
		}
		window.location.href = ipay.utils.toGetURL(js_Obj['path_proxy_wx_ali'], getParams);
	} else {
		if (!ipay.utils.url.isHttp(url)) {
			abComfirm({
				'url' : url,
				'tt' : tt_params.tt,
				'ot' : ot,
				'isRechr' : isRechr
			}, getPayTypeInfo(obj));
			if ('safari' == ipay.utils.browser.getUa() || ipay.utils.browser.isHappyApp()) {
				top.location.href = url;
			} else {
				_ifr(url);
			}
		} else {
			_invoke_WEB(tt_params, ot, channelKey, inv, payParam, isRechr, obj);
		}

	}
}

function _invoke_WEB_aliwap_qr(tt_params, ot, channelKey, inv, payParam, isRechr, obj) {
	_invoke_WEB_aliwap(tt_params, ot, channelKey, inv, payParam, isRechr, obj);
}

function _ifr(url) {
	$("body").append("<iframe id='ifr_create' src='" + url + "' style='display: none;'></iframe>");
	setTimeout(function() {
		$("#ifr_create").remove();
	}, 4000);
}

function abComfirm(params, conf) {
	var confirm = function() {
		$.unmask();
		$.abAlert({
			showCloseBtn : true,
			showOk : false, // 点击确定按钮
			showCancel : false,// 点击取消按钮
			title : '<p style="display: none;border-bottom: none;"></p>',
			msg : '<p>请在' + conf.appName + '内完成支付</p>',
			butTxt : '付款已完成',
			showFooter : true,
			footerTxt : '无法打开' + conf.payTypeName + '？<span class="callWX"  style="color: #0c6cc5;text-decoration:underline;">点此尝试</span>',
			onFooter : function() {
				window.open(params.url);
			},
			onButton : function() {
				// 立刻领取
				var $QueryUrl = js_Obj.qrp_path;
				if (params.isRechr) {
					$QueryUrl = js_Obj.rechrQrp_path;
				}
				$QueryUrl = ipay.utils.toGetURL($QueryUrl, {
					"tt" : params.tt,
					"ot" : params.ot
				});
				window.location.href = $QueryUrl;
			}
		}, 'pop_alertMsg');
	}
	$.mask({
		tip : '打开' + conf.appName + '中...'
	});
	setTimeout(confirm, 1000);
}

// 确认充值
function confirmRechrPay(tag, tt, orderParams) {
	var obj = tag;
	var key = obj.attr("data-key");
	var rechrMethod = findMethod('_rechrPay', key);
	if (!rechrMethod) {
		alert('js异常, 无相应支付下单方式');
	}
	opts = $.extend({}, _orderParams, orderParams || {});
	rechrMethod(obj, tt, opts);
}

// 默认充值下单
function _rechrPay(obj, tt, orderParams) {
	var url = js_Obj.rechrPay_path;
	var params = {
		tt : tt,
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx"),
		money : obj.attr("data-money")
	};
	if (orderParams && orderParams.payInfo) {
		if (typeof orderParams.payInfo === 'string') {
			params['payInfo'] = orderParams.payInfo;
		} else {
			params['payInfo'] = JSON.stringify(orderParams.payInfo);
		}
	}
	if (orderParams && orderParams.auth) {
		params['auth.' + orderParams.auth.type] = orderParams.auth.data;
	}
	if (orderParams && orderParams.money) {
		params['money'] = orderParams.money;
	}
	rechr_common(url, tt, params, obj, orderParams);

}

function rechr_common(url, tt, params, obj, orderParams) {
	var tt_params = params;
	var rechr_orderParams = orderParams;
	var rechr_obj = obj;
	$.ajax4lock({
		url : url,
		dataType : "json",
		cache : false,// 不保存缓存
		data : params,
		type : "get",
		success : function(data) {
			if (data.code != 0) {
				// 异常处理
				var m = findMethod('_onerr', data.code);
				if (rechr_orderParams && rechr_orderParams.onerr && typeof rechr_orderParams.onerr === 'function') {
					rechr_orderParams.onerr(rechr_obj, tt_params.tt, data, m);
				} else {
					if (m) {
						m(data.code, data.msg, rechr_obj, tt_params.tt);
					}
				}
				return;
			}
			var payData = data.data;
			var inv = payData.invoke;
			var ot = payData.ot;
			var payParm = payData.payParam;
			var channelKey = payData.channelKey;
			var mInvoke = findMethod("_invoke_" + inv, channelKey);
			if (!mInvoke) {
				alert('js异常, 无相应支付调起方式, invoke=' + inv + ", channelKey=" + channelKey);
				return;
			}
			try {
				mInvoke(tt_params, ot, channelKey, inv, payParm, true);
				if (rechr_orderParams && rechr_orderParams.onok && typeof rechr_orderParams.onok === 'function') {
					rechr_orderParams.onok(obj, tt_params.tt, data);
				}
			} catch (e) {
				console.log(e);
			}
		}
	});
}

function _rechrPay_tc(obj, tt) {
	var url = js_Obj.rtc_path;
	var objs = {
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx"),
		tt : tt
	}
	window.location.href = ipay.utils.toGetURL(url, objs);
}

// 点卡充值定义
function _rechrPay_gc(obj, tt) {
	var url = js_Obj.rgc_path;
	var objs = {
		payType : obj.attr("data-payType"),
		payEx : obj.attr("data-payEx"),
		tt : tt
	}
	window.location.href = ipay.utils.toGetURL(url, objs);

}

/**
 * 支付限额判断.使用方法： 限额表达式 exp. 应支付金额pay（单位元） paylimit(exp).isValid(pay)，返回值true为可支付。
 */
paylimit = function(exp) {
	var fs = {
		ENUM : function(t) {
			this._limits = [];
			this.isValid = function(v) {
				if (v) {
					for (i in this._limits) {
						if (this._limits[i] == v) {
							return true;
						}
					}
				}
				return false;
			}
			this.build = function(t) {
				if (t == undefined) {
					return;
				}
				var s = t.split(",");
				for (i in s) {
					this._limits.push(s[i].trim());
				}
			}
			this.toString = function() {
				return 'ENUM:' + this._limits.join(',');
			}
			this.build(t);
			this.max = function() {
				if (this._limits.length > 0) {
					var a = this._limits[0];
					for (i = 1; i < this._limits.length; i++) {
						a = Math.max(a, this._limits[i]);
					}
					return a;
				}
			};
			this.min = function() {
				if (this._limits.length > 0) {
					var a = this._limits[0];
					for (i = 1; i < this._limits.length; i++) {
						a = Math.min(a, this._limits[i]);
					}
					return a;
				}
			};
			this.getLimits = function() {
				return [].concat(this._limits);
			}
		},

		RANGE : function(t) {
			this._min = null;
			this._max = null;

			this.isValid = function(v) {
				if (v) {
					if (this._min == null || this._min <= v) {
						if (this._max == null || this._max >= v) {
							return true;
						}
					}
				}
				return false;
			}
			this.build = function(t) {
				if (t == undefined) {
					return;
				}
				if (t.length == 0) {
					return;
				}
				var s = t.split("-");
				if (s.length > 0) {
					s[0] = s[0].trim();
					if (s[0].length > 0) {
						this._min = new Number(s[0]);
					}
					if (s.length > 1) {
						s[1] = s[1].trim();
						if (s[1].length > 0) {
							this._max = new Number(s[1]);
						}
					}
				}
			}
			this.toString = function() {
				return 'RANGE:' + this._min + '-' + this._max;
			}
			this.build(t);
			this.max = function() {
				if (this._max != null) {
					return this._max;
				}
			};
			this.min = function() {
				if (this._min != null) {
					return this._min;
				}
			};
		},

		IRANGE : function(t) {
			this._min = null;
			this._max = null;

			this.isValid = function(v) {
				if (v) {
					if (Math.floor(v) == v) {
						if (this._min == null || this._min <= v) {
							if (this._max == null || this._max >= v) {
								return true;
							}
						}
					}
				}
				return false;
			}
			this.build = function(t) {
				if (t == undefined) {
					return;
				}
				if (t.length == 0) {
					return;
				}
				var s = t.split("-");
				if (s.length > 0) {
					s[0] = s[0].trim();
					if (s[0].length > 0) {
						this._min = new Number(s[0]);
					}
					if (s.length > 1) {
						s[1] = s[1].trim();
						if (s[1].length > 0) {
							this._max = new Number(s[1]);
						}
					}
				}
			}
			this.toString = function() {
				return 'IRANGE:' + this._min + '-' + this._max;
			}
			this.build(t);
			this.max = function() {
				if (this._max != null) {
					return this._max;
				}
			};
			this.min = function() {
				if (this._min != null) {
					return this._min;
				}
			};
		},
		MULTI : function(t) {
			this._limits = [];
			this.build = function(t) {
				if (t == undefined) {
					return;
				}
				if (t.length == 0) {
					return;
				}
				var s = t.split(";");
				for (i in s) {
					var exp = s[i].trim();
					if (exp.lenth == 0) {
						continue;
					}
					var i = exp.indexOf(':');
					if (i >= 0) {
						var v = exp.substring(i + 1)
						var k = exp.substring(0, i);
						if (fs[k]) {
							this._limits.push(new fs[k](v));
						} else {
							console.log("No paylimit for:[" + k + "], MULTI:" + t);
						}
					}
				}
			}
			this.isValid = function(v) {
				if (v) {
					if (this._limits.length == 0) {
						return true;
					}
					for (i in this._limits) {
						if (this._limits[i].isValid(v)) {
							return true;
						}
					}
				}
				return false;
			}
			this.build(t);
			this.max = function() {
				if (this._limits.length > 0) {
					var a = this._limits[0].max();
					if (a == undefined || a == null) {
						return a;
					}
					for (i = 1; i < this._limits.length; i++) {
						var t = this._limits[i].max();
						if (t == undefined || t == null) {
							return t;
						}
						a = Math.max(a, t);
					}
					return a;
				}
			};
			this.min = function() {
				if (this._limits.length > 0) {
					var a = this._limits[0].min();
					if (a == undefined || a == null) {
						return a;
					}
					for (i = 1; i < this._limits.length; i++) {
						var t = this._limits[i].min();
						if (t == undefined || t == null) {
							return t;
						}
						a = Math.min(a, this._limits[i].min());
					}
					return a;
				}
			};
		},
		NO : function(t) {
			this.isValid = function(v) {
				if (v) {
					return true;
				}
				return false;
			}
			this.build = function(t) {
			}
			this.toString = function() {
				return 'NOLIMIT';
			}
			this.build(t);
			this.max = function() {
			};
			this.min = function() {
			};
		}
	};

	if (exp) {
		var i = exp.indexOf(':');
		if (i >= 0) {
			var v = exp.substring(i + 1)
			var k = exp.substring(0, i);
			if (fs[k]) {
				return new fs[k](v);
			} else {
				console.log("No paylimit for:[" + k + "], will use NO LIMIT");
			}
		}
	}
	return new fs['NO']();
}

rangelimit = function(exp) {
	if (exp == undefined) {
		return paylimit();
	}
	return paylimit("RANGE:" + exp);
}

/**
 * 充值手续费.使用方式。 手续费表达式 exp. 充值金额money（单位元） payfee(exp).calc(money)，返回值{pay: xxx,
 * fee: xxx, money: xxx}, pay为应支付金额，fee为手续费。money为到账金额，单位均为元。
 */
payfee = function(exp) {
	var fee = {
		RATE : function(r) {
			this.rate = 0;
			this.calc = function(v) {
				if (this.rate == 0) {
					return {
						pay : v,
						fee : 0,
						money : v
					};
				}
				var pay = Math.ceil(10000 * v / (100 - this.rate)) / 100;
				var r = (pay - v).toFixed(2);
				return {
					pay : pay,
					fee : r,
					money : v
				};
			}
			this.calcByPay = function(pay) {
				if (this.rate == 0) {
					return {
						pay : pay,
						fee : 0,
						money : pay
					};
				}
				var money = Math.floor(pay * (100 - this.rate)) / 100;
				var r = (pay - money).toFixed(2);
				return {
					pay : pay,
					fee : r,
					money : money
				};
			}
			this.build = function(r) {
				if (r && r >= 0 && r < 100) {
					this.rate = r;
				}
			}
			this.toString = function() {
				return 'RATE:' + this.rate + '%';
			}
			this.build(r);
		},
		NO : function() {
			this.calc = function(v) {
				return {
					pay : v,
					fee : 0,
					money : v
				};
			}
			this.calcByPay = function(v) {
				return {
					pay : v,
					fee : 0,
					money : v
				};
			}
		}
	};

	if (exp) {
		var i = exp.indexOf(':');
		if (i >= 0) {
			var v = exp.substring(i + 1)
			var k = exp.substring(0, i);
			if (fee[k]) {
				return new fee[k](v);
			} else {
				console.log("No payfee for:[" + k + "], will use NO FEE");
			}
		}
	}
	return new fee['NO']();
}

errcode = {
	PAY_PASSWD_INVALID : '6005'
}

// 查询支付结果处理
function __qr_result(qrurl, qr_pageurl, condition, tt, ot) {
	$.ajax({
		url : qrurl,
		dataType : "json",
		data : {
			ot : ot,
			tt : tt
		},
		type : "get",
		cache : false,
		success : function(data) {
			if (data.code != 0) {
				__pay_result_PAY_ING(null, qrurl, qr_pageurl, condition, tt, ot);
				return;
			}
			__pay_result(data.data.pay, qrurl, qr_pageurl, condition, tt, ot);
		},
		error : function() {
			__pay_result_PAY_ING(null, qrurl, qr_pageurl, condition, tt, ot);
		}
	});
}

var __pay_result = function(pay, qrurl, qr_pageurl, condition, tt, ot) {
	var m = ipay.utils.findMethod('__pay_result_' + pay.result);
	m(pay, qrurl, qr_pageurl, condition, tt, ot);
}

function _onResultErr(pay) {
	error_msg(pay.err);
}

function __pay_result_WAIT(pay, qrurl, qr_pageurl, condition, tt, ot) {
	__pay_result_PAY_ING(pay, qrurl, qr_pageurl, condition, tt, ot);
}

function __pay_result_PAY_FAIL(pay, qrurl, qr_pageurl, condition, tt, ot) {
	if (condition.onfinish) {
		condition.onfinish();
	}
	_onResultErr(pay);
}

function __pay_result_PAY_CANCEL(pay, qrurl, qr_pageurl, condition, tt, ot) {
	__pay_result_PAY_FAIL(pay, qrurl, qr_pageurl, condition, tt, ot);
}

function __pay_result_PLAT_FAIL(pay, qrurl, qr_pageurl, condition, tt, ot) {
	if (condition.onfinish) {
		condition.onfinish();
	}
	window.location.href = qr_pageurl;
}

function __pay_result_PAY_OK(pay, qrurl, qr_pageurl, condition, tt, ot) {
	if (condition.onfinish) {
		condition.onfinish();
	}
	window.location.href = qr_pageurl;
}

function __pay_result_PAY_ING(pay, qrurl, qr_pageurl, condition, tt, ot) {
	if (new Date().getTime() > (condition.s + condition.timeout)) {// 达到查询时间
		if (condition.onfinish) {
			condition.onfinish();
		}
		$.abAlert({
			okKnow : '',
			okTxt : "继续",
			cancTxt : "放弃",
			showTitle : false,
			showOk : true, // 是否显示确定按钮
			showCancel : true,// 是否显示取消按钮
			msg : '<p>未查询到支付结果，请确认是否继续查询</p>',
			onAccept : function() {
				if (condition.onstart) {
					condition.onstart();
				}
				condition.s = new Date().getTime();
				__pay_result_PAY_ING(pay, qrurl, qr_pageurl, condition, tt, ot);
			},
			onCancel : function() {
				if (condition.onfinish) {
					condition.onfinish();
				}
			}
		}, 'pop_confirm');
	} else {
		t = setTimeout(function() {
			__qr_result(qrurl, qr_pageurl, condition, tt, ot);
		}, condition.i);
	}
}

function getPayTypeInfo(obj) {
	if (typeof obj == "undefined") {
		return;
	}
	var payTypeName = obj.attr("data-name");
	var appName = obj.attr("data-app");
	appName = appName ? appName : payTypeName;
	var app = {
		'appName' : appName,
		'payTypeName' : payTypeName
	}
	return app;
}