package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/smartwalle/alipay/v3"
)

var client *alipay.Client

const (
	// 应用id
	AppId = "9021000122695171"
	// 应用私钥
	PrivateKey = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCBUVcfoTu+lDxWIQLOXkfOqZdOWDLUGNCjMeEVILV/mcr9CjSvakjcMxlwuGyK+vMMijk7uEXBb9VeC+clflSGaHjt0MssNDiYqHJ/4vl8W1PDLyYdOItJshvrZMb8FTrtELLHMFCAj6QgeMW1PxDGHx3QSzbECmJKT3kZ6YbOVwWYiWUQiV72SiJiWCp6RpRNNd7BeaFLY3R078vaxRz//uZN9Oipmq96XhHA2iLRjV8/3WMDDUJb1b7Yz70xePbV3VugmD6d7t1K3/GXxYkmqW6+OVTFifiDb93abiVhwRu/WjB7NPMuejTE1cu7/T1u++LOQs+fGKZA11lozhspAgMBAAECggEAPDdZH2RfwpWaJu0GNnxWaZg02YleWb8YX/WY/tKVUU6W5A89izUtqkOtI/zspibGyF8Q2YssCDAAJePbBW13BINfVDX2daX3eLZvYreZPtvz/B5XXCH5Uh66u3jY44crQjqVVQVdZw/7+Bbk0UoFkIvqwRRU8yJ2tG2fUX6ZhDj97EhYSe6tuOowD3n5vx447gKIOEy3nMLatqtxtk9gySPZmLTvuhAsIwcWgljGkwq/Is70CS7MDsokg/YPnthTQsr9cRG5iywceHnwmPQD2PztmwuRls3z/w7u17x0majhSCO+xnLy1nXsUQoG8ogWwITHGldxwy1T4qE0ih+qAQKBgQC46QEAF2+7QBCEKIcXigZJQ0TmaJMfaEqm4HDvyfDhGAvrDKkfP3LYCB3Oa+vCNhj0RUQJMwzoPwElKuZOtsj/X5U4D/VuYAF4RCwcIBwL4Lm06N8tedOAHAN+vy4ZoVLQRPxV8JX9jAaTsXE9n0OK4HHmjHRmDZJjmB36W8ob6QKBgQCzCOM/lP0ndffGXEtSVOQ1j8e6n8cTTWdga8iV1kh3zCkEBmkrQuFZXq2bKJPN9NNz5l96YGYMCI/ASr5O9H4i8KWJmNjh3hjRYQMzCWDRsp1im3OQ7q9oncIJ1mA7FIQMIl5D8/6bF1e/LvBJdC+1UerWPPhDWBGsO2jJkZK9QQKBgQC06jwxU7zc7zR5qZFrRX8bBTcPW/e+PfL0TRoScnk8MqPOiKebzB9YILDQ+yRC81z8+hw0B/+z55j+PXfyQcJsoZ9Ep9CQ+lvVyJWDuyLVuDzaNRHO06hMapw80V6QcxescCKXDvohhXQV4wGRshaKdUjbskZcZyD4UqfaAR7AqQKBgHjOYHkA0amU4nJIyNJvUeYKdN0q/yu5KS5YzGq+wvuDGZILuV9lq6WgS0jNIp7wutYT9w0eiv1HsagxRyUDuTFebHTiXEZclSaDbaM8isY03hoxhtOfG2FeQhZdP2XePBPsBOuZco24PI9W3vDRo3eYJPwW+/aFMLelBtosjnWBAoGAYS3VnpoHJ/vuhaXlZR9lmoqlJ2yHAo/qt6ruEL+XOOKsiFLtOZJg0GkjxoMiGcoA8DFGm06FvcaJfxTmdUPo6+77jiObnxTS02dIJp2y9oijrJOB2STOe1tuz/YvJZTsyiT14CPdarju8fBl+r2++5PdfFfvv+7C6nmwnaDir64="
)

const (
	// 回调地址
	ServerDomain = "https://3409-183-14-29-68.ngrok-free.app"
	// 支付完成后支付宝同步通知的地址
	ReturnUrl = ServerDomain + "/callback"
	// 支付成功支付宝异步通知的地址
	NotifyUrl = ServerDomain + "/notify"
)

func init() {
	var err error

	if client, err = alipay.New(AppId, PrivateKey, false); err != nil {
		fmt.Println("初始化支付宝失败", err)
		return
	}

	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile("./cert/appPublicCert.crt"); err != nil {
		fmt.Println("加载应用公钥证书失败", err)
		return
	}

	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile("./cert/alipayRootCert.crt"); err != nil {
		fmt.Println("加载支付宝根证书失败", err)
		return
	}

	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile("./cert/alipayPublicCert.crt"); err != nil {
		fmt.Println("加载支付宝公钥证书失败", err)
		return
	}
}

func main() {
	http.HandleFunc("/pay", pay)   // 手机网页支付
	http.HandleFunc("/page", page) // 电脑网页支付 扫码支付

	http.HandleFunc("/close", close) // 关闭 只能关闭未付款的交易

	http.HandleFunc("/refund", refund) // 退款
	http.HandleFunc("/cancel", cancel) // 撤销 撤销交易(包含退款动作)

	http.HandleFunc("/notify", notify)     // 支付宝异步通知 通知函数
	http.HandleFunc("/callback", callback) // 支付宝同步通知 回调函数

	http.ListenAndServe(":8080", nil)
}

