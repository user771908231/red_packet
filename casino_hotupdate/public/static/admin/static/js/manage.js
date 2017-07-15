$(document).ready(function () {
    layui.config({
        base: "/static/admin/plugins/layui/lay/modules/"
    });
    layui.use(["laytpl","layer", "laypage"],function () {
        layui.laytpl.config({
            open: "<{",
            close: "}>"
        });
    });
});

common = {
    msg: function (title, content) {
        layui.layer.open({
            title: title,
            content: content
        });
    },
    render: function (data, tpl_obj, view_obj) {
        layui.laytpl($(tpl_obj).html()).render(data, function (html) {
            $(view_obj).html(html)
        })
    },
    renderLayer: function (data, tpl_obj) {       //渲染一个弹出框
        var width = $(tpl_obj).data("width");
        var height = $(tpl_obj).data("height");
        layui.laytpl($(tpl_obj).html()).render(data, function (html) {
            layui.layer.open({
                type: 1 //Page层类型
                ,area: [width, height]
                ,title: $(tpl_obj).data("title")
                ,shade: 0.6 //遮罩透明度
                ,maxmin: true //允许全屏最小化
                ,anim: 1 //0-6的动画形式，-1不开启
                ,content: html
            });
        });
    },
    bindAjax: function (bind_event, bind_dom, event_func, ajax_func) {
        $(document).on(bind_event, bind_dom, function (event) {
            var event_obj = this;
            if(event_func.call(event_obj, event)){
                //处理表单事件
                if(bind_event == "submit"){
                    $.ajax({
                        url: $(event_obj).attr("action"),
                        type: $(event_obj).attr("method"),
                        data: $(event_obj).serialize(),
                        async: true,
                        success: function (res) {
                            if (ajax_func) {
                                ajax_func.bind(event_obj);
                                ajax_func.call(event_obj ,$(event_obj).serialize(), res)
                            }
                        },
                        error: function () {
                            layui.layer.open({title: "消息", content: "网络异常！"});
                        }
                    });
                }
                //处理超链接
                if(bind_event == "click"){
                    $.ajax({
                        url: $(event_obj).attr("href"),
                        type: "GET",
                        async: true,
                        success: function (res) {
                            if(ajax_func){
                                ajax_func.bind(event_obj);
                                ajax_func.call(event_obj, $(event_obj).attr("href"), res)
                            }
                        },
                        error: function () {
                            layui.layer.open({title: "消息", content: "网络异常！"});
                        }
                    });
                }
            }
        })
    },
    ajaxForm: function (form_obj, data, ajax_func) {
        var form = $(form_obj);
        layui.each(data, function (key, value) {
            form.find("input[name='"+key+"']").val(value)
        });
        $.ajax({
            url: form.attr("action"),
            type: form.attr("method"),
            data: form.serialize(),
            async: false,
            success: function (res) {
                if (ajax_func) {
                    ajax_func(res)
                }
            },
            error: function () {
                layui.layer.open({title: "消息", content: "网络异常！"});
            }
        });
    },
    getFormData: function (form_obj) {
        var form = $(form_obj).serializeArray();
        var data = {};
        for(var i=0;i<form.length;i++){
            data[form[i].name] = form[i].value
        }
        return data;
    }
};