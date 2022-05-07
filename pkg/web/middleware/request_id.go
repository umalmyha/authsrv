package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/umalmyha/authsrv/pkg/web/response"
)

type ctxReqIdKey string

const reqIdHeader = "X-Request-Id"
const CtxReqId ctxReqIdKey = "http-request-id"

func RequestId(nextFn HttpHandlerFn) HttpHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		requestId := r.Header.Get(reqIdHeader)
		if requestId == "" {
			requestId = uuid.NewString()
			ctx = context.WithValue(ctx, CtxReqId, requestId)
		}
		defer response.SetHeader(w, reqIdHeader, requestId)

		return nextFn(w, r.WithContext(ctx))
	}
}
