package httprpc

import (
	"context"
	"net/http"
)

type NextMiddleware func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type Middleware interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddleware) error
}

type MiddlewareFunc func(ctx context.Context,
	w http.ResponseWriter, r *http.Request, next NextMiddleware) error

func (f MiddlewareFunc) ServeHTTP(ctx context.Context,
	w http.ResponseWriter, r *http.Request, next NextMiddleware) error {
	return f(ctx, w, r, next)
}
