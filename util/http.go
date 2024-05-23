/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-23 19:13:13
 * @LastEditTime: 2024-05-23 17:32:41
 * @LastEditors: yujiajie
 */
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func HttpGet(url string, params map[string]interface{}) ([]byte, error) {
	targetUrl := url
	if len(params) > 0 {
		targetUrl = url + "?" + BuildParams(params)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get error, url: %s, params: %v, err: %v", url, params, err)
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	return res, nil
}

func HttpPost(url string, data map[string]interface{}) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("http post error, url: %s, data: %v, err: %v", url, data, err)
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	return res, nil
}

func BuildParams(params map[string]interface{}) string {
	var res string
	var val string
	for k, v := range params {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Float32, reflect.Float64: //浮点数需要先转成字符串，防止精度丢失
			val = fmt.Sprintf("%s=%s", k, strconv.FormatFloat(rv.Float(), 'f', -1, 64))
		default:
			val = fmt.Sprintf("%s=%v", k, v)
		}
		res = res + val + "&"
	}
	return strings.Trim(res, "&")
}
