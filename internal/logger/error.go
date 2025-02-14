package logger

import (
	"context"
	"errors"
)

type errorWithLogCtx struct {
	next error
	ctx  requestContext
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := requestContext{}
	if x, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var errCtx *errorWithLogCtx
	if errors.As(err, &errCtx) {
		return context.WithValue(ctx, requestCtxKey, errCtx.ctx)
	}
	return ctx
}
