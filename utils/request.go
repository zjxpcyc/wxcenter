package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request 请求远程数据
func Request(method, addr string, query url.Values, body io.Reader) (result map[string]interface{}, err error) {
	// 请求 Body
	var bodyData io.Reader
	if method != http.MethodGet && body != nil {
		bodyData = body

		b := &bytes.Buffer{}
		io.Copy(b, body)
	} else {
		bodyData = nil
	}

	// 请求参数
	if query != nil {
		params := query.Encode()

		if strings.Index(addr, "?") > -1 {
			addr += "&" + params
		} else {
			addr += "?" + params
		}
	}

	// 构造 http 请求
	var req *http.Request
	var res *http.Response
	client := new(http.Client)

	req, err = http.NewRequest(method, addr, bodyData)
	if err != nil {
		return
	}

	res, err = client.Do(req)
	if err != nil {
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}

	// 格式化结果
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}
