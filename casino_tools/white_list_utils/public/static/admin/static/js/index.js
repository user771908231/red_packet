/**
 * Created by david on 2017/1/13.
 */
layui.config({
    base: '/static/admin/static/modules/'
}).use(['jquery', 'tab','navbar', 'layer'], function () {
    var $ = layui.jquery,
        layer = layui.layer,
        navbar = layui.navbar(),
        tab = layui.tab({
            elem: '.admin-nav-card' //设置选项卡容器
        });



    // 顶级菜单 点击切换侧边导航
    $('.myTabs>li').click(function () {
        var ids = $(this).attr('aria-controls');

        $('.nav_wraper ul').css('display', 'none');
        $('#' + ids).css('display', 'block');

        $('.myTabs li').removeClass('active');
        $('.' + ids).addClass('active');
    });

    // 点击用户
    $('#infoDesc').hover(function () {
        $(this).find('ul').show();
    }, function () {
        $(this).find('ul').hide();
    });

    //菜单
    // $('.menu>li').on('click', function () {
    //     $(this).find('.submenu').toggle(300);
    // });





    //iframe自适应
    $(window).on('resize', function () {
        var $content = $('.admin-nav-card .layui-tab-content');
        $content.height($(this).height() - 147);
        $content.find('iframe').each(function () {
            $(this).height($content.height());
        });
    }).resize();

    //设置navbar
    navbar.set({
        spreadOne: false,
        elem: '#admin-navbar-side',
        cached:false,
        url: '/static/admin/datas/nav.json'
    });
    //渲染navbar
    navbar.render();
    //监听点击事件
    navbar.on('click(side)', function(data) {
        tab.tabAdd(data.field);
    });

    // $('.openapp').on('click', function () {
    //     var url = $(this).attr('data-href');
    //     var icon = $(this).attr('data-icon');
    //     var title = $(this).attr('data-title');
    //     openTab(url, title, icon);
    // });

    // function openTab(url, title, icon) {
    //     if (!url || !title) {
    //         layer.msg('参数不能为空');
    //         return false;
    //     }
    //     var str = {
    //         'href': url,
    //         'icon': icon,
    //         'title': title
    //     };
    //     tab.tabAdd(str);
    // };



});