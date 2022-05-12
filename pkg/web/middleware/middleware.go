package middleware

import "net/http"

type HttpHandlerFn func(http.ResponseWriter, *http.Request) error
type MiddlewareFn func(HttpHandlerFn) HttpHandlerFn

func Wrap(fn HttpHandlerFn, middlewares ...MiddlewareFn) func(http.ResponseWriter, *http.Request) error {
	h := fn
	for i := len(middlewares) - 1; i >= 0; i-- {
		midFn := middlewares[i]
		if midFn != nil {
			h = midFn(h)
		}
	}
	return h
}
