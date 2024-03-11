/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-11 22:47:01
 * @LastEditTime: 2024-03-11 22:48:57
 * @LastEditors: yujiajie
 */
package response

type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (res *APIResponse) SetData(data interface{}) {
	res.Data = data
}

func (res *APIResponse) SetMsg(msg string) {
	res.Msg = msg
}

func (res *APIResponse) SetCode(code int) {
	res.Code = code
}

func (res APIResponse) Clone() Response {
	return &res
}
