function submitForm() {
    var nameInput = document.getElementById("name");
    var telphone = document.getElementById("telphone");
    var wxid = document.getElementById("wxid");
    var comment = document.getElementById("comment");

    ajax({
        url: "./apply",
        type: "POST",                      
        data: { 
            name: nameInput.value, 
            telphone: telphone.value,
            wxid: wxid.value,
            comment: comment.value
        },      
        dataType: "json",
        success: function (response, xml) {
            alert("谢谢加入，稍微我们会与您联系！ \n你可先下载来一圈棋牌游戏体验！");
        },
        fail: function (status) {
            alert("谢谢加入，稍微我们会与您联系！ \n你可先下载来一圈棋牌游戏体验！");
        }
    });
}

function ajax(options) {
    options = options || {};
    options.type = (options.type || "GET").toUpperCase();
    options.dataType = options.dataType || "json";
    var params = formatParams(options.data);
    if (window.XMLHttpRequest) {
        var xhr = new XMLHttpRequest();
    } else { 
        var xhr = new ActiveXObject('Microsoft.XMLHTTP');
    }
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            var status = xhr.status;
            if (status >= 200 && status < 300) {
                options.success && options.success(xhr.responseText, xhr.responseXML);
            } else {
                options.fail && options.fail(status);
            }
        }
    }
    if (options.type == "GET") {
        xhr.open("GET", options.url + "?" + params, true);
        xhr.send(null);
    } else if (options.type == "POST") {
        xhr.open("POST", options.url, true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send(params);
    }
}

function formatParams(data) {
    var arr = [];
    for (var name in data) {
        arr.push(encodeURIComponent(name) + "=" + encodeURIComponent(data[name]));
    }
    arr.push(("v=" + Math.random()).replace(".",""));
    return arr.join("&");
}