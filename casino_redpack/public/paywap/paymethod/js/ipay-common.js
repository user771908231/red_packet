if (!String.prototype.trim) {
	String.prototype.trim = function() {
		return this.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, '');
	};
}

var ipay = {};
ipay['cfg'] = {};
ipay['cfg']['icon_loading'] = '/paywap/paymethod/img/loading.gif';

/*
 * iapppay公共方法
 */
/*
 * 公共方法1： ! maskUI plugin Version 1.0-2016.09.12 Requires jQuery v1.7 or later
 * 
 * Copyright (c) 2007-2013 M. Alsup Dual licensed under the MIT and GPL
 * licenses: http://www.opensource.org/licenses/mit-license.php
 * http://www.gnu.org/licenses/gpl.html
 * 
 * Thanks to Amir-Hossein Sobhi for some excellent contributions!
 * 
 * changed by zhuhaishen@20160912
 */

;
(function() {
	/* jshint eqeqeq:false curly:false latedef:false */
	"use strict";

	function setup($) {
		$.fn._fadeIn = $.fn.fadeIn;

		var noOp = $.noop || function() {
		};

		// this bit is to ensure we don't call setExpression when we shouldn't
		// (with extra muscle to handle
		// confusing userAgent strings on Vista)
		var msie = /MSIE/.test(navigator.userAgent);
		var ie6 = /MSIE 6.0/.test(navigator.userAgent) && !/MSIE 8.0/.test(navigator.userAgent);
		var mode = document.documentMode || 0;
		var setExpr = $.isFunction(document.createElement('div').style.setExpression);

		// global $ methods for blocking/unblocking the entire page
		$.maskUI = function(opts) {
			install(window, opts);
		};
		$.unmaskUI = function(opts) {
			remove(window, opts);
		};

		$.maskUI.version = 1.0; // 2nd generation blocking at no extra cost!
		// override these in your code to change the default behavior and style
		$.maskUI.defaults = {
			id : 'def',
			// message displayed when blocking (use null for no message)
			message : '<h1>请稍后...</h1>',
			css : {
				padding : 0,
				margin : 0,
				width : '70%',
				top : '40%',
				left : '15%',
				textAlign : 'center',
				color : '#000',
				border : '3px solid #aaa',
				// backgroundColor : '#fff',
				cursor : 'wait'
			},

			// styles for the overlay
			overlayCSS : {
				backgroundColor : '#000',
				opacity : 0.6,
				cursor : 'wait'
			},

			// style to replace wait cursor before unblocking to correct issue
			// of lingering wait cursor
			cursorReset : 'default',

			// z-index for the blocking overlay
			baseZ : 1000,

			// set these to true to have the message automatically centered
			centerX : true,
			// <-- only effects element blocking (page block
			// controlled via css above)
			centerY : true,

			// allow body element to be stetched in ie6; this makes blocking
			// look better
			// on "short" pages. disable if you wish to prevent changes to the
			// body height
			allowBodyStretch : true,

			// enable if you want key and mouse events to be disabled for
			// content that is blocked
			bindEvents : true,

			// be default maskUI will supress tab navigation from leaving
			// blocking content
			// (if bindEvents is true)
			constrainTabKey : true,

			// fadeIn time in millis; set to 0 to disable fadeIn on block
			fadeIn : 200,

			// fadeOut time in millis; set to 0 to disable fadeOut on unblock
			fadeOut : 400,

			// time in millis to wait before auto-unblocking; set to 0 to
			// disable auto-unblock
			timeout : 0,

			// disable if you don't want to show the overlay
			showOverlay : true,

			// if true, focus will be placed in the first available input field
			// when
			// page blocking
			focusInput : true,

			// elements that can receive focus
			focusableElements : ':input:enabled:visible',

			// suppresses the use of overlay styles on FF/Linux (due to
			// performance issues with opacity)
			// no longer needed in 2012
			// applyPlatformOpacityRules: true,
			// callback method invoked when fadeIn has completed and blocking
			// message is visible
			onBlock : null,

			// callback method invoked when unblocking has completed; the
			// callback is
			// passed the element that has been unblocked (which is the window
			// object for page
			// blocks) and the options that were passed to the unblock call:
			// onUnblock(element, options)
			onUnblock : null,

			// callback method invoked when the overlay area is clicked.
			// setting this will turn the cursor to a pointer, otherwise cursor
			// defined in overlayCss will be used.
			onOverlayClick : null,

			// don't ask; if you really must know:
			// http://groups.google.com/group/jquery-en/browse_thread/thread/36640a8730503595/2f6a79a77a78e493#2f6a79a77a78e493
			quirksmodeOffsetHack : 4,

			// class name of the message block
			blockMsgClass : 'blockMsg'
		};

		// private data and functions follow...
		var pageBlock = null;
		var pageBlockEls = [];

		function install(el, opts) {
			var css;
			var full = (el == window);
			var msg = (opts && opts.message !== undefined ? opts.message : undefined);
			opts = $.extend({}, $.maskUI.defaults, opts || {});
			var pre = getPre(opts);
			if ($(el).data(pre + 'isBlocked')) {
				return;
			}
			$(el).data(pre + 'isBlocked', 1);

			opts.overlayCSS = $.extend({}, $.maskUI.defaults.overlayCSS, opts.overlayCSS || {});
			css = $.extend({}, $.maskUI.defaults.css, opts.css || {});
			if (opts.onOverlayClick)
				opts.overlayCSS.cursor = 'pointer';

			msg = msg === undefined ? opts.message : msg;

			// if an existing element is being used as the blocking content then
			// we capture
			// its current place in the DOM (and current display style) so we
			// can restore
			// it when we unblock
			if (msg && typeof msg != 'string' && (msg.parentNode || msg.jquery)) {
				var node = msg.jquery ? msg[0] : msg;
				var data = {};
				$(el).data(pre + 'history', data);
				data.el = node;
				data.parent = node.parentNode;
				data.display = node.style.display;
				data.position = node.style.position;
				if (data.parent)
					data.parent.removeChild(node);
			}

			$(el).data(pre + 'onUnblock', opts.onUnblock);
			var z = opts.baseZ;

			// maskUI uses 3 layers for blocking, for simplicity they are all
			// used on every platform;
			// layer1 is the iframe layer which is used to supress bleed through
			// of underlying content
			// layer2 is the overlay layer which has opacity and a wait cursor
			// (by default)
			// layer3 is the message content that is displayed while blocking
			var iframeSrc = /^https/i.test(window.location.href || '') ? 'javascript:false' : 'about:blank';
			var lyr1, lyr2, lyr3, s;
			if (msie || opts.forceIframe)
				lyr1 = $('<iframe data-id=\'' + opts.id + '\' class="maskUI" style="z-index:' + (z++) + ';display:none;border:none;margin:0;padding:0;position:absolute;width:100%;height:100%;top:0;left:0" src="' + iframeSrc + '"></iframe>');
			else
				lyr1 = $('<div data-id=\'' + opts.id + '\' class="maskUI" style="display:none"></div>');

			lyr2 = $('<div data-id=\'' + opts.id + '\' class="maskUI blockOverlay" style="z-index:' + (z++) + ';display:none;border:none;margin:0;padding:0;width:100%;height:100%;top:0;left:0"></div>');

			if (full) {
				s = '<div data-id=\'' + opts.id + '\' class="maskUI ' + opts.blockMsgClass + ' blockPage" style="z-index:' + (z + 10) + ';display:none;position:fixed"></div>';
			} else {
				s = '<div data-id=\'' + opts.id + '\' class="maskUI ' + opts.blockMsgClass + ' blockElement" style="z-index:' + (z + 10) + ';display:none;position:absolute"></div>';
			}
			lyr3 = $(s);

			// if we have a message, style it
			if (msg) {
				lyr3.css(css);
			}

			// style the overlay
			lyr2.css(opts.overlayCSS);
			lyr2.css('position', full ? 'fixed' : 'absolute');

			// make iframe layer transparent in IE
			if (msie || opts.forceIframe)
				lyr1.css('opacity', 0.0);

			// $([lyr1[0],lyr2[0],lyr3[0]]).appendTo(full ? 'body' : el);
			var layers = [ lyr1, lyr2, lyr3 ], $par = full ? $('body') : $(el);
			$.each(layers, function() {
				this.appendTo($par);
			});

			// ie7 must use absolute positioning in quirks mode and to account
			// for activex issues (when scrolling)
			var expr = setExpr && (!$.support.boxModel || $('object,embed', full ? null : el).length > 0);
			if (ie6 || expr) {
				// give body 100% height
				if (full && opts.allowBodyStretch && $.support.boxModel)
					$('html,body').css('height', '100%');

				// fix ie6 issue when blocked element has a border width
				if ((ie6 || !$.support.boxModel) && !full) {
					var t = sz(el, 'borderTopWidth'), l = sz(el, 'borderLeftWidth');
					var fixT = t ? '(0 - ' + t + ')' : 0;
					var fixL = l ? '(0 - ' + l + ')' : 0;
				}

				// simulate fixed position
				$.each(layers, function(i, o) {
					var s = o[0].style;
					s.position = 'absolute';
					if (i < 2) {
						if (full)
							s.setExpression('height', 'Math.max(document.body.scrollHeight, document.body.offsetHeight) - (jQuery.support.boxModel?0:' + opts.quirksmodeOffsetHack + ') + "px"');
						else
							s.setExpression('height', 'this.parentNode.offsetHeight + "px"');
						if (full)
							s.setExpression('width', 'jQuery.support.boxModel && document.documentElement.clientWidth || document.body.clientWidth + "px"');
						else
							s.setExpression('width', 'this.parentNode.offsetWidth + "px"');
						if (fixL)
							s.setExpression('left', fixL);
						if (fixT)
							s.setExpression('top', fixT);
					} else if (opts.centerY) {
						if (full)
							s.setExpression('top', '(document.documentElement.clientHeight || document.body.clientHeight) / 2 - (this.offsetHeight / 2) + (blah = document.documentElement.scrollTop ? document.documentElement.scrollTop : document.body.scrollTop) + "px"');
						s.marginTop = 0;
					} else if (!opts.centerY && full) {
						var top = (opts.css && opts.css.top) ? parseInt(opts.css.top, 10) : 0;
						var expression = '((document.documentElement.scrollTop ? document.documentElement.scrollTop : document.body.scrollTop) + ' + top + ') + "px"';
						s.setExpression('top', expression);
					}
				});
			}

			// show the message
			if (msg) {
				lyr3.append(msg);
				if (msg.jquery || msg.nodeType)
					$(msg).show();
			}

			if ((msie || opts.forceIframe) && opts.showOverlay)
				lyr1.show(); // opacity is zero
			if (opts.fadeIn) {
				var cb = opts.onBlock ? opts.onBlock : noOp;
				var cb1 = (opts.showOverlay && !msg) ? cb : noOp;
				var cb2 = msg ? cb : noOp;
				if (opts.showOverlay)
					lyr2._fadeIn(opts.fadeIn, cb1);
				if (msg)
					lyr3._fadeIn(opts.fadeIn, cb2);
			} else {
				if (opts.showOverlay)
					lyr2.show();
				if (msg)
					lyr3.show();
				if (opts.onBlock)
					opts.onBlock.bind(lyr3)();
			}

			// bind key and mouse events
			bind(1, el, opts);

			if (full) {
				if (pageBlock) {
					var data = {};
					$(el).data(pre + 'pageBlock', data);
					data.pageBlock = pageBlock;
					data.pageBlockEls = pageBlockEls;
				}
				pageBlock = lyr3[0];
				pageBlockEls = $(opts.focusableElements, pageBlock);
				if (opts.focusInput)
					setTimeout(focus, 20);
			} else
				center(lyr3[0], opts.centerX, opts.centerY);

			if (opts.timeout) {
				// auto-unblock
				var to = setTimeout(function() {
					if (full)
						$.unmaskUI(opts);
					else
						$(el).unblock(opts);
				}, opts.timeout);
				$(el).data(pre + 'timeout', to);
			}
		}

		function getPre(opts) {
			return "maskUI." + opts.id + ".";
		}

		// remove the block
		function remove(el, opts) {
			opts = $.extend({}, $.maskUI.defaults, opts || {});
			var pre = getPre(opts);
			var count;
			var full = (el == window);
			var $el = $(el);
			if (!$el.data(pre + 'isBlocked')) {
				return;
			}
			$el.removeData(pre + 'isBlocked');

			var data = $el.data(pre + 'history');
			var to = $el.data(pre + 'timeout');
			if (to) {
				clearTimeout(to);
				$el.removeData(pre + 'timeout');
			}
			opts = $.extend({}, $.maskUI.defaults, opts || {});
			bind(0, el, opts); // unbind events
			if (opts.onUnblock === null) {
				opts.onUnblock = $el.data(pre + 'onUnblock');
				$el.removeData(pre + 'onUnblock');
			}

			var els = $("[data-id='" + opts.id + "']");
			// fix cursor issue
			if (opts.cursorReset) {
				if (els.length > 1)
					els[1].style.cursor = opts.cursorReset;
				if (els.length > 2)
					els[2].style.cursor = opts.cursorReset;
			}

			if (full)
				pageBlock = pageBlockEls = null;

			if (opts.fadeOut) {
				count = els.length;
				els.stop().fadeOut(opts.fadeOut, function() {
					if (--count === 0)
						reset(els, data, opts, el);
				});
			} else
				reset(els, data, opts, el);
		}

		// move blocking element back into the DOM where it started
		function reset(els, data, opts, el) {
			var pre = getPre(opts);
			var $el = $(el);

			els.each(function(i, o) {
				// remove via DOM calls so we don't lose event handlers
				if (this.parentNode)
					this.parentNode.removeChild(this);
			});

			if (data && data.el) {
				data.el.style.display = data.display;
				data.el.style.position = data.position;
				data.el.style.cursor = 'default'; // #59
				if (data.parent)
					data.parent.appendChild(data.el);
				$el.removeData(pre + 'history');
			}

			if ($el.data(pre + 'static')) {
				$el.css('position', 'static'); // #22
			}

			var pb = $el.data(pre + 'pageBlock');
			if (pb) {
				pageBlock = pb.pageBlock;
				pageBlockEls = pb.pageBlockEls;
				$el.removeData(pre + 'pageBlock');
			}

			if (typeof opts.onUnblock == 'function')
				opts.onUnblock(el, opts);
			$el.removeData(pre + 'isBlocked');
			// fix issue in Safari 6 where block artifacts remain until reflow
			var body = $(document.body), w = body.width(), cssW = body[0].style.width;
			body.width(w - 1).width(w);
			body[0].style.width = cssW;
		}

		// bind/unbind the handler
		function bind(b, el, opts) {
			var pre = getPre(opts);
			var full = el == window, $el = $(el);

			// don't bother unbinding if there is nothing to unbind
			if (!b && (full && !pageBlock || !full && !$el.data(pre + 'isBlocked'))) {
				return;
			}

			$el.data(pre + 'isBlocked', b);

			// don't bind events when overlay is not in use or if bindEvents is
			// false
			if (!full || !opts.bindEvents || (b && !opts.showOverlay))
				return;

			// bind anchors and inputs for mouse and key events
			var events = 'mousedown mouseup keydown keypress keyup touchstart touchend touchmove';
			if (b)
				$(document).bind(events, opts, handler);
			else
				$(document).unbind(events, handler);

			// former impl...
			// var $e = $('a,:input');
			// b ? $e.bind(events, opts, handler) : $e.unbind(events, handler);
		}

		// event handler to suppress keyboard/mouse events when blocking
		function handler(e) {
			// allow tab navigation (conditionally)
			if (e.type === 'keydown' && e.keyCode && e.keyCode == 9) {
				if (pageBlock && e.data.constrainTabKey) {
					var els = pageBlockEls;
					var fwd = !e.shiftKey && e.target === els[els.length - 1];
					var back = e.shiftKey && e.target === els[0];
					if (fwd || back) {
						setTimeout(function() {
							focus(back);
						}, 10);
						return false;
					}
				}
			}
			var opts = e.data;
			var target = $(e.target);
			if (target.hasClass('blockOverlay') && opts.onOverlayClick)
				opts.onOverlayClick(e);

			// allow events within the message content
			if (target.parents('div.' + opts.blockMsgClass).length > 0)
				return true;

			// allow events for content that is not being blocked
			return target.parents().children().filter('div.maskUI').length === 0;
		}

		function focus(back) {
			if (!pageBlockEls)
				return;
			var e = pageBlockEls[back === true ? pageBlockEls.length - 1 : 0];
			if (e)
				e.focus();
		}

		function center(el, x, y) {
			var p = el.parentNode, s = el.style;
			var l = ((p.offsetWidth - el.offsetWidth) / 2) - sz(p, 'borderLeftWidth');
			var t = ((p.offsetHeight - el.offsetHeight) / 2) - sz(p, 'borderTopWidth');
			if (x)
				s.left = l > 0 ? (l + 'px') : '0';
			if (y)
				s.top = t > 0 ? (t + 'px') : '0';
		}

		function sz(el, p) {
			return parseInt($.css(el, p), 10) || 0;
		}

	}

	/* global define:true */
	if (typeof define === 'function' && define.amd && define.amd.jQuery) {
		define([ 'jquery' ], setup);
	} else {
		setup(jQuery);
	}

})();

