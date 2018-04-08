collectPasswd = function(options, onCollected, clickChange) {
	if (typeof options == 'function') {
		var tmp = options;
		options = onCollected;
		onCollected = tmp;
	}
	// 密码弹框
	var passOptions = {
		title : '',
		tip : '',
		err : null,
		numPad : '', // 支付密码的底部提示
		clickChange : '', // 点此修改
		price : '', // 支付密码价格提示
		showHide : true
	// 显示隐藏
	};
	var settings = $.extend({}, passOptions, options);
	this.id = $.nexTagID('cp');

	this.pop = $('<div class="'
			+ id
			+ ' alert-pwd-mask"></div><div class="'
			+ id
			+ ' alert-pwd-passWord" data-i="0">'
			+ '<div class="numPas" style="display: block"><div class="alert-pwd-password_top" style="text-align:center"><a href="javascript:void(0)" class="alert-pwd-close ac_Close" data-id="'
			+ id
			+ '"></a><p style="line-height: 18px;padding-top: 5px;">'
			+ settings.title
			+ '</p><p class="alert-pwd-price" style="color: #709ac4;text-align: center; font-size: 12px;line-height: 15px;">'
			+ settings.price
			+ '</p></div><form id="alert-pwd-password"> <input readonly="readonly" class="alert-pwd-pass alert-pwd-pass_left" type="password" maxlength="1" value=""> <em class="alert-pwd-pass-line" style="left: 42px;"> </em> <input readonly="readonly" class="alert-pwd-pass" type="password" maxlength="1" value=""><em class="alert-pwd-pass-line" style="left: 82px;"> </em> <input readonly="readonly" class="alert-pwd-pass" type="password" maxlength="1" value=""><em class="alert-pwd-pass-line" style="left: 124px;"> </em> <input readonly="readonly" class="alert-pwd-pass" type="password" maxlength="1" value=""> <em class="alert-pwd-pass-line" style="left: 164px;"> </em><input readonly="readonly" class="alert-pwd-pass" type="password" maxlength="1" value=""> <em class="alert-pwd-pass-line" style="left: 206px;"> </em><input readonly="readonly" class="alert-pwd-pass alert-pwd-pass_right" type="password" maxlength="1" value=""><em class="alert-pwd-pass-line" style="left: 246px;"> </em></form><p class="alert-pwd-popformErrBox"></p><p class="alert-pwd-password_footer numPad"  style="text-align: center; color:#767676;">'
			+ settings.numPad + '<span style="color:#709ac4">' + settings.tip + '</span> ' + settings.clickChange + ' </p></div>' + '<div class="noNumPad" style="display: none"><div class="alert-pwd-password_top" style="text-align:center"><a style="display: none;" href="javascript:void(0)" class="alert-pwd-close ac_Close" data-id="' + id + '"></a><p>' + settings.title
			+ '</p></div><div class="noNumPadDiv"><input class="passwordE" maxlength="24" style="border-radius: 8px; width: 100%;" type="password" value="" placeholder="请输入原支付密码" ></div><p class="alert-nonumpwd-popformErrBox"></p><div class="footerBtn"><a href="javascript:void(0);" class="ac_popCancel">取消</a><a href="javascript:void(0)" class="ac_popAccept">确定</a></div></div></div>');

	$('body').append(this.pop);
	var pass = $('.alert-pwd-pass');
	var passwordE = $('.passwordE');
	var _this = this;
	this.pop_passWord = $('.alert-pwd-passWord');
	this.pop_passWord.css({
		'margin-top' : -pop_passWord.height() / 2
	});
	this.pop_passWord.find('.alert-pwd-close').click(function() {
		_this.close();
	});
	this.pop_passWord.find('.ac_popCancel').click(function() {
		_this.close();
	});
	this.onCollected = onCollected;
	if (settings.showHide) {
		$('.numPas').show();
		$('.noNumPad').hide();
	} else {
		$('.numPas').hide();
		$('.noNumPad').show();
	}
	this.pop_passWord.find('.ac_popAccept').click(function() {
		var padval = passwordE.val();
		if (padval.length > 0) {
			$('.passwordE').val('');
			_this.onCollected(padval, _this);
		} else {
			$('.passwordE').select();
			$('.alert-nonumpwd-popformErrBox').html('请输入原支付密码')
		}
	});
	passwordE.focus();
	passwordE.on('keyup', function(e) {
		if (passwordE.val().length > 0) {
			$('.alert-nonumpwd-popformErrBox').html('')
		}
	});
	this.settings = settings;
	this.close = function() {
		$('.' + this.id).remove();
		$('#alert-key').remove();
	};
	this.errorTip = function(err) {
		$('.passwordE').focus();
		this.pop.find('.alert-nonumpwd-popformErrBox').html(err);
	};
	this.count = 6;
	this.reset = function(changeOptions, onCollected) {
		if (changeOptions) {
			this.settings = $.extend({}, this.settings, changeOptions);
			if (changeOptions.title) {
				this.pop.find('.alert-pwd-password_top p').html(changeOptions.title);
			}
			if (changeOptions.tip) {
				this.pop.find('.alert-pwd-password_footer span').html(changeOptions.tip);
			}
			if (changeOptions.numPad) {
				this.pop.find('.numPad').html(changeOptions.numPad);
			}
			if (changeOptions.clickChange) {
				this.pop.find('.clickChange').html(changeOptions.clickChange);
			}
			if (changeOptions.price) {
				this.pop.find('.alert-pwd-price').html(changeOptions.price);
			}
			if (changeOptions.err) {
				this.pop.find('.alert-pwd-popformErrBox').html(changeOptions.err);
			} else {
				this.pop.find('.alert-pwd-popformErrBox').html('');
			}
		}
		if (onCollected) {
			this.onCollected = onCollected;
		}
		this.pop.find('.numPas').show();
		this.pop.find('.noNumPad').hide();
		this.pop_passWord.attr('data-i', 0);
		keyboard(this.callback);
		for (var i = 0; i < this.count; i++) {
			$('.alert-pwd-pass:eq(' + i + ')').val('');
			$('.alert-pwd-pass:eq(' + i + ')').removeClass('current');
			$('.alert-pwd-pass:eq(' + i + ')').next().css({
				'display' : 'none'
			})
		}
	};

	this.callback = function(type, val) {
		var alertIndex = this.pop_passWord.attr('data-i');
		var passwords = $('#alert-pwd-password').get(0);
		if (type === 'char') {
			if (alertIndex < this.count) {
				$(passwords.elements[alertIndex]).val(val);
				$(passwords.elements[alertIndex]).addClass('current');
				$(passwords.elements[alertIndex]).next().css({
					'display' : 'block'
				});
				alertIndex++;
				this.pop_passWord.attr('data-i', alertIndex);
			}
			if (alertIndex >= this.count) {
				var temp_rePass_word = '';
				for (var j = 0; j < passwords.elements.length; j++) {
					temp_rePass_word += $(passwords.elements[j]).val();
				}
				this.onCollected(temp_rePass_word, this);
			}
		} else if (type === 'del') {
			if (alertIndex == 0) {
				return;
			}
			$(passwords.elements[alertIndex - 1]).val('');
			$(passwords.elements[alertIndex - 1]).removeClass('current');
			$(passwords.elements[alertIndex - 1]).next().css({
				'display' : 'none'
			});
			this.pop_passWord.attr('data-i', alertIndex - 1);
		}
	};
	// 调起键盘
	keyboard(this.callback);
	function keyboard(callback) {
		// 键盘弹框
		var div = '\
		<div id="alert-key" style="position:fixed;background-color:#A8A8A8;width:100%;bottom:0;z-index: 10;display: none;">\
			<ul id="alert-keyboard" style="font-size:20px;margin:2px -2px 1px 2px">\
				<li class="alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="off">1</span></li>\
				<li class="alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="off">2</span></li>\
				<li class="alert-symbol alert-btn_number_" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="alert-off">3</span></li>\
				<li class="alert-tab alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="alert-off">4</span></li>\
				<li class="alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="off">5</span></li>\
				<li class="alert-symbol alert-btn_number_" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="off">6</span></li>\
				<li class="alert-tab alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="alert-off">7</span></li>\
				<li class="alert-symbol" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="off">8</span></li>\
				<li class="alert-symbol alert-btn_number_" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;"><span class="alert-off">9</span></li>\
				<li class="alert-cancle alert-btn_number_"  unselectable="on" onselectstart="return false;" style="-moz-user-select:none;background: #bbbec3;margin-bottom: 0;"></li>\
				<li class="alert-symbol zero" unselectable="on" onselectstart="return false;" style="-moz-user-select:none;margin-bottom: 0;"><span class="alert-off">0</span></li>\
			</ul>\
			<p class=" alert-delete alert-lastitem" style="margin-bottom: 0;"></p>\
		</div>\
		';
		$('body').append(div);
		$('#alert-key').slideDown("slow");
		$('.alert-symbol').unbind('click');
		$('.alert-symbol').click(function() {
			var c = $(this).text();
			callback('char', c);
		});
		$('.alert-delete').unbind('click');
		$('.alert-delete').click(function() {
			callback('del');
		});
	}

	if ($(".noNumPad").css('display') !== 'none') {
		$('#alert-key').remove();
	} else {
		keyboard(this.callback);
	}
};