func pay(writer http.ResponseWriter, request *http.Request) {
	oid := uuid.New().String()
	p := alipay.TradeWapPay{}
	p.NotifyURL = NotifyUrl // 异步通知
	p.ReturnURL = ReturnUrl // 同步通知
	p.Subject = "手机网站支付" + oid
	p.OutTradeNo = oid
	p.TotalAmount = "100.00"
	p.ProductCode = "QUICK_WAP_PAY"

	var url, err = client.TradeWapPay(p)
	if err != nil {
		fmt.Println("调用手机网站接口失败", err)
		return
	}

	http.Redirect(writer, request, url.String(), http.StatusTemporaryRedirect)
}

func page(writer http.ResponseWriter, request *http.Request) {
	oid := uuid.New().String()
	p := alipay.TradePagePay{}
	p.NotifyURL = NotifyUrl
	p.ReturnURL = ReturnUrl
	p.Subject = "电脑网站支付" + oid
	p.OutTradeNo = oid
	p.TotalAmount = "100.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println("调用电脑网站支付接口失败", err)
		return
	}

	http.Redirect(writer, request, url.String(), http.StatusTemporaryRedirect)
}

func close(writer http.ResponseWriter, request *http.Request) {
	p := alipay.TradeClose{}
	p.OutTradeNo = "123456abcdefddd"

	res, err := client.TradeClose(p)
	if err != nil {
		fmt.Println("调用关闭交易接口失败", err)
		return
	}
	if res.Code != "10000" {
		fmt.Printf("订单%s关闭失败:%s---%s\n", res.OutTradeNo, res.Msg, res.SubMsg)
		return
	}

	fmt.Printf("订单%s关闭成功\n交易号%s\n", res.OutTradeNo, res.TradeNo)
}

func refund(writer http.ResponseWriter, request *http.Request) {
	rid := uuid.New().String()
	p := alipay.TradeRefund{}
	p.OutRequestNo = rid
	p.OutTradeNo = "123456abcdefddd"
	p.RefundAmount = "50.00"

	var res, err = client.TradeRefund(p)
	if err != nil {
		fmt.Println("调用退款接口失败", err)
		return
	}
	if res.Code != "10000" {
		fmt.Println("订单", res.OutTradeNo, "退款失败")
		return
	}
	fmt.Printf("订单%s退款成功\n交易号%s\n", res.OutTradeNo, res.TradeNo)
}

func cancel(writer http.ResponseWriter, request *http.Request) {
	p := alipay.TradeCancel{}
	p.OutTradeNo = "123456abcdefddd"

	res, err := client.TradeCancel(p)
	if err != nil {
		fmt.Println("调用撤销交易接口失败", err)
		return
	}
	if res.Code != "10000" {
		fmt.Printf("订单%s撤销失败:%s---%s\n", res.OutTradeNo, res.Msg, res.SubMsg)
		return
	}

	fmt.Printf("订单%s撤销成功\n", res.OutTradeNo)
}

func notify(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("\n*****支付宝平台调用notify函数*****")

	// 解析表单
	err := request.ParseForm()
	if err != nil {
		fmt.Println("解析表单数据错误", err)
		return
	}

	// 解析异步通知
	notification, err := client.DecodeNotification(request.Form)
	if err != nil {
		fmt.Println("解析异步通知错误", err)
		return
	}

	// 验证订单
	var rsp *alipay.TradeQueryRsp
	var p = alipay.NewPayload("alipay.trade.query")
	p.Set("out_trade_no", notification.OutTradeNo)
	if err = client.Request(p, &rsp); err != nil {
		fmt.Printf("订单%s验证错误: %s\n", notification.OutTradeNo, err.Error())
		return
	}
	if rsp.IsFailure() {
		fmt.Printf("订单%s验证错误: %s---%s\n", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)
		return
	}

	// 返回异步通知成功处理的消息给支付宝
	client.ACKNotification(writer)

	// 打印结果
	if notification.TradeStatus == alipay.TradeStatusSuccess || notification.TradeStatus == alipay.TradeStatusFinished {
		fmt.Printf("订单%s处理成功\n支付宝交易号:%s\n异步通知ID:%s\n", notification.OutTradeNo, notification.TradeNo, notification.NotifyId)
	} else {
		fmt.Printf("订单%s处理失败\n支付宝交易号:%s\n", notification.OutTradeNo, notification.TradeNo)
	}
}

func callback(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("\n*****支付宝平台回调callback函数*****")

	// 解析表单
	err := request.ParseForm()
	if err != nil {
		fmt.Println("解析表单数据错误", err)
		return
	}

	// 验签
	if err := client.VerifySign(request.Form); err != nil {
		fmt.Println("验证签名错误", err) // 状态码400 签名错误
		return
	}
	fmt.Println("签名验证通过")

	// 验证订单
	outTradeNo := request.Form.Get("out_trade_no")
	p := alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := client.TradeQuery(p)
	if err != nil {
		fmt.Printf("订单%s验证错误: %s\n", outTradeNo, err.Error())
		return
	}
	if rsp.IsFailure() {
		fmt.Printf("订单%s验证错误: %s---%s\n", outTradeNo, rsp.Msg, rsp.SubMsg)
		return
	}

	// 打印结果
	fmt.Printf("订单%s完成支付\n支付宝交易号:%s\n", outTradeNo, request.Form.Get("trade_no"))
}
