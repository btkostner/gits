// server/response.go
// Common JSON API responses

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

// Response is a base JSON API response
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []*Error    `json:"errors,omitempty"`

	Meta *Meta `json:"meta,omitempty"`
}

// Meta represents the JSON API meta object
type Meta struct {
	Version string    `json:"version,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

// Error is a simple struct to hold server and controller errors
type Error struct {
	Status int `json:"status,omitempty"`

	Code   int    `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`

	Source *Source `json:"source,omitempty"`
}

// Source defines a place a request error can occure
type Source struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// NewResponse returns a new JSON API response with meta fields
func NewResponse() *Response {
	return &Response{
		Meta: &Meta{
			Version: "2.0.0",
			Date:    time.Now(),
		},
	}
}

// NewError returns a new Error from status and title
func NewError(s int, t string) *Response {
	e := &Error{
		Status: s,
		Title:  t,
	}

	r := NewResponse()
	r.Errors = []*Error{e}
	return r
}

// Handle writes the error to a http response
func (r Response) Handle(w http.ResponseWriter) {
	j, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		logrus.Error(err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Errors != nil && r.Errors[0] != nil {
		w.WriteHeader(r.Errors[0].Status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(j)
}