// 登录密码弹框
loginPassWord = function(options, onCollected) {
	// 密码弹框
	var passOptions = {
		title : '',
		tip : '',
		oneStyle : '',
		twoStyle : 'display:none',
		err : null
	};
	var settings = $.extend({}, passOptions, options);
	this.id = $.nexTagID('cp');

	this.pop = $('<div class="'
			+ id
			+ ' alert-pwd-mask"></div><div class="'
			+ id
			+ ' alert-pwd-passWord" data-i="0"><div class="alert-pwd-password_top" style="text-align:center"><p>'
			+ settings.title
			+ '</p></div><div class="old"><div class="onepwd" style="margin: 18px 25px 7px 25px;'
			+ settings.oneStyle
			+ '"><input class="oldPassword" style="border-radius: 8px; width: 100%;" type="password" value="" placeholder="请输入6-24位登录密码" maxlength="24"></div>'
			+ '<div class="graphVerifyCode" style="margin: 9px 95px 7px 25px;position: relative; display: none;"><input class="graphVerifyCodeInp" style="border-radius: 8px; width: 100%;" type="text" placeholder="请输入右图字符" name="verifyCode" maxlength="4" id="verifyCode"><span><img src=""></span></div>'
			+ '<p class="alert-pwd-popformErrBox"></p><div class="footerBtn"><a href="javascript:void(0);" class="ac_popCancel">取消</a><a href="javascript:void(0)" class="ac_popAccept oldAccept">确定</a></div></div>'
			+ '<div class="new" style="'
			+ settings.twoStyle
			+ '"><div class="twopwd" style="margin: 18px 25px 7px 25px;"><input class="newPassword passwordOne" style="border-radius: 8px; width: 100%;margin-bottom: 10px;" type="password" value="" placeholder="请输入新的6-24位登录密码" maxlength="24" ><input class="newPassword passwordTwo" style="border-radius: 8px; width: 100%;" type="password" value="" placeholder="请再次输入登录密码" maxlength="24" ></div><p class="alert-pwd-popformErrBox"></p><div class="footerBtn"><a href="javascript:void(0);" class="ac_popCancel">取消</a><a href="javascript:void(0)" class="ac_popAccept newAccept">确定</a></div></div></div>');
	$('body').append(this.pop);
	var oldPassword = $('.oldPassword');
	var oldAccept = $('.oldAccept');
	oldPassword.focus();
	this.pop_passWord = $('.alert-pwd-passWord');
	this.pop_passWord.css({
		'margin-top' : -pop_passWord.height() / 2
	});

	this.onCollected = onCollected;
	var _this = this;
	this.pop_passWord.find('.ac_popCancel').click(function() {
		_this.close();
	});
	oldAccept.click(function() {
		var padval = oldPassword.val();
		if (padval.length == 0) {
			oldPassword.focus();
			$('.alert-pwd-popformErrBox').html('密码不能为空')
		} else if (padval.length > 5) {
			$('.oldPassword').val('');
			_this.onCollected(padval, _this);
			$('.passwordOne').focus();
		} else {
			$('.alert-pwd-popformErrBox').html('请输入正确6-24位登录密码')
		}
	});
	$("input").focus(function() {
		$('.alert-pwd-popformErrBox').html('')
	});
	// 显示图形验证码
	this.showGraphVerifyCode = function(url) {
		$('.graphVerifyCode img').attr('src', url);
		$('.graphVerifyCode').show();
	};
	this.close = function() {
		$('.' + this.id).remove();
		$('#alert-key').remove();
	};
	this.errorTip = function(err) {
		this.pop.find('.alert-pwd-popformErrBox').html(err);
	};
	this.count = 6;
	this.reset = function(onCollected) {
		if (onCollected) {
			this.onCollected = onCollected;
		}

		var oldpwd = this.pop.find('.old');
		var newpaw = this.pop.find('.new');
		oldpwd.hide();
		newpaw.show();
		var passwordOne = this.pop.find('.passwordOne');
		var passwordTwo = this.pop.find('.passwordTwo');

		this.pop.find('.newAccept').click(function() {
			if (passwordOne.val().length == 0) {
				$('.passwordOne').focus();
				$('.alert-pwd-popformErrBox').html('密码不能为空')
			} else if (passwordOne.val().length < 6) {
				$('.passwordOne').focus();
				$('.alert-pwd-popformErrBox').html('请输入正确6-24位登录密码')
			} else if (passwordOne.val() == passwordTwo.val()) {
				_this.onCollected(passwordOne.val(), _this);
			} else {
				$('.passwordOne').focus();
				$('.newPassword').val('');
				$('.alert-pwd-popformErrBox').html('密码输入不一致请重新输入')
			}
		});
	};
};