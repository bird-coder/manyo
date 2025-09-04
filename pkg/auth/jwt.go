/*
 * @Author: yujiajie
 * @Date: 2024-05-13 17:28:55
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-07-04 10:18:06
 * @FilePath: /Go-Base/pkg/auth/jwt.go
 * @Description:
 */
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	Data map[string]any
	jwt.RegisteredClaims
}

type JwtAuth struct {
	opts *jwtOptions
}

func NewJwtAuth(opts ...jwtOptionFunc) *JwtAuth {
	o := defaultJwtOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &JwtAuth{
		opts: o,
	}
}

// /////////////////对称加密//////////////////////
func (a *JwtAuth) Issue(data map[string]any, secretKey []byte) (string, *MyCustomClaims, error) {
	now := time.Now()
	claims := &MyCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(a.opts.expire)),
			Issuer:    a.opts.issuer,
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//token加密
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", claims, err
	}
	return signedToken, claims, nil
}

func (a *JwtAuth) Auth(signedToken string, secretKey []byte) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithLeeway(a.opts.leeway), jwt.WithIssuedAt(), jwt.WithIssuer(a.opts.issuer))
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// /////////////////非对称加密//////////////////////
func (a *JwtAuth) IssueRSA(data map[string]any, privateKeyPem string) (string, *MyCustomClaims, error) {
	if len(privateKeyPem) == 0 {
		return "", nil, errNoPrivateKey
	}
	now := time.Now()
	claims := &MyCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.opts.expire)),
			Issuer:    a.opts.issuer,
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	privateKeyData, err := os.ReadFile(privateKeyPem)
	if err != nil {
		return "", claims, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", claims, err
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//token加密
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", claims, err
	}
	return signedToken, claims, nil
}

func (a *JwtAuth) AuthRSA(signedToken string, publicKeyPem string) (*MyCustomClaims, error) {
	if len(publicKeyPem) == 0 {
		return nil, errNoPublicKey
	}
	token, err := jwt.ParseWithClaims(signedToken, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		publicKeyData, err := os.ReadFile(publicKeyPem)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	}, jwt.WithLeeway(a.opts.leeway), jwt.WithIssuedAt(), jwt.WithIssuer(a.opts.issuer))
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
