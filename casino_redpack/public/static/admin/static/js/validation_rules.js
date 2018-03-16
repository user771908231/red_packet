// 判断是否为电话号码
function isTelAvailable($poneInput) {
    var myreg= /^0?(13[0-9]|15[012356789]|17[013678]|18[0-9]|14[57])[0-9]{8}$/;
    if (!myreg.test($poneInput.val())) {
        return false;
    } else {
        return true;
    }
}

$('.layui-btn').submit(function () {
    var name_values = $("input [name='name']").val();
    var passwd_one = $("input [name='passwd_one']").val();
    var passwd_two = $("input [name='passwd_two']").val();
    if (name_values.length() < 11) {
        return console.log("手机号码少于11位")
    }else if (name_values.length() >11){
        return console.log("手机号码大于11位")
    }else if (name_values == 11) {
        var rel = isTelAvailable(name_values)
        if (!rel) {
            console.log("不是一个有效手机号码")
        }
        if (passwd_one != passwd_two) {
            console.log("两次密码不一致")
        }
    }
})

