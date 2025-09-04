/*
 * @Author: yujiajie
 * @Date: 2025-07-03 15:16:43
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-07-03 17:12:55
 * @FilePath: /Go-Base/pkg/token/options.go
 * @Description:
 */
package token

import (
	"time"
)

type tokenOptions struct {
	appName    string
	issuer     string
	leeway     time.Duration
	expire     time.Duration
	secret     string
	prevSecret string

	resetTime     time.Time
	resetDuration time.Duration
}

type tokenOptionFunc func(o *tokenOptions)

func defaultTokenOptions() *tokenOptions {
	return &tokenOptions{
		appName:       "rock",
		issuer:        "nbgame",
		leeway:        5 * time.Second,
		expire:        2 * time.Hour,
		resetTime:     time.Now(),
		resetDuration: 24 * time.Hour,
	}
}

func WithAppName(appName string) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.appName = appName
	}
}

func WithIssuer(issuer string) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.issuer = issuer
	}
}

func WithLeeway(leeway time.Duration) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.leeway = leeway
	}
}

func WithExpire(expire time.Duration) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.expire = expire
	}
}

func WithSecret(secret string) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.secret = secret
	}
}

func WithPrevSecret(prevSecret string) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.prevSecret = prevSecret
	}
}

func WithResetDuration(duration time.Duration) tokenOptionFunc {
	return func(o *tokenOptions) {
		o.resetDuration = duration
	}
}
