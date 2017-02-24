// controller/hook.go
// Handles GitHub hooks

package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"

	"github.com/btkostner/gits/src/config"
	"github.com/btkostner/gits/src/server"
)

func handleHook(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("x-github-event")
	signature := r.Header.Get("x-hub-signature")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Debugf("Unable to read payload: %s", err)
		server.NewError(http.StatusBadRequest, "Unable to read payload").Handle(w)
		return
	}

	payload := github.WebHookPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		logrus.Debugf("Unable to decode payload: %s", err)
		server.NewError(http.StatusBadRequest, "Unable to decode payload").Handle(w)
		return
	}

	if payload.Repo == nil {
		server.NewError(http.StatusBadRequest, "Invalid payload").Handle(w)
		return
	}

	project := config.FindProject(*payload.Repo.FullName)
	if project == nil {
		server.NewError(http.StatusNotFound, "Unknown repository").Handle(w)
		return
	}

	if project.Secret == "" {
		logrus.Warnf("%s is not using a secret", *project.Repo.FullName)
	} else {
		if err := checkSignatures(project.Secret, signature, body); err != nil {
			logrus.Warnf("Invalid signature for %s", project.Repo.FullName)
			server.NewError(http.StatusBadRequest, "Unable to read payload").Handle(w)
			return
		}
	}

	if event == "ping" {
		handlePing(w, r)
		return
	}

	server.NewError(http.StatusBadRequest, "Unknown GitHub event").Handle(w)
	return
}

func checkSignatures(stored string, expected string, body []byte) error {
	mac := hmac.New(sha1.New, []byte(stored))
	mac.Write(body)
	expSig := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expSig), []byte(expected)) {
		return errors.New("Missmatched signatures")
	}

	return nil
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
