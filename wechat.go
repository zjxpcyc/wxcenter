package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zjxpcyc/wechat-scheduler/jobs"

	"github.com/zjxpcyc/wxcenter/utils"

	"github.com/zjxpcyc/tinylogger"
	"github.com/zjxpcyc/wechat.v2/mini"
	"github.com/zjxpcyc/wechat.v2/wx"
)

var logger tinylogger.LogService
var wxClients map[string]*wx.Client
var miniClients map[string]*mini.Client

// var payClients map[string]*pay.Client

// Scheduler 只有一个 String() 方法的实现
// 实现 wechat.v2 Scheduler 接口
type Scheduler struct {
	url string
}

// Result 实现 wechat.v2 Scheduler 接口
func (t *Scheduler) Result() string {
	config := utils.GetConf("callback")
	res, err := utils.Request(config.Get("method").String(), t.url, nil, nil)
	if err != nil {
		logger.Error("获取 " + t.url + " 结果出错: " + err.Error())
		return ""
	}

	status := int(res["code"].(float64))
	if status != 200 {
		logger.Error("获取 " + t.url + " 结果出错: " + res["message"].(string))
		return ""
	}

	return res["result"].(string)
}

// NewScheduler new Scheduler
func NewScheduler(appid string, typ int) *Scheduler {
	config := utils.GetConf("callback")
	addr := config.Get("url").String()

	addr = strings.Replace(addr, "{appid}", appid, -1)
	addr = strings.Replace(addr, "{type}", strconv.Itoa(typ), -1)

	return &Scheduler{
		url: addr,
	}
}

// RegisterWxClient 公众号注册
func RegisterWxClient(appid string, cert map[string]string) {
	if wxClients == nil {
		wxClients = map[string]*wx.Client{}
	}

	// 初始化并入库
	m, err := NewModel("wx-" + appid)
	if err != nil {
		logger.Error("加载数据库失败: " + err.Error())
	} else {
		cf, _ := json.Marshal(cert)
		m.Update("cert", string(cf))
	}

	cli, ok := wxClients[appid]
	if !ok {
		cli = wx.NewClient(cert)
		wxClients[appid] = cli
	}

	// 注册中控中心
	registeWxScheduler(appid, cert)
	// 启动定时任务
	cli.SetAccessToken(NewScheduler(appid, jobs.JOB_ACCESS_TOKEN))
	cli.SetJSAPITicket(NewScheduler(appid, jobs.JOB_JSAPI_TICKET))

	logger.Info("注册公众号服务成功 appid=" + appid)
}

// RegisterWxMiniClient 小程序注册
func RegisterWxMiniClient(appid string, cert map[string]string) {
	// if miniClients == nil {
	// 	miniClients = map[string]*mini.Client{
	// 		appid: mini.NewClient(cert),
	// 	}

	// 	NewModel("mini-" + appid)
	// } else {
	// 	if _, ok := miniClients[appid]; !ok {
	// 		miniClients[appid] = mini.NewClient(cert)
	// 	}
	// }

	// logger.Info("注册微信小程序服务 appid=" + appid)
}

// GetWeChatClient returns the wx client of appid
func GetWeChatClient(appid string) (*wx.Client, error) {
	cli, ok := wxClients[appid]
	if !ok {
		return nil, errors.New("请先进行公众号客户端注册")
	}

	return cli, nil
}

// GetMiniClient returns the mini-programe client of appid
func GetMiniClient(appid string) (*mini.Client, error) {
	cli, ok := miniClients[appid]
	if !ok {
		return nil, errors.New("请先进行小程序客户端注册")
	}

	return cli, nil
}

// registeWxScheduler 注册定时任务
// 注册 access-token 与 jsapi_ticket
func registeWxScheduler(appid string, cert map[string]string) {
	params := map[string]interface{}{
		"appid":     appid,
		"appsecret": cert["secret"],
		"tasks": []map[string]int{
			map[string]int{
				"type": 0,
			},
			map[string]int{
				"type": 2,
			},
		},
	}

	paramsBytes, _ := json.Marshal(params)
	paramsStr := string(paramsBytes)

	logger.Info("开始注册定时任务 ...")
	logger.Info("注册内容: " + paramsStr)

	config := utils.GetConf("registe")
	res, err := utils.Request(config.Get("method").String(), config.Get("url").String(), nil, strings.NewReader(paramsStr))
	if err != nil {
		logger.Error("注册定时任务失败: " + err.Error())
	}

	var code int
	if c, ok := res["code"]; ok {
		code = int(c.(float64))
	}
	if code != http.StatusOK {
		logger.Error("注册定时任务失败: " + res["message"].(string))
		logger.Error(fmt.Sprintf("错误详情: %v", res))
	}
}

// WechatInit 微信初始化
func WechatInit() {
	logger.Info("开始微信客户端初始化...")

	for schema, m := range AllModel {
		cert := map[string]string{}
		res, err := m.Query("cert")
		if err != nil {
			logger.Error("初始化数据库内容失败: " + err.Error())
			continue
		}

		err = json.Unmarshal([]byte(res), &cert)
		if err != nil {
			logger.Error("初始化数据库内容失败: " + err.Error())
			continue
		}

		wxInfo := strings.Split(schema, "-")
		wxTyp := wxInfo[0]
		appid := wxInfo[1]

		if wxTyp == "wx" {
			RegisterWxClient(appid, cert)
		} else if wxTyp == "mini" {
			RegisterWxMiniClient(appid, cert)
		}
	}
}

func init() {
	logger = tinylogger.NewLogger()

	wx.SetLogger(logger)
	mini.SetLogger(logger)
	utils.SetLogger(logger)
}
