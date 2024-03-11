/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-06 22:49:35
 * @LastEditTime: 2024-03-11 21:32:28
 * @LastEditors: yujiajie
 */
package constant

type Code int

type ErrorCode interface {
	String() string
}

const (
	SUCCESS Code = 100 + iota
	PARAMS_ERROR
	NOT_FOUND
	FORBID
	SERVER_ERROR
)

func (code Code) String() string {
	switch code {
	case SUCCESS:
		return "success"
	case PARAMS_ERROR:
		return "params error"
	case NOT_FOUND:
		return "not found"
	case FORBID:
		return "forbid"
	case SERVER_ERROR:
		return "server error"
	}
	return "unknown error"
}