/**
 * 公共方法2.<BR>
 * 执行ajax锁屏：ajax4lock<BR>
 * 执行某个方法时锁屏：exec4lock
 */
(function() { // jquery相关JS
	$.extend({
		// 全局函数定义
		// ID生成器，用于生成唯一标识ID，目前用于mask标签的标识
		nexTagID : function(pre) {
			return (pre ? pre : 'k') + '_' + new Date().getTime() + "_" + Math.floor(Math.random() * 100000000);
		},
		// ajax封装
		ajax4lock : function(setting) {
			var msg = '正在加载...';
			if (setting.lock && setting.lock.msg) {
				msg = setting['lock']['msg'];
			}
			var maskOpts = {
				id : $.nexTagID(),
				message : '<span style="width:33px; height:33px; background:url(' + ipay['cfg']['icon_loading'] + ') no-repeat; background-size:100%; position:relative;  display:block;margin: 20px auto 10px auto;">&nbsp;</span><span style="text-align: center; color:#eee; font-size:14px;">' + msg + '</span>',
				overlayCSS : {
					backgroundColor : 'rgba(0,0,0,0)'
				},
				css : {
					// backgroundColor : 'rgba(0,0,0,0.5)',
					border : 'none',
					// color : '#FFF',
					// height : '40px',
					// 'line-height' : '40px',
					// 'border-radius' : '10px',
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
			};
			var newSetting = $.extend({}, setting);
			newSetting['beforeSend'] = function(xhr) {
				$.maskUI(maskOpts);
				try {
					if (typeof setting.beforeSend == 'function') {
						setting.beforeSend(xhr);
					}
				} catch (e) {
					try {
						console.log(e);
					} catch (ex) {
					}
				}
			};
			newSetting['complete'] = function(xhr, status) {
				try {
					if (typeof setting.complete == 'function') {
						setting.complete(xhr, status);
					}
				} catch (e) {
					try {
						console.log(e);
					} catch (ex) {
					}
				}
				$.unmaskUI(maskOpts);
			};
			$.ajax(newSetting);
		},

		// 执行时锁定屏幕。
		exec4lock : function(run) {
			var maskOpts = {
				id : $.nexTagID(),
				message : null,
				overlayCSS : {
					backgroundColor : 'rgba(0,0,0,0)'
				}
			};
			try {
				$.maskUI(maskOpts);
				return run();
			} finally {
				$.unmaskUI(maskOpts);
			}
		}
	});

	$.fn.extend({
		// 执行时锁定屏幕。
		exec4lock : function(run) {
			var maskOpts = {
				id : $.nexTagID(),
				message : null,
				overlayCSS : {
					backgroundColor : 'rgba(0,0,0,0)'
				}
			};
			try {
				$.maskUI(maskOpts);
				return run();
			} finally {
				$.unmaskUI(maskOpts);
			}
		}
	});

})();

/**
 * 公共方法3：Formater 使用示例。var formater = new Formater('{name1} is {name2}'); var
 * out = formater.exec({ name1: 'Iapppay', name2: 'a family' }); then, out is
 * "Iapppay is a family"。
 */
var Formater = function(template, patten) {
	this._template = template;
	this._patten = patten || /{([^{}]+)}/gm;
	this._store = [];
	this._data = {};
	var _this = this;
	this.exec = exec;
	this.build = build;
	this.build();
	return this;

	function exec(collection) {
		this._data = collection || {};
		return this._store.join("");
	}

	function build(newTemplate, newPatten) {
		if (newTemplate) {
			this._template = newTemplate;
		}
		if (newPatten) {
			this._patten = newPatten;
		}
		if (this._template) {
			this._store.length = 0;
			var tmp = this._template, lastLastIndex = 0, flags = (this._patten.ignoreCase ? "i" : "") + (this._patten.multiline ? "m" : "") + (this._patten.sticky ? "y" : ""), separator = RegExp(this._patten.source, flags + "g");
			while (match = separator.exec(tmp)) {
				var lastIndex = match.index + match[0].length; // `separator.lastIndex`
				// is not
				// reliable
				// cross-browser
				var s = tmp.slice(lastLastIndex, match.index);
				if (s != '') {
					this._store.push(s);
				}
				lastLastIndex = lastIndex;
				if (match) {
					this._store.push(new struct(match[1]));
				}
				if (match && separator.lastIndex === match.index)
					separator.lastIndex++; // avoid an infinite loop
			}
			if (lastLastIndex < tmp.length) {
				this._store.push(tmp.slice(lastLastIndex));
			}
		}
	}

	function struct(str) {
		this.str = str.trim();
		this.toString = function() {
			return _this._data[this.str] || "";
		}
	}

};

/**
 * 公共方法区
 */
var _ipay_rquery = (/\?/);
var _ipay_rts = /([?&])_=[^&]*/;
ipay['utils'] = {
	// 生成一个GET url
	toGetURL : function(url, params) {
		if (params && typeof params !== 'string') {
			params = $.param(params);
		}
		return params ? url += (_ipay_rquery.test(url) ? "&" : "?") + params : url;
	},
	// 获取一个刷新url
	refreshURL : function(url) {
		return _ipay_rts.test(url) ?
		// If there is already a '_' parameter, set its value
		url.replace(_ipay_rts, "$1_=" + new Date().getTime()) :
		// Otherwise add one to the end
		url + (_ipay_rquery.test(url) ? "&" : "?") + "_=" + new Date().getTime();
	},
	// 刷新当前页面
	reload : function() {
		window.location.reload(true);
	},
	formater : function(ex) {
		return new Formater(ex);
	},
	log : function(msg, e) {
		if (msg) {
			console.log(msg);
		}
		if (e) {
			console.log(e);
		}
	},
	findMethod : function(baseMethodName, key) {
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
	},
	isWXBrowser : function() {
		return /MicroMessenger/.test(navigator.userAgent);
	},

	isInIframe : function() {
		return self != top;
	},
	browser : {
		// 当前页面是否在iframe中
		isInIframe : function() {
			return self != top;
		},
		// 是否为微信浏览器
		isWX : function() {
			return /MicroMessenger/.test(navigator.userAgent);
		},
		// 是否为QQ浏览器
		isQQ : function() {
			return /MQQBrowser/.test(navigator.userAgent);
		},
		// 是否为iOS系统
		isIOS : function() {
			return /Mac OS X/.test(navigator.userAgent) || /iOS/.test(navigator.userAgent);
		},
		// 是否为Android系统
		isAndroid : function() {
			return /Android/.test(navigator.userAgent);
		},
		// 是否为爱贝APP
		isHappyApp : function() {
			return /For IAppPay/i.test(navigator.userAgent) || /4HPYP/i.test(navigator.userAgent);
		},
		getUa : function() {
			var p_mac = /Mac OS/i;
			var p_version = /Version|Mobile/i;
			var p_safari = /Safari/i;
			var p_chrome = /Chrome|CriOS/i;
			var p_qh360 = /QHBrowser/i;
			var p_samsung = /SamsungBrowser/i;
			var ua = navigator.userAgent;
			var uax = [];
			// $('#ua').html(ua);
			var uaxs = ua.replace(/.*\(KHTML.*Gecko\)(.*)/, '$1').trim();
			var p_no_chrome = /^Version/i;
			var mayChrome = true;
			if (p_no_chrome.test(uaxs)) {
				mayChrome = false;
			}
			uaxs = uaxs.split(' ');
			for (i in uaxs) {
				if (!p_version.test(uaxs[i])) {
					uax.push(uaxs[i]);
				}
			}
			uaxs = uax.join(' ');
			if (p_mac.test(ua)) {// iOS系统
				if (uax.length == 1) {
					if (p_safari.test(uax[0])) {
						return 'safari';
					}
				}
			}
			if (mayChrome && uax.length == 2) {
				if (p_chrome.test(uaxs) && p_safari.test(uaxs)) {
					return 'chrome';
				}
			}
			if (p_samsung.test(uaxs)) {
				return 'samsung';
			}
			if (p_qh360.test(uaxs)) {
				return 'qh360';
			}
		}

	},
	cookie : {
		// 写cookies
		set : function(name, value, time, path) {
			var exp = new Date();
			if (time) {
				exp.setTime(exp.getTime() + time * 1000);// 过期时间
			} else {
				exp.setTime(exp.getTime() + 60 * 60 * 1000);// 默认一小时
			}
			if (path) {
				document.cookie = name + "=" + escape(value) + ";Path=" + path + ";expires=" + exp.toGMTString();
			} else {
				document.cookie = name + "=" + escape(value) + ";expires=" + exp.toGMTString();
			}
		},
		// 读取cookies
		get : function(name) {
			var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
			if (arr = document.cookie.match(reg))
				return unescape(arr[2]);
			else
				return null;
		},
		// 删除cookies
		del : function(name, path) {
			if (path) {
				document.cookie = name + "=;Path=" + path + ";expires=Thu, 01 Jan 1970 00:00:00 GMT";
				return;
			}
			document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
		}
	},
	url : {
		isHttp : function(url) {
			return /^((http:\/\/)|(https:\/\/))/i.test(url);
		}
	}
}

/**
 * 通用ajax初始化配置。 配置的默认的error错误处理。及timeout。
 */
$(function() { // 页面加载后的初始化数据
	// ajax封装
	$.ajaxSetup({
		error : function(response, textStatus, errorThrown) {
			var status = response.status;
			if (status == 200) {
				alert('服务暂不可用(响应数据异常)!(200)');
				return;
			}
			if (500 <= status) {
				alert("服务器异常！(" + status + ")");
			} else {
				if (status >= 400 && status < 500) {
					alert("服务器拒绝服务！(" + status + ")");
				} else {
					alert("网络或服务异常(" + status + ")！请稍后再试！");
				}
			}
		},
		timeout : 15000
	});
});
