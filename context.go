package log

import (
	"context"
)

type contextKey string

const (
	ContextKeyLogger contextKey = "github.com/go-haru/log.Logger"
)

type contextStore interface {
	context.Context
	Set(string, any)
	Get(string) (any, bool)
}

func typeAssert[T any](val any) (dst T, ok bool) {
	dst, ok = val.(T)
	return dst, ok
}

func typeAssertWithOk[T any](val any, ok bool) (dst T, _ bool) {
	if !ok {
		return dst, false
	}
	dst, ok = val.(T)
	return dst, ok
}

func contextGetValue[T any](ctx context.Context, key contextKey) (value T, ok bool) {
	var store contextStore
	if store, ok = ctx.(contextStore); ok {
		return typeAssertWithOk[T](store.Get(string(key)))
	}
	return typeAssert[T](ctx.Value(key))
}

func contextSetValue(ctx context.Context, key contextKey, value any) context.Context {
	if store, ok := ctx.(contextStore); ok {
		store.Set(string(key), value)
		return store
	}
	return context.WithValue(ctx, key, value)
}

func Context(ctx context.Context, l Logger) context.Context {
	if l == nil {
		l = logger
	}
	return contextSetValue(ctx, ContextKeyLogger, l)
}

func C(ctx context.Context) (l Logger) {
	var ok bool
	if l, ok = contextGetValue[Logger](ctx, ContextKeyLogger); ok {
		return l
	}
	return logger
}
