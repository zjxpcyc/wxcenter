package utils

import (
	"io/ioutil"
	"os"

	"github.com/tidwall/gjson"
)

// Configer 配置文件读取器
var Configer gjson.Result

// ConfInit 配置初始化
func ConfInit() {
	f, err := os.Open("./conf/wxcenter.conf")
	if err != nil {
		Configer = gjson.Parse("{}")
		return
	}
	defer f.Close()

	dt, err := ioutil.ReadAll(f)
	if err != nil {
		Configer = gjson.Parse("{}")
		return
	}

	Configer = gjson.Parse(string(dt))
	return
}

// GetConf 配置读取
func GetConf(key string) gjson.Result {
	return Configer.Get(key)
}
