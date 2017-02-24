// controller/controller.go
// Handles route splitting

package controller

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/guregu/kami"
	"github.com/zenazn/goji/web/mutil"
)

// New returns a new server node with all routes
func New() *kami.Mux {
	r := kami.New()

	r.Post("/", handleHook)

	r.LogHandler = func(c context.Context, w mutil.WriterProxy, r *http.Request) {
		logrus.Debugf("%v - %v %v", w.Status(), r.Method, r.URL.Path)
	}

	r.PanicHandler = handleErrorPanic
	r.NotFound(handleErrorFound)
	r.MethodNotAllowed(handleErrorMethod)

	return r
}
