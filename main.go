package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/zjxpcyc/wxcenter/utils"
)

const Version = "0.1.0"

var ver = flag.Bool("v", false, "version")
var port = flag.Int("p", 9000, "http listen port default is 9000")

// NewRouter 路由表
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	router := new(Router)

	// 服务注册
	mux.HandleFunc("/wx/register", router.Registe)

	// 公众号接入
	mux.HandleFunc("/wx/signature", router.Signature)

	// 发送被动文本消息
	mux.HandleFunc("/wx/notifiy-text", router.NotifiyText)

	// 发送模板消息
	mux.HandleFunc("/wx/notify/tpl", router.NotifiyTpl)

	// jssdk 校验
	mux.HandleFunc("/wx/jssdk", router.Jssdk)

	// 获取临时二维码
	mux.HandleFunc("/wx/qrcode-temp", router.QrcodeTmp)

	// 依据 code 获取用户信息
	mux.HandleFunc("/wx/userinfo", router.UserInfo)

	// 依据 openid 获取用户信息
	mux.HandleFunc("/wx/userdetail", router.UserDetail)

	return mux
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Println(Version)
		os.Exit(0)
	}

	// 部分初始化工作
	logger.Info("开始进行系统初始化 ...")
	utils.ConfInit()
	DBInit()
	WechatInit()

	addr := ":" + strconv.Itoa(*port)
	serv := http.Server{Addr: addr, Handler: NewRouter()}

	logger.Info("开始进行系统初始化完成 .")
	logger.Info("启动成功 http://" + addr)

	serv.ListenAndServe()
}
