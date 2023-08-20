package auth

import (
	"context"
	"errors"
)

var (
	// ErrInvalidToken ..
	ErrInvalidToken = errors.New("invalid token")
	// ErrDestroyedToken ..
	ErrDestroyedToken = errors.New("destroyed token")
	// ErrUserForbidden ..
	ErrUserForbidden = errors.New("user disabled")
)

// TokenInfo 令牌信息
type TokenInfo interface {
	// GetAccessToken 获取访问令牌
	GetAccessToken() string
	// GetTokenType 获取令牌类型
	GetTokenType() string
	// GetExpiresAt 获取令牌到期时间戳
	GetExpiresAt() int64
	// EncodeToJSON JSON编码
	EncodeToJSON() ([]byte, error)
}

// Auther 认证接口
type Auther interface {
	// GenerateToken 生成令牌
	GenerateToken(ctx context.Context, userUID, realName, phone, email string) (TokenInfo, error)

	// DestroyToken 销毁令牌
	DestroyToken(ctx context.Context, accessToken string) error

	// ParseUserInfo 解析用户info
	ParseUserInfo(ctx context.Context, accessToken string, noCheckInRedis ...bool) (userUID, realName, phone, email string, err error)

	// Release 释放资源
	Release() error

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, userUID, realName string, active bool)
}
