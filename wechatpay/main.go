package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var client *core.Client
var mchkey *rsa.PrivateKey
var handler *notify.Handler

const (
	// APPID
	Appid = "wx49xxxxxxxxxxx975"
	// 商户号
	MchID = "1xxxxxxxx2"
	// 商户证书序列号
	MchCertificateSerialNumber = "4CAE5AxxxxxxxxxxxxxxxxxxxxxxxxxxxxBA9"
	// 商户APIv3密钥
	MchAPIv3Key = "9f39xxxxxxxxxxxxxxxxxx7fa1"
)

var ctx = context.TODO()
var err error

// 加载私钥
func InitMchKey() {
	mchkey, err = utils.LoadPrivateKeyWithPath("./cert/apiclient_key.pem")
	if err != nil {
		fmt.Printf("load merchant private key error")
		return
	}
}

// 创建微信支付客户端
func InitClient() {
	client, err = core.NewClient(context.TODO(), []core.ClientOption{option.WithWechatPayAutoAuthCipher(MchID, MchCertificateSerialNumber, mchkey, MchAPIv3Key)}...)
	if err != nil {
		fmt.Printf("new wechat pay client err:%s", err)
		return
	}
}

func InitHandler() {
	// 创建证书下载器管理器mgr
	mgr := downloader.MgrInstance()
	// 向mgr注册商户平台的证书下载器
	err := mgr.RegisterDownloaderWithPrivateKey(ctx, mchkey, MchCertificateSerialNumber, MchID, MchAPIv3Key)
	if err != nil {
		fmt.Printf("注册证书下载器失败:%v\n", err)
		return
	}
	// 使用mgr获取商户平台证书访问器
	certificateVisitor := mgr.GetCertificateVisitor(MchID)
	// 创建通知处理器
	handler, err = notify.NewRSANotifyHandler(MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		fmt.Printf("创建通知处理器失败:%v\n", err)
		return
	}
}

func main() {
	InitMchKey()
	InitClient()
	InitHandler()

	http.HandleFunc("/jspay", jspay)
	http.HandleFunc("/h5pay", h5pay)
	http.HandleFunc("/napay", napay)
	http.HandleFunc("/close", close)
	http.HandleFunc("/callback", callback)
	http.ListenAndServe(":8080", nil)
}

func jspay(w http.ResponseWriter, r *http.Request) {
	svc := jsapi.JsapiApiService{Client: client}
	resp, result, err := svc.PrepayWithRequestPayment(context.TODO(),
		jsapi.PrepayRequest{
			Appid:       core.String(Appid),
			Mchid:       core.String(MchID),
			Description: core.String("支付测试"),
			OutTradeNo:  core.String("123456abcdef"),
			NotifyUrl:   core.String("https://vczs.top/callback"),
			Amount: &jsapi.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(1),
			},
			Payer: &jsapi.Payer{
				Openid: core.String("wechat123456789"),
			},
		},
	)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(resp)
	fmt.Println(result.Response)
}

func h5pay(w http.ResponseWriter, r *http.Request) {
	svc := h5.H5ApiService{Client: client}
	resp, result, err := svc.Prepay(context.TODO(),
		h5.PrepayRequest{
			Appid:       core.String(Appid),
			Mchid:       core.String(MchID),
			Description: core.String("支付测试"),
			OutTradeNo:  core.String("123456abcdef"),
			NotifyUrl:   core.String("https://vczs.top/callback"),
			Amount: &h5.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(1),
			},
			SceneInfo: &h5.SceneInfo{
				H5Info: &h5.H5Info{
					Type: core.String("iOS"),
				},
				PayerClientIp: core.String("123.45.67.89"),
			},
		},
	)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(resp)
	fmt.Println(result.Response)
}

func napay(w http.ResponseWriter, r *http.Request) {
	svc := native.NativeApiService{Client: client}
	resp, result, err := svc.Prepay(context.TODO(),
		native.PrepayRequest{
			Appid:       core.String(Appid),
			Mchid:       core.String(MchID),
			Description: core.String("支付测试"),
			OutTradeNo:  core.String("123456abcdef"),
			NotifyUrl:   core.String("https://vczs.top/callback"),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(1),
			},
			Attach: core.String("平台币"),
		},
	)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(resp)
	fmt.Println(result.Response)
}

func close(w http.ResponseWriter, r *http.Request) {
	svc := native.NativeApiService{Client: client}
	result, err := svc.CloseOrder(context.TODO(),
		native.CloseOrderRequest{
			OutTradeNo: core.String("123456abcdef"),
			Mchid:      core.String(MchID),
		},
	)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(result.Response)
}

func callback(w http.ResponseWriter, r *http.Request) {
	notifyData := make(map[string]any)
	notifyReq, err := handler.ParseNotifyRequest(ctx, r, &notifyData)
	if err != nil {
		fmt.Printf("解析微信支付通知失败:%v\n", err)
		return
	}

	if notifyData["trade_state"] != "SUCCESS" {
		fmt.Printf("交易不成功:%v\n", err)
		return
	}

	fmt.Println(*notifyReq.Resource)
	fmt.Println(notifyReq.Summary)
	fmt.Println(notifyData)
}
