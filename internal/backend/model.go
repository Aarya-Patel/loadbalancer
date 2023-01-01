package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

/* Types */
//go:generate stringer -type=Status -output=status_string.go
type Status uint8

/* Enums */
const (
	Stale Status = iota
	Alive
	Dead
)

type Backend struct {
	// URL to start the backend server
	URL    *url.URL
	Server *http.Server
	Status Status

	// ReverseProxy which forwards requests to `Server`
	ReverseProxy *httputil.ReverseProxy
}
