package logger

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// requestCtxKey is a custom type to avoid key collisions in context.
type keyType int

const (
	requestCtxKey keyType = 0 // Unique key for storing request context in context.Context.
)

// LogRequestHandler wraps a slog.Handler to enrich log records with request-specific data.
type LogRequestHandler struct {
	next slog.Handler
}

// NewLogRequestHandler creates a new LogRequestHandler that wraps the provided slog.Handler.
func NewLogRequestHandler(next slog.Handler) *LogRequestHandler {
	return &LogRequestHandler{next: next}
}

// Enabled is for implementing the slog.Handler interface.
func (h *LogRequestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

// Handle processes the log record, enriching it with request-specific data if available.
func (h *LogRequestHandler) Handle(ctx context.Context, rec slog.Record) error {
	if rCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		if rCtx.CallerIP != "" {
			rec.Add("caller_ip", rCtx.CallerIP) // Add caller IP to the log record.
		}
		if rCtx.Method != "" {
			rec.Add("method", rCtx.Method)
		}
		if rCtx.Endpoint != "" {
			rec.Add("endpoint", rCtx.Endpoint) // Add request URL to the log record.
		}
		if rCtx.UID != uuid.Nil {
			rec.Add("uid", rCtx.UID.String()) // Add user ID to the log record.
		}
		if rCtx.ResponseStatus != 0 {
			rec.Add("status_code", rCtx.ResponseStatus)
			rec.Add("text_status", http.StatusText(rCtx.ResponseStatus))
		}
		if rCtx.Duration != 0 {
			rec.Add("duration", rCtx.Duration)
		}
	}
	return h.next.Handle(ctx, rec)
}

// WithAttrs is for implementing the slog.Handler interfaces for.
func (h *LogRequestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogRequestHandler{next: h.next.WithAttrs(attrs)}
}

// WithGroup is for implementing the slog.Handler interface.
func (h *LogRequestHandler) WithGroup(name string) slog.Handler {
	return &LogRequestHandler{next: h.next.WithGroup(name)}
}

// requestContext holds request-specific data to be added to log records.
type requestContext struct {
	CallerIP       string
	Method         string
	Endpoint       string
	UID            uuid.UUID
	ResponseStatus int
	Duration       time.Duration
}

// WithCallerIP adds or updates the caller IP in the request context.
func WithCallerIP(ctx context.Context, ip string) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.CallerIP = ip
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{CallerIP: ip})
}

// WithEndpoint adds or updates the endpoint path in the request context.
func WithEndpoint(ctx context.Context, path string) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.Endpoint = path
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{Endpoint: path})
}

// WithUID adds or updates the user ID in the request context.
func WithUID(ctx context.Context, uid uuid.UUID) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.UID = uid
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{UID: uid})
}

// WithStatusCode adds or updates the status code of the response.
func WithStatusCode(ctx context.Context, code int) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.ResponseStatus = code
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{ResponseStatus: code})
}

// WithMethod adds or updates the method of the request.
func WithMethod(ctx context.Context, method string) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.Method = method
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{Method: method})
}

// WithDuration adds or updates the processing duration of request.
func WithDuration(ctx context.Context, d time.Duration) context.Context {
	if reqCtx, ok := ctx.Value(requestCtxKey).(requestContext); ok {
		reqCtx.Duration = d
		return context.WithValue(ctx, requestCtxKey, reqCtx)
	}
	return context.WithValue(ctx, requestCtxKey, requestContext{Duration: d})
}
