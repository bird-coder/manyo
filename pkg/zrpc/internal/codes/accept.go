/*
 * @Author: yujiajie
 * @Date: 2025-01-09 23:45:28
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-09 23:46:42
 * @FilePath: /Go-Base/pkg/zrpc/internal/codes/accept.go
 * @Description:
 */
package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss,
		codes.Unimplemented, codes.ResourceExhausted:
		return false
	default:
		return true
	}
}
