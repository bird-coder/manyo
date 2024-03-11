/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-11 22:50:01
 * @LastEditTime: 2024-03-11 22:54:00
 * @LastEditors: yujiajie
 */
package response

import "net/http"

var Default = &APIResponse{}

func Error(code int, msg string) Response {
	res := Default.Clone()
	if msg != "" {
		res.SetMsg(msg)
	}
	res.SetCode(code)
	return res
}

func OK(data interface{}, msg string) Response {
	res := Default.Clone()
	if msg != "" {
		msg = "success"
	}
	res.SetMsg(msg)
	res.SetCode(http.StatusOK)
	res.SetData(data)
	return res
}
