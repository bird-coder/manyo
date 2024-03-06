/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-06 22:46:08
 * @LastEditTime: 2024-03-06 22:46:58
 * @LastEditors: yujiajie
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
	ContentSecurity  = "X-Content-Security"
	RequestUriHeader = "X-Request-Uri"
	ApplicationJson  = "application/json"
	ContentType      = "Content-Type"
	JsonContentType  = "application/json; charset=utf-8"
	KeyField         = "key"
	SecretField      = "secret"
	TypeField        = "type"
	SignatureField   = "signature"
	TimeField        = "time"
	CryptionType     = 1
)
