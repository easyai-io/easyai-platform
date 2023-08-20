package contextx

import (
	"context"
)

type (
	transCtx     struct{}
	noTransCtx   struct{}
	transLockCtx struct{}
	userIDCtx    struct{}
	userNameCtx  struct{}
	userPhoneCtx struct{}
	userEmailCtx struct{}
	traceIDCtx   struct{}
)

// NewTrans Wrap transaction context
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// FromTrans  ...
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// NewNoTrans  ...
func NewNoTrans(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTransCtx{}, true)
}

// FromNoTrans  ...
func FromNoTrans(ctx context.Context) bool {
	v := ctx.Value(noTransCtx{})
	return v != nil && v.(bool)
}

// NewTransLock  ...
func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, true)
}

// FromTransLock  ...
func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil && v.(bool)
}

// NewUserUID  ...
func NewUserUID(ctx context.Context, userUID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userUID)
}

// FromUserUID  ...
func FromUserUID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserName  ...
func NewUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userNameCtx{}, userName)
}

// FromUserName  ...
func FromUserName(ctx context.Context) string {
	v := ctx.Value(userNameCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserPhone  ...
func NewUserPhone(ctx context.Context, phone string) context.Context {
	return context.WithValue(ctx, userPhoneCtx{}, phone)
}

// FromUserPhone  ...
func FromUserPhone(ctx context.Context) string {
	v := ctx.Value(userPhoneCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserEmail  ...
func NewUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailCtx{}, email)
}

// FromUserEmail  ...
func FromUserEmail(ctx context.Context) string {
	v := ctx.Value(userEmailCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewTraceID  ...
func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

// FromTraceID  ...
func FromTraceID(ctx context.Context) (string, bool) {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}
