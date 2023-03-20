package auth

import (
	"io/ioutil"
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
	privateKeyData, err := ioutil.ReadFile(privateKeyPem)
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
		publicKeyData, err := ioutil.ReadFile(publicKeyPem)
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
