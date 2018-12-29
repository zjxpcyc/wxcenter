package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zjxpcyc/gen"
	"github.com/zjxpcyc/wechat.v2/wx"
	"github.com/zjxpcyc/wxcenter/utils"
)

// Router 路由
type Router struct{}

// Registe 注册公众号或者小程序等
func (t *Router) Registe(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	secret := ctx.GetParam("secret")
	token := ctx.GetParam("token")
	aeskey := ctx.GetParam("aeskey")

	cert := map[string]string{
		"appid":  appid,
		"secret": secret,
		"token":  token,
		"aeskey": aeskey,
	}

	RegisterWxClient(appid, cert)
	ctx.Response("ok")
}

// Signature 公众号初始接入
func (t *Router) Signature(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	timestamp := ctx.GetParam("timestamp")
	nonce := ctx.GetParam("nonce")

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	ctx.Response(wxCli.Signature(timestamp, nonce))
}

// NotifiyText 被动返回文本
func (t *Router) NotifiyText(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	to := ctx.GetParam("to")
	message := ctx.GetParam("message")

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	res, err := wxCli.ResponseMessageText(to, message)
	if err != nil {
		logger.Error("处理被动消息失败: " + err.Error())
		ctx.Response("处理被动消息失败", http.StatusInternalServerError)
	}

	ctx.Response(res)
}

// NotifiyTpl 发送模板消息
func (t *Router) NotifiyTpl(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	to := ctx.GetParam("to")
	tplID := ctx.GetParam("tpl")
	link := ctx.GetParam("link")
	if link != "" {
		var err error
		link, err = url.QueryUnescape(link)
		if err != nil {
			logger.Error("发送模板消息失败: " + err.Error())
			ctx.Response(errors.New("发送消息链接不正确"))
		}
	}

	var dtRaw []byte
	dt := ctx.GetParam("data")
	if dt == "" {
		ctx.Response(errors.New("发送消息内容不能为空"))
	} else {
		var err error
		dtRaw, err = gen.Base64Decode(dt)
		if err != nil {
			logger.Error("处理模板消息失败: " + err.Error())
			ctx.Response("处理模板消息失败", http.StatusInternalServerError)
		}
	}

	data := make(map[string]wx.TplMessageData)
	if err := json.Unmarshal(dtRaw, &data); err != nil {
		logger.Error("处理模板消息失败: " + err.Error())
		ctx.Response("处理模板消息失败", http.StatusInternalServerError)
	}

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	err = wxCli.SendTplMessage(to, tplID, link, data)
	if err != nil {
		logger.Error("模板消息发送失败: " + err.Error())
		ctx.Response("模板消息发送失败", http.StatusInternalServerError)
	}

	if err := json.Unmarshal(dtRaw, &data); err != nil {
		logger.Error("发送模板消息失败: " + err.Error())
		ctx.Response("发送模板消息失败", http.StatusInternalServerError)
	}

	ctx.Response("ok")
}

// Jssdk 微信 JSAPI SDK 校验
func (t *Router) Jssdk(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	link, err := url.QueryUnescape(ctx.GetParam("url"))
	if err != nil {
		logger.Error("JSSDK校验失败: " + err.Error())
		ctx.Response(errors.New("JSSDK校验失败"))
	}

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	logger.Info("获取到分享链接 - ", link)

	res := wxCli.GetJSTicketSignature(link)
	ctx.Response(res)
}

// QrcodeTmp 临时二维码
func (t *Router) QrcodeTmp(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	qrStr, err := url.QueryUnescape(ctx.GetParam("qrstr"))
	if err != nil {
		logger.Error("获取临时二维码失败: " + err.Error())
		ctx.Response(errors.New("获取临时二维码失败"))
	}

	exp := ctx.GetParam("expire")

	var expire int64
	if exp != "" {
		expire, _ = strconv.ParseInt(exp, 10, 64)
	}

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	var qrURL string
	if expire != 0 {
		qrURL, err = wxCli.GetTempStrQRCode(qrStr, time.Duration(expire))
	} else {
		qrURL, err = wxCli.GetTempStrQRCode(qrStr)
	}

	ctx.Response(qrURL)
}

// UserInfo Oauth2 用户信息获取
func (t *Router) UserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	code := ctx.GetParam("code")

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}

	res, err := wxCli.GetUserInfo(code)
	if err != nil {
		logger.Error("获取微信个人信息失败: " + err.Error())
		ctx.Response("获取微信个人信息失败", http.StatusInternalServerError)
	}

	ctx.Response(res)
}

// UserDetail 依据 openid 获取用户信息
func (t *Router) UserDetail(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContext(w, r)

	appid := ctx.GetParam("appid")
	openid := ctx.GetParam("openid")

	wxCli, err := GetWeChatClient(appid)
	if err != nil {
		ctx.Response(err)
	}
	res, err := wxCli.GetUserDetail(openid)
	if err != nil {
		logger.Error("获取微信个人详情失败: " + err.Error())
		ctx.Response("获取微信个人详情失败", http.StatusInternalServerError)
	}

	ctx.Response(res)
}
