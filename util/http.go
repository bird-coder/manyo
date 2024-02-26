/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-23 19:13:13
 * @LastEditTime: 2023-09-23 22:28:48
 * @LastEditors: yuanshisan
 */
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpGet(url string, params map[string]string) ([]byte, error) {
	targetUrl := url
	if len(params) > 0 {
		targetUrl = url + "?" + buildParams(params)
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

func buildParams(params map[string]string) string {
	var res string
	for k, v := range params {
		res = res + k + "=" + v + "&"
	}
	return strings.Trim(res, "&")
}
