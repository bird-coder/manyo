/*
 * @Author: yujiajie
 * @Date: 2025-07-03 14:51:48
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-09-04 10:47:23
 * @FilePath: /manyo/pkg/token/manager.go
 * @Description:
 */
package token

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bird-coder/manyo/pkg/auth"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
)

var (
	jtiKey        = "nbgame:%s:jti:%s"
	validVal      = "valid"
	errInvalidJTI = errors.New("invalid jti")
)

type TokenManager struct {
	opts *tokenOptions

	history sync.Map

	jwtAuth *auth.JwtAuth

	redisCli cache.AdapterCache

	lock sync.RWMutex
}

func NewTokenManager(redisCli cache.AdapterCache, opts ...tokenOptionFunc) *TokenManager {
	o := defaultTokenOptions()
	for _, opt := range opts {
		opt(o)
	}
	tm := &TokenManager{
		opts:     o,
		jwtAuth:  auth.NewJwtAuth(auth.WithIssuer(o.issuer), auth.WithLeeway(o.leeway), auth.WithExpire(o.expire)),
		redisCli: redisCli,
	}

	return tm
}

func (tm *TokenManager) SetSecret(secret string) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	tm.opts.secret = secret
}

func (tm *TokenManager) GetSecret() string {
	tm.lock.RLock()
	defer tm.lock.RUnlock()

	return tm.opts.secret
}

func (tm *TokenManager) SetPrevSecret(prevSecret string) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	tm.opts.prevSecret = prevSecret
}

func (tm *TokenManager) GetPrevSecret() string {
	tm.lock.RLock()
	defer tm.lock.RUnlock()

	return tm.opts.prevSecret
}

func (tm *TokenManager) GenerateToken(data map[string]any) (string, error) {
	secret := tm.GetSecret()
	tk, claims, err := tm.jwtAuth.Issue(data, []byte(secret))
	if err != nil {
		return "", err
	}
	if claims.ID != "" {
		tm.setCacheJTI(claims.ID)
	}

	return tk, nil
}

func (tm *TokenManager) RefreshToken(ctx *gin.Context) (string, error) {
	claims, err := tm.ParseToken(ctx)
	if err != nil {
		return "", err
	}
	tk, err := tm.GenerateToken(claims.Data)
	if err != nil {
		return "", err
	}
	tm.delCacheJTI(claims.ID)
	return tk, nil
}

func (tm *TokenManager) ParseToken(ctx *gin.Context) (claims *auth.MyCustomClaims, err error) {
	secret := tm.GetSecret()
	prevSecret := tm.GetPrevSecret()
	if len(prevSecret) > 0 {
		count := tm.loadCount(secret)
		prevCount := tm.loadCount(prevSecret)

		var first, second string
		if count > prevCount {
			first = secret
			second = prevSecret
		} else {
			first = prevSecret
			second = secret
		}

		claims, err = tm.doParseToken(ctx, first)
		if err != nil {
			claims, err = tm.doParseToken(ctx, second)
			if err != nil {
				return nil, err
			}

			tm.incrementCount(second)
		} else {
			tm.incrementCount(first)
		}
	} else {
		claims, err = tm.doParseToken(ctx, secret)
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func (tm *TokenManager) doParseToken(ctx *gin.Context, secret string) (*auth.MyCustomClaims, error) {
	tokenString, err := request.AuthorizationHeaderExtractor.ExtractToken(ctx.Request)
	if err != nil {
		return nil, err
	}
	claims, err := tm.jwtAuth.Auth(tokenString, []byte(secret))
	if err != nil {
		return nil, err
	}
	if claims.ID == "" || !tm.validJTI(claims.ID) {
		return nil, errInvalidJTI
	}
	return claims, nil
}

func (tm *TokenManager) incrementCount(secret string) {
	if time.Since(tm.opts.resetTime) > tm.opts.resetDuration {
		tm.history.Range(func(key, value any) bool {
			tm.history.Delete(key)
			return true
		})
	}
	val, ok := tm.history.Load(secret)
	if ok {
		atomic.AddUint64(val.(*uint64), 1)
	} else {
		var count uint64 = 1
		tm.history.Store(secret, &count)
	}
}

func (tm *TokenManager) loadCount(secret string) uint64 {
	val, ok := tm.history.Load(secret)
	if ok {
		return *val.(*uint64)
	}
	return 0
}

func (tm *TokenManager) validJTI(id string) bool {
	res := tm.getCacheJTI(id)
	return res == validVal
}

func (tm *TokenManager) getJTIKey(id string) string {
	return fmt.Sprintf(jtiKey, tm.opts.appName, id)
}

func (tm *TokenManager) getCacheJTI(id string) string {
	if tm.redisCli == nil {
		return validVal
	}
	redisKey := tm.getJTIKey(id)
	res, err := tm.redisCli.Get(redisKey)
	if err != nil {
		return ""
	}
	return res
}

func (tm *TokenManager) setCacheJTI(id string) error {
	if tm.redisCli == nil {
		return nil
	}
	redisKey := tm.getJTIKey(id)
	err := tm.redisCli.Set(redisKey, validVal, int(tm.opts.resetDuration))
	return err
}

func (tm *TokenManager) delCacheJTI(id string) error {
	if tm.redisCli == nil {
		return nil
	}
	redisKey := tm.getJTIKey(id)
	err := tm.redisCli.Del(redisKey)
	return err
}
