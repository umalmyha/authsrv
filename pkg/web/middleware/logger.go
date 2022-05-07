package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RequestLogger(logger *log.Logger) MiddlewareFn {
	return func(nextFn HttpHandlerFn) HttpHandlerFn {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			reqId, ok := ctx.Value(CtxReqId).(ctxReqIdKey)
			if !ok {
				reqId = ctxReqIdKey(uuid.NewString())
				logger.Printf("request id middleware is not applied, request id %s has been generated", reqId)
			}

			logger.Printf("request %s : begin : %s %s -> %s", reqId, r.Method, r.URL.Path, r.RemoteAddr)

			start := time.Now()
			err := nextFn(w, r)
			finish := time.Since(start)

			if err != nil {
				logger.Printf("request %s had error during execution: %v", reqId, err)
			}

			logger.Printf("request %s : end : %s %s -> %s (%s)", reqId, r.Method, r.URL.Path, r.RemoteAddr, finish)
			return err
		}
	}
}
