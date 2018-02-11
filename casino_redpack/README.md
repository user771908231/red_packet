# casino_redpack
微信项目

## 依赖
mongodb redis
## 内网监听端口
127.0.0.1:9090
## 外网地址及端口
wx.tondeen.com:80  通过Nginx反向代理至内网地址
## 公众号支付回调地址
/mp/pay/callback

## 微信支付证书配置
apiclient_cert.pem
apiclient_key.pem
这两个证书文件存放于../conf目录下

## 微信参数配置
//公共
<br/>WX_APP_ID
<br/>WX_MCH_ID
<br/>//用于公众号支付 mch
<br/>WX_API_KEY
<br/>//用于Oauth2.0登录，收发消息等 mp
<br/>WX_APP_SECRET
<br/>
//存放于casino_redpack/model/weixinModel/init.go中
