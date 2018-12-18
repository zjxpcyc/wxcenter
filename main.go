package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/zjxpcyc/wxcenter/utils"
)

const Version = "0.1.0"

var ver = flag.Bool("v", false, "version")
var port = flag.Int("p", 9000, "http listen port default is 9000")

// ServeHTTP 路由分配
func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			if msg, ok := p.(string); ok {
				if msg == "" {
					return
				}
			}

			logger.Error(string(debug.Stack()))

			http.Error(w, "微信服务中心: 严重错误", http.StatusInternalServerError)
			return
		}
	}()

	path := r.URL.Path

	switch true {
	case strings.Index(path, "/wx/register") > -1:
		// 服务注册
		t.Registe(w, r)

	case strings.Index(path, "/wx/signature") > -1:
		// 公众号接入
		t.Signature(w, r)

	case strings.Index(path, "/wx/notifiy/text") > -1:
		// 发送被动文本消息
		t.NotifiyText(w, r)

	case strings.Index(path, "/wx/notify/tpl") > -1:
		// 发送模板消息
		t.NotifiyTpl(w, r)

	case strings.Index(path, "/wx/jssdk") > -1:
		// jssdk 校验
		t.Jssdk(w, r)

	case strings.Index(path, "/wx/qrcode-temp") > -1:
		// 获取临时二维码
		t.QrcodeTmp(w, r)

	case strings.Index(path, "/wx/userinfo") > -1:
		// 依据 code 获取用户信息
		t.UserInfo(w, r)

	case strings.Index(path, "/wx/userdetail") > -1:
		// 依据 openid 获取用户信息
		t.UserDetail(w, r)

	default:
		http.NotFound(w, r)
	}
}

// NewHandler http handler
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", new(Router))
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
	serv := http.Server{Addr: addr, Handler: NewHandler()}

	logger.Info("开始进行系统初始化完成 .")
	logger.Info("启动成功 http://" + addr)

	serv.ListenAndServe()
}
