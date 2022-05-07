package middleware

import (
	"net/http"

	"github.com/pkg/errors"
	webErrs "github.com/umalmyha/authsrv/pkg/web/errors"
)

func PanicRecover(nextFn HttpHandlerFn) HttpHandlerFn {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.Wrapf(webErrs.HttpInternalServerErr, "panic occurred on processing http request - %v", rvr)
			}
		}()

		err = nextFn(w, r)
		return err
	}
}
