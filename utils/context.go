package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Context struct {
	w    http.ResponseWriter
	r    *http.Request
	Body []byte
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	r.ParseForm()

	ctx := &Context{
		w: w,
		r: r,
	}

	ctx.Body, _ = ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	return ctx
}

func (c *Context) GetParam(k string) string {
	return c.r.FormValue(k)
}

func (c *Context) Response(data interface{}, code ...int) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		return
	// 	}
	// }()

	status := http.StatusOK
	if code != nil && len(code) > 0 {
		status = code[0]
	}

	var message string
	var result interface{}

	switch d := data.(type) {
	case error:
		message = d.Error()
		if status == http.StatusOK {
			status = http.StatusBadRequest
		}
	case string:
		if status != http.StatusOK {
			message = d
		} else {
			result = d
		}
	default:
		result = d
	}

	mapData := map[string]interface{}{
		"code":    status,
		"message": message,
		"result":  result,
	}

	rtn, err := json.Marshal(mapData)
	if err != nil {
		logger.Error("转换待返回数据失败: " + err.Error())
		// c.w.Write([]byte("内部错误"))
		// c.w.WriteHeader(http.StatusInternalServerError)
		http.Error(c.w, "转换待返回数据失败", http.StatusInternalServerError)
	} else {
		c.w.Header().Set("Content-Type", "application/json")
		c.w.WriteHeader(http.StatusOK)
		c.w.Write(rtn)
	}

	// panic("")
}
