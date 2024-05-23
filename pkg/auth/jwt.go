/*
 * @Author: yujiajie
 * @Date: 2024-03-07 09:10:05
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:33:37
 * @FilePath: /manyo/pkg/auth/jwt.go
 * @Description:
 */
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Data any
	jwt.RegisteredClaims
}

// /////////////////对称加密//////////////////////
func Issue(data any, secretKey []byte) (string, error) {
	claims := MyCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			Issuer:    "aireal",
		},
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//token加密
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Auth(signedToken string, secretKey []byte) (any, error) {
	token, err := jwt.ParseWithClaims(signedToken, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithLeeway(5*time.Second))
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.Data, nil
	}
	return nil, err
}

// /////////////////非对称加密//////////////////////
func IssueRSA(data any, privateKeyPem string) (string, error) {
	claims := MyCustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			Issuer:    "aireal",
		},
	}
	privateKeyData, err := os.ReadFile(privateKeyPem)
	if err != nil {
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//token加密
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func AuthRSA(signedToken string, publicKeyPem string) (any, error) {
	token, err := jwt.ParseWithClaims(signedToken, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		publicKeyData, err := os.ReadFile(publicKeyPem)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	})
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.Data, nil
	}
	return nil, err
}
