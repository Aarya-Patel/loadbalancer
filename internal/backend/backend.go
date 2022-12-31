package backend

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	// URL to start the backend server
	Url    *url.URL
	Server *http.Server

	// ReverseProxy which forwards requests to `Server`
	ReverseProxy *httputil.ReverseProxy
}

func serveHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hello from sadf")
}

func New(serverUrl string) (*Backend, error) {
	targetURL, err := url.Parse(serverUrl)
	if err != nil {
		return &Backend{}, err
	}

	backend := Backend{
		Url: targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(serveHTTP),
		},
		ReverseProxy: httputil.NewSingleHostReverseProxy(targetURL),
	}

	return &backend, nil
}

func InitBackends(backends []*Backend) {
	for _, backend := range backends {
		server := backend.Server
		log.Print("Starting server on ", backend.Url.String())
		go server.ListenAndServe()
	}
}
