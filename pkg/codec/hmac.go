/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-06 22:33:50
 * @LastEditTime: 2024-03-06 22:34:02
 * @LastEditors: yujiajie
 */
package codec

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Md5(body string) string {
	h := md5.New()
	io.WriteString(h, body)
	return hex.EncodeToString(h.Sum(nil))
}

func Hmac(key []byte, body string) []byte {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, body)
	return h.Sum(nil)
}

func HmacBase64(key []byte, body string) string {
	return base64.StdEncoding.EncodeToString(Hmac(key, body))
}
