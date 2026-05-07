/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-06 22:46:08
 * @LastEditTime: 2026-05-07 14:45:00
 * @LastEditors: yujiajie 1037297660@qq.com
 */
package constant

const (
	AllowOrigin      = "Access-Control-Allow-Origin"
	AllOrigins       = "*"
	AllowMethods     = "Access-Control-Allow-Methods"
	AllowHeaders     = "Access-Control-Allow-Headers"
	AllowCredentials = "Access-Control-Allow-Credentials"
	ExposeHeaders    = "Access-Control-Expose-Headers"
	RequestMethod    = "Access-Control-Request-Method"
	RequestHeaders   = "Access-Control-Request-Headers"
	AllowHeadersVal  = "Content-Type, Origin, X-CSRF-Token, Authorization, AccessToken, Token, Range"
	ExposeHeadersVal = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers"
	Methods          = "GET, HEAD, POST, PATCH, PUT, DELETE, OPTIONS"
	AllowTrue        = "true"
	MaxAgeHeader     = "Access-Control-Max-Age"
	MaxAgeHeaderVal  = "86400"
	VaryHeader       = "Vary"
	OriginHeader     = "Origin"
)

const (
	ContentEncoding  = "Content-Encoding"
	ContentLength    = "Content-Length"
	ContentSecurity  = "X-Content-Security"
	ContentSignature = "X-Content-Signature"
	RequestUriHeader = "X-Request-Uri"
	ApplicationJson  = "application/json"
	ApplicationForm  = "application/x-www-form-urlencoded"
	ContentType      = "Content-Type"
	JsonContentType  = "application/json; charset=utf-8"
	KeyField         = "key"
	SecretField      = "secret"
	TypeField        = "type"
	SignatureField   = "signature"
	TimeField        = "time"
	CryptionType     = 1
)
