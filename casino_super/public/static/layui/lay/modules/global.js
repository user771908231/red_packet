/**

 layui官网

 */

layui.define(['layer', 'code', 'form', 'element', 'util'], function(exports){
    var $ = layui.jquery
        ,layer = layui.layer
        ,form = layui.form()
        ,util = layui.util
        ,device = layui.device();

    //阻止IE7以下访问
    if(device.ie && device.ie < 8){
        layer.alert('Layui最低支持ie8，您当前使用的是古老的 IE'+ device.ie + '，你丫的肯定不是程序猿！');
    }


    //首页banner
    setTimeout(function(){
        $('.site-zfj').addClass('site-zfj-anim');
        setTimeout(function(){
            $('.site-desc').addClass('site-desc-anim')
        }, 5000)
    }, 100);


    for(var i = 0; i < $('.adsbygoogle').length; i++){
        (adsbygoogle = window.adsbygoogle || []).push({});
    }




    //展示当前版本
    $('.site-showv').html(layui.v);

    //获取下载数
    $.get('http://fly.layui.com/api/handle?id=10&type=find', function(res){
        $('.site-showdowns').html(res.number);
    }, 'jsonp');

    //记录下载
    $('.site-down').on('click',function(){
        $.get('http://fly.layui.com/api/handle?id=10');
    });

    // //固定Bar
    // util.fixbar({
    //     bar1: true
    //     ,click: function(type){
    //         if(type === 'bar1'){
    //             location.href = 'http://fly.layui.com/';
    //         }
    //     }
    // });

    //窗口scroll
    ;!function(){
        var main = $('.site-tree').parent(), scroll = function(){
            var stop = $(window).scrollTop();
            if($(window).width() <= 750) return;
            var bottom = $('.footer').offset().top - $(window).height();
            if(stop > 61 && stop < bottom){
                if(!main.hasClass('site-fix')){
                    main.addClass('site-fix');
                }
                if(main.hasClass('site-fix-footer')){
                    main.removeClass('site-fix-footer');
                }
            } else if(stop >= bottom) {
                if(!main.hasClass('site-fix-footer')){
                    main.addClass('site-fix site-fix-footer');
                }
            } else {
                if(main.hasClass('site-fix')){
                    main.removeClass('site-fix').removeClass('site-fix-footer');
                }
            }
            stop = null;
        };
        scroll();
        $(window).on('scroll', scroll);
    }();

    //代码修饰
    layui.code({
        elem: 'pre'
    });

    //目录
    var siteDir = $('.site-dir');
    if(siteDir[0] && $(window).width() > 750){
        layer.open({
            type: 1
            ,content: siteDir
            ,skin: 'layui-layer-dir'
            ,area: 'auto'
            ,title: '目录'
            ,offset: 'r'
            ,shade: false
            ,success: function(layero, index){
                layer.style(index, {
                    marginLeft: -15
                })
            }
        });
        siteDir.find('li').on('click', function(){
            var othis = $(this);
            othis.find('a').addClass('layui-this');
            othis.siblings().find('a').removeClass('layui-this');
        });
    }

    //在textarea焦点处插入字符
    var focusInsert = function(str){
        var start = this.selectionStart
            ,end = this.selectionEnd
            ,offset = start + str.length

        this.value = this.value.substring(0, start) + str + this.value.substring(end);
        this.setSelectionRange(offset, offset);
    };

    //演示页面
    $('body').on('keydown', '#LAY_editor, .site-demo-text', function(e){
        var key = e.keyCode;
        if(key === 9 && window.getSelection){
            e.preventDefault();
            focusInsert.call(this, '  ');
        }
    });

    var iframe = $('#LAY_demo').prop('contentWindow');
    $('#LAY_demo_run').on('click', function(){
        iframe.layui.jquery('body').html($('#LAY_editor').val());
        if(iframe.layui.form){
            iframe.layui.form().render();
        }
        if(iframe.layui.element){
            iframe.layui.element().init();
        }
    });


    //查看代码
    $(function(){
        var DemoCode = $('#LAY_democode');
        DemoCode.val([
            DemoCode.val()
            ,'<body>'
            ,global.preview
            ,'\n<script src="//res.layui.com/layui/build/layui.js" charset="utf-8"></script>'
            ,$('#LAY_democodejs').html()
            ,'\n</body>\n</html>'
        ].join(''));
    });


    //手机设备的简单适配
    var treeMobile = $('.site-tree-mobile')
        ,shadeMobile = $('.site-mobile-shade')

    treeMobile.on('click', function(){
        $('body').addClass('site-mobile');
    });

    shadeMobile.on('click', function(){
        $('body').removeClass('site-mobile');
    });

    exports('global', {});
});