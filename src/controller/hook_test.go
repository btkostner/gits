// controller/hook_test.go
// Test important controller functions

package controller

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"

	"github.com/btkostner/gits/src/config"
)

// TestHookRead tests web server action on misformated body requests
func TestHookRead(t *testing.T) {
	req := make([]byte, 4)

	mac := hmac.New(sha1.New, []byte("testing"))
	mac.Write(req)
	enc := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	request, err := http.NewRequest("POST", "/", bytes.NewBuffer(req))
	if err != nil {
		t.Error(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-github-event", "ping")
	request.Header.Set("x-hub-signature", enc)

	handler := New()
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 400 {
		t.Error("Did not return 400 status code")
	}
}

// TestHookSignature ensures all signatures are accuratly checked
func TestHookSignature(t *testing.T) {
	viper.Set("projects.btkostner/gits.secret", "testing")
	config.ReadProjects()

	repoName := "btkostner/gits"

	body := github.WebHookPayload{
		Repo: &github.Repository{
			FullName: &repoName,
		},
	}
	req, err := json.Marshal(body)
	if err != nil {
		t.Error(err)
		return
	}

	mac := hmac.New(sha1.New, []byte("testing"))
	mac.Write(req)
	enc := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	request, err := http.NewRequest("POST", "/", bytes.NewBuffer(req))
	if err != nil {
		t.Error(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-github-event", "ping")
	request.Header.Set("x-hub-signature", enc)

	handler := New()
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Error("Did not return 200 status code")
	}
}

// TestCheckSignatures ensures we acuratly are testing signatures
func TestCheckSignatures(t *testing.T) {
	sig := "testing"
	str := []byte("body")

	mac := hmac.New(sha1.New, []byte(sig))
	mac.Write(str)
	enc := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	err := checkSignatures(sig, enc, str)

	if err != nil {
		t.Error(err)
	}
}

// BenchmarkHookPing benchmarks basic PING requests
func BenchmarkHookPing(b *testing.B) {
	body := github.WebHookPayload{}
	req, err := json.Marshal(body)
	if err != nil {
		b.Error(err)
	}

	handler := New()
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "/", bytes.NewBuffer(req))
	if err != nil {
		b.Error(err)
	}

	request.Header.Set("x-github-event", "ping")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(recorder, request)
	}
}
