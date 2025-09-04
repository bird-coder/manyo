package auth

import (
	"errors"
	"time"
)

var (
	errNoPrivateKey = errors.New("缺少私钥文件")
	errNoPublicKey  = errors.New("缺少公钥文件")
)

func defaultJwtOptions() *jwtOptions {
	return &jwtOptions{
		issuer: "nbgame",
		leeway: 5 * time.Second,
		expire: 2 * time.Hour,
	}
}

type jwtOptions struct {
	issuer string
	leeway time.Duration
	expire time.Duration
}

type jwtOptionFunc func(o *jwtOptions)

func WithIssuer(issuer string) jwtOptionFunc {
	return func(o *jwtOptions) {
		o.issuer = issuer
	}
}

func WithLeeway(leeway time.Duration) jwtOptionFunc {
	return func(o *jwtOptions) {
		o.leeway = leeway
	}
}

func WithExpire(expire time.Duration) jwtOptionFunc {
	return func(o *jwtOptions) {
		o.expire = expire
	}
}
