// 去个人中心，包含url中的参数
function goPersonCenter() {
	return ipay.utils.toGetURL(js_Obj.person_path, "tt=" + $("body").attr("data-tt"));
}
var myScroll = null;
var initList = null;
var timeAc = null;
function loaded() {
	if (!myScroll) {
		myScroll = new iScroll('wrapper', {
			scrollbarClass : 'myScrollbar'
		});
		initList = setTimeout(function() {
			adjustPaymentList();
		}, 200);
	}
	// 限制推荐标签字符个数
	$(".rec").each(function() {
		$(this).text($(this).text().substr(0, 8));
	});
}
// 用于在支付列表可用高度变化时，改变高度
function paymentListResize() {
	var hgt = $(".btnPay").offset().top - $("#wrapper").offset().top;
	if ($(window).height() < $(window).width()) {
		hgt = $(window).height() - $(".btnPay").outerHeight();
	}
	$("#wrapper").css("height", hgt);
	if (myScroll) {
		myScroll.refresh();
	}
}
// 用于在支付列表可用高度变化时，调整显示支付方式数目
function adjustPaymentList() {
	paymentListResize();
	var len = $(".paymentList li").length;
	if (len > 3) {
		var liNum = Math.floor($("#wrapper").outerHeight() / $(".paymentList li").outerHeight());
		liNum = liNum == 0 ? 1 : liNum;
		if (len > liNum) {
			liNum = ($("#wrapper").outerHeight() % $(".paymentList li").outerHeight()) > $(".ac_togglePayment").outerHeight() ? liNum : liNum - 1;
			$(".paymentList li").show();
			$(".paymentList li:gt(" + (liNum - 1) + ")").hide();
			$(".morePay").show();
		}
		if (myScroll) {
			myScroll.refresh();
		}
	}
}
function listChgByMorePay(t1, t2) {
	if ($(".morePay").is(":visible")) {
		setTimeout(adjustPaymentList, t1);
	} else {
		setTimeout(paymentListResize, t2);
	}
}

