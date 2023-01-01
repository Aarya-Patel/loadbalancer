package backend

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(serverUrl string) (*Backend, error) {
	targetURL, err := url.Parse(serverUrl)
	if err != nil {
		return &Backend{}, err
	}

	bknd := Backend{
		URL: targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(generateBackendHandler(targetURL)),
		},
		Status:       Stale,
		ReverseProxy: httputil.NewSingleHostReverseProxy(targetURL),
	}

	return &bknd, nil
}

func InitBackends(backends []*Backend) {
	for _, bknd := range backends {
		server := bknd.Server
		log.Print("Starting backend on ", bknd.URL.String())
		go server.ListenAndServe()
		updateBackendStatus(bknd, Alive)
	}
}

func updateBackendStatus(backend *Backend, newStatus Status) {
	log.Printf("Backend on %s changed status from %s -> %s",
		backend.URL.String(),
		backend.Status.String(),
		newStatus.String(),
	)
	backend.Status = newStatus
}

func generateBackendHandler(url *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, fmt.Sprintf("Hello from %s", url.Host))
	}
}
