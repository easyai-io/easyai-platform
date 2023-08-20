package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	lc "github.com/easyai-io/easyai-platform/pkg/localcache"
)

// tokenInfo 令牌信息
type tokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 令牌到期时间
}

func (t *tokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *tokenInfo) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}

const defaultKey = "easyai-platform"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       3600 * 24 * 365 * 2,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

// Option 定义参数项
type Option func(*options)

// SetSigningMethod 设定签名方式
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningKey 设定签名key
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc 设定验证key的回调函数
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// SetExpired 设定令牌过期时长(单位秒，默认7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// CheckUserStatusFn 检查用户状态函数
type CheckUserStatusFn func(userUID, realName string) (bool, error)

// NoCheckUserStatus 不检查用户状态
var NoCheckUserStatus = func(userUID, realName string) (bool, error) { return true, nil }

// New 创建认证实例
func New(cli *redis.Client, checkUserStatus CheckUserStatusFn, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts:            &o,
		redis:           cli,
		lcCache:         lc.New("jwt_easyai_user_status"),
		checkUserStatus: checkUserStatus,
	}
}

// JWTAuth jwt认证
type JWTAuth struct {
	opts            *options
	redis           *redis.Client
	lcCache         *lc.LocalCache
	checkUserStatus CheckUserStatusFn
}

// UserClaims user
type UserClaims struct {
	UserUID  string `json:"user_id"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(ctx context.Context, userUID, realName, phone, email string) (TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second)

	token := jwt.NewWithClaims(a.opts.signingMethod, &UserClaims{
		UserUID:  userUID,
		RealName: realName,
		Phone:    phone,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userUID,
		},
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt.Unix(),
		TokenType:   a.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

// 解析令牌
func (a *JWTAuth) parseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, a.opts.keyfunc)
	if err != nil || !token.Valid {
		fmt.Println("parse token error:", err)
		return nil, ErrInvalidToken
	}

	return token.Claims.(*UserClaims), nil
}

// DestroyToken 销毁令牌
func (a *JWTAuth) DestroyToken(ctx context.Context, tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// 如果设定了redis存储，则将失效的token标记在redis中
	key := "destroyed-token:" + tokenString
	if a.redis != nil {
		cmd := a.redis.Set(ctx, key, claims.Subject+"-deleted", time.Hour*24*180)
		return cmd.Err()
	}

	return nil
}

// ParseUserInfo 解析用户ID
func (a *JWTAuth) ParseUserInfo(ctx context.Context, tokenString string, noCheckInRedis ...bool) (userUID, realName, phone, email string, err error) {
	if tokenString == "" {
		return "", "", "", "", ErrInvalidToken
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", "", "", "", err
	}
	if len(noCheckInRedis) == 0 || !noCheckInRedis[0] {
		key := "destroyed-token:" + tokenString
		cmd := a.redis.Exists(ctx, key)
		if err := cmd.Err(); err != nil {
			return "", "", "", "", err
		}
		if cmd.Val() > 0 {
			return "", "", "", "", ErrDestroyedToken
		}
	}
	if !a.isUserActive(ctx, claims.UserUID, claims.RealName) {
		return "", "", "", "", ErrUserForbidden
	}
	return claims.UserUID, claims.RealName, claims.Phone, claims.Email, nil
}

// Release 释放资源
func (a *JWTAuth) Release() error {
	return a.redis.Close()
}

// UpdateUserStatus 更新用户状态
func (a *JWTAuth) UpdateUserStatus(ctx context.Context, userUID, realName string, active bool) {
	a.lcCache.SetKV(userUID, active, time.Minute*5)
	a.lcCache.SetKV(realName, active, time.Minute*5)
	if a.redis != nil {
		key := "user-status:" + userUID
		a.redis.Set(ctx, key, active, time.Hour*24*180)
	}
}

func (a *JWTAuth) isUserActive(ctx context.Context, userUID, realName string) (res bool) {

	defer func() {
		// 未命中或者用户被禁用, double check and update status
		if !res && a.checkUserStatus != nil {
			if status, err := a.checkUserStatus(userUID, realName); err == nil {
				res = status
				a.UpdateUserStatus(ctx, userUID, realName, status)
			}
		}
	}()

	// 本地cache命中
	if v, ok := a.lcCache.GetValue(userUID); ok {
		if active, ok := v.(bool); ok {
			return active
		}
	}

	// redis命中
	if a.redis != nil {
		key := "user-status:" + userUID
		cmd := a.redis.Get(ctx, key)
		if err := cmd.Err(); err == nil {
			if val, err := cmd.Bool(); err == nil {
				return val
			}
		}
	}
	return false
}