// 初始化
$(function() {
	loaded();
	window.addEventListener('orientationchange', function(event) {
		listChgByMorePay(200, 200);
		listChgByMorePay(500, 500);
		if (timeAc != null) {
			clearTimeout(timeAc);
			$(".acBox li").css("left", 0);
		}
		sysScroll();
	});
	$(".bill").resize(function() {
		listChgByMorePay(0, 0);
	})
	// 横屏时，浏览器自带banner出现时，高度变小，因此在有“更多”按钮时，需要改变显示支付方式个数
	$("#wrapper").resize(function() {
		paymentListResize();
	})

	// 商品信息展示更多文字
	$('.billName').click(function() {
		$(this).toggleClass("current");
		if ($(window).height() > $(window).width()) {
			listChgByMorePay(0, 0);
		}
	});

	// 支付方式拓展按钮
	$(".morePay").click(function() {
		$(".paymentList li").show();
		$(this).hide();
		$(".btnPay + div.footer").hide();
		$(".paymentList div.footer").show();
		$('.btnPay').addClass("down");
		$(".btnPayBg").show();
		paymentListResize();
	});

	var btnPay = $(".btnPay");
	$('.paymentList li:not(.disabledPay)').click(function() {
		$(".paymentList li span").removeClass('current');
		$(this).find("span").addClass('current');
		// 获取价格
		var realFee = $(this).find("font").attr("data-realFee");
		if (realFee) {
			$(".btnPay span").html('¥' + realFee);
		} else {
			$(".btnPay span").html('');
		}
	});

	// 点击返回按钮
	$('.ac_back').click(function() {
		$.abAlert({
			showKnow : false,
			showTitle : false,
			showOk : true, // 是否显示确定按钮
			showCancel : true,// 是否显示取消按钮
			msg : '<p>是否取消支付？</p>',
			onAccept : function() {
				// 点击确定跳转页面
				location.href = $('body').attr('data-cpurl');
			}
		}, 'pop_confirm');
	});

	// 确认支付click事件
	$(".btnPay").click(function() {
		var $CURRENT = $(".secbtn.current").parent();
		var pre = $CURRENT.attr('data-pre');
		// 应付价格
		var realFee = $CURRENT.attr('data-realFee');
		if (!realFee) {
			return;
		}
        var paytype = $CURRENT.attr('data-paytype');
        $(".p9_paymethod").val(paytype);
		document.formpay.submit();
	});

	function goPayObj(paypwd, actions) {
		var t = {
			auth : {
				type : 'pwd',
				data : paypwd
			},
			onerr : function(tag, tt, data, errMethod) {
				if (data.code == errcode['PAY_PASSWD_INVALID']) {
					var passOptions = {
						err : data.msg
					};
					actions.reset(passOptions);
					return;
				}
				actions.close();
				if (errMethod)
					errMethod(data.code, data.msg, tag, tt);
			},
			onok : function(tag, tt, data) {
				actions.close();
			}
		};
		return t;
	}

	// 根据key是否进行展开
	function toExpand(key) {
		var done = null;
		var a = $("li:not(.disabledPay) font." + key);
		var li = a.parent('li:first');
		// 对象存在就选中。
		if (a.size()) {
			done = key;
			li.click();
		} else {
			return done;
		}
		// 选中对象是否为可见。
		var i = li.index();
		if (i > 2) {
			clearTimeout(initList);
			$(".morePay").click();
			myScroll.scrollToElement(li[0], 0);
			myScroll.scrollTo(0, -(li.height() / 2), 0, true);// 上面留出半块儿
		}
		return done;
	}

	// 系统公告滚动显示
	function sysScroll() {
		var acBoxLeft = 0;
		var acBox_ul_w = $('.acBox ul').width();
		var $acBox_li = $('.acBox li');
		var acBox_li_w = $acBox_li.width();
		function acBoxScroll() {
			if (acBoxLeft < acBox_li_w) {
				acBoxLeft += 1;
			} else {
				acBoxLeft = -acBox_ul_w;
			}
			$acBox_li.css('left', -acBoxLeft);
			timeAc = setTimeout(acBoxScroll, 20);
		}
		if (acBox_li_w > acBox_ul_w) {
			timeAc = setTimeout(acBoxScroll, 1000);
		}
	}

	// 走马灯效果
	var getMsg = function() {
		var msg_url = js_Obj.msg_path;
		var data = {
			tt : $('body').attr('data-tt'),
			cpid : $('body').attr('data-cpid')
		}
		$.ajax({
			url : msg_url,
			cache : false,// 不保存缓存
			type : "get",
			data : data,
			success : function(data) {
				if (data.code != 0) {
					console.log("getMsg() request error,errorCode:" + data.code);
					return;
				}
				$('.acBox li').html("");
				var msg_var = "";
				for (var i = 0; i < data.data.length; i++) {
					msg_var += (data.data[i].adTitle);
				}
				if (msg_var) {
					$(".acBox.msg").show();
					$('.acBox li').html(msg_var);
					sysScroll();
				} else {
					$(".acBox.msg").hide();
				}
				listChgByMorePay(0, 0);
			},
			error : function(XMLHttpRequest, textStatus, errorThrown) {
				console.log("getMsg() error:" + XMLHttpRequest.status + "," + textStatus + "," + errorThrown);
				return;
			}
		});
	}
	// 对第一个进行click调用
	$('.paymentList li:not(.disabledPay):first').click();
	// window.addEventListener("DOMContentLoaded", loaded, false);
	// 判断R用户
	if (isRUser($('body').attr("data-userId"))) {
		window.location.href = goPersonCenter();
		return;
	}
	try {
		getMsg();
	} catch (e) {
		console.log("getMsg() error:" + e);
	}
	// 自动调起某支付方式
	var selected = null;
	if ($('body').attr('data-autoPay')) {
		if (toExpand("vc")) {
			selected = "vc";
			var cookieName = "home_cookie" + $('body').attr('data-userId') + "_" + $('body').attr('data-tid');
			var cook = ipay.utils.cookie.get(cookieName);
			if (!cook) {
				ipay.utils.cookie.set(cookieName, "1");
				$(".btnPay").click();
			}
		}
	}
	if (!selected) {
		if (ipay['utils'].browser.isWX()) {
			if (toExpand("wx")) {
				selected = "wx";
			}
		}
	}
});
function getPaytype() {
	if (ipay.utils.isWXBrowser()) {
		return "wx";
	} else if (/QQ\/\d/.test(navigator.userAgent)) {
		return "qq";
	}
	// 支付宝
	else if (/AlipayClient/.test(navigator.userAgent)) {
		return "ali";
	}
}
// 针对百度浏览器的处理
$(function() {
	if (self != top && /baidubrowser/i.test(navigator.userAgent)) {
		$("a").click(function() {
			if ($(this).attr("href")) {
				self.location.href = $(this).attr("href");
			}
			return false;
		});
	}
});