package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

/* Types */
//go:generate stringer -type=Status -output=statusstring.go
type Status uint8

/* Enums */
const (
	Initial Status = iota
	Healthy
	Unhealthy
)

type Backend struct {
	// URL to start the backend server
	URL    *url.URL
	Server *http.Server
	Status Status
	Mux    sync.Mutex

	// ReverseProxy which forwards requests to `Server`
	ReverseProxy *httputil.ReverseProxy
}
