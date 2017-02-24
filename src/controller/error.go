// controller/error.go
// Handles all web server errors

package controller

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/guregu/kami"

	"github.com/btkostner/gits/src/server"
)

func handleErrorPanic(c context.Context, w http.ResponseWriter, r *http.Request) {
	err := kami.Exception(c)

	logrus.WithFields(logrus.Fields{
		"METHOD": r.Method,
		"PATH":   r.URL.Path,
	}).Error(err)

	server.NewError(http.StatusInternalServerError, "Internal Server Error").Handle(w)
}

func handleErrorFound(w http.ResponseWriter, r *http.Request) {
	server.NewError(http.StatusNotFound, "Not Found").Handle(w)
}

func handleErrorMethod(w http.ResponseWriter, r *http.Request) {
	server.NewError(http.StatusMethodNotAllowed, "Method Not Allowed").Handle(w)
}
