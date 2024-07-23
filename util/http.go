/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-23 19:13:13
 * @LastEditTime: 2024-07-23 15:14:15
 * @LastEditors: yujiajie
 */
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	defaultTimeout = 30 * time.Second
)

func HttpGet(targetUrl string, params map[string]interface{}) ([]byte, error) {
	if len(params) > 0 {
		targetUrl = targetUrl + "?" + BuildParams(params)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get error, url: %s, params: %v, err: %v", targetUrl, params, err)
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	return res, nil
}

func HttpPost(targetUrl string, data map[string]interface{}) ([]byte, error) {
	client := &http.Client{Timeout: defaultTimeout}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(targetUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("http post error, url: %s, data: %v, err: %v", targetUrl, data, err)
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	return res, nil
}

func HttpPostForm(targetUrl string, data map[string]interface{}) ([]byte, error) {
	client := &http.Client{Timeout: defaultTimeout}

	postData := url.Values{}
	for k, v := range data {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Float32, reflect.Float64:
			postData.Set(k, strconv.FormatFloat(rv.Float(), 'f', -1, 64))
		default:
			postData.Set(k, fmt.Sprintf("%v", v))
		}
	}
	resp, err := client.PostForm(targetUrl, postData)
	if err != nil {
		return nil, fmt.Errorf("http post error, url: %s, data: %v, err: %v", targetUrl, data, err)
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
