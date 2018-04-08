// 判断是否为电话号码
function isTelAvailable(poneInput) {
    var myreg= /^0?(13[0-9]|15[012356789]|17[013678]|18[0-9]|14[57])[0-9]{8}$/;
    if (!myreg.test(poneInput.val())) {
        return false;
    } else {
        return true;
    }
}
function isPhoneNo(phone) {
    var pattern = /^0?(13[0-9]|15[012356789]|17[013678]|18[0-9]|14[57])[0-9]{8}$/;
    return pattern.test(phone);
}


    // console.log($('.captcha_img').trigger("onclick").val());



function isPasswd(passwd_one,passwd_two) {

    if(passwd_one.length == 0){
        layer.msg("密码不能为空！")
        return false;
    }else if(passwd_one.length >12){
        layer.msg("密码长度不能超过12位！")
        return false;
    }else if(passwd_one.length < 6){
        layer.msg("密码长度最少12位！")
        return false;
    }

    if(passwd_two.length == 0){
        layer.msg("重复密码不能为空！")
        return false;
    }else if(passwd_two.length >12){
        layer.msg("重复密码长度不能超过12位！")
        return false;
    }else if(passwd_two.length < 6){
        layer.msg("重复密码长度最少6位！")
        return false;
    }

    if(passwd_one != passwd_two){
        layer.msg("两次密码输入不一致！")
        return false;
    }
    return true;
}

// $('#sbmit').click(function () {
//     var name_values = $('#name').val();
//     var passwd_one = $("#passwd_one").val();
//     var passwd_two = $('#passwd_two').val();
//     var captcha = $('#captcha').val();
//     console.log("name:"+name_values,"passwd_one:"+passwd_one,"passwd_two:"+passwd_two,"captcha:"+captcha)
//     if (name_values.length < 11) {
//         return layer.msg('手机号小于11位！');
//     }else if (name_values.length >11){
//         return layer.msg('手机号大于11位！');
//     }else if (name_values.length == 11) {
//         if (!isPhoneNo(name_values)) {
//             return layer.msg("不是一个有效手机号码")
//         }
//         isPasswd(passwd_one,passwd_two)
//         if (captcha.length == 0){
//             return layer.msg("请输入验证码！")
//         }
//
//         if(captcha.length != 4){
//             return layer.msg("验证码输入有误！")
//         }
//
//         // $.ajax({
//         //     type:'POST',
//         //     url:window.location,
//         //     data:{"name":name_values,"passwd_one":passwd_one,"passwd_two":passwd_two,"captcha":captcha},
//         //     success:function (data) {
//         //         console.log(data)
//         //     }
//         //
//         // })
//
//         $("#form-reg").submit();
//
//         common.ajaxForm("#form-reg", {}, function (data) {
//             console.log(data)
//         })
//
//     }else{
//         console.log("dfsdf")
//     }
// })

common.bindAjax("submit", "#form-reg", function (e) {
    e.preventDefault();
    var name_values = $('#name').val();
    var passwd_one = $("#passwd_one").val();
    var passwd_two = $('#passwd_two').val();
    var captcha = $('#captcha').val();
    if(name_values.length == 0) {
        layer.msg('请输入手机号！')
        return false;
    }
    if (name_values.length < 11) {
        layer.msg('手机号小于11位！');
        return false;
    }else if (name_values.length >11){
        layer.msg('手机号大于11位！');
        return false;
    }else if (name_values.length == 11) {
        if (!isPhoneNo(name_values)) {
            layer.msg("不是一个有效手机号码")
            return false;
        }
        if(!isPasswd(passwd_one, passwd_two)){
            return false;
        }

        if (captcha.length == 0) {
            layer.msg("请输入验证码！")
            return false;
        }

        if (captcha.length != 4) {
            layer.msg("验证码输入有误！")
            return false;
        }
        return true;
    }

}, function (req, res) {
    //返回提示！
    layer.msg(res.msg)
    if(res.code == 304){
        window.location.replace("/home");
}

        console.log($(".captcha_img"))
        // $('.captcha_img')[0].src='/captcha/0ajVxw3tiJCPSLa.png?reload='+(new Date()).getTime();
    $('.captcha_img').trigger("myClick")
    console.log("ddd")


});


