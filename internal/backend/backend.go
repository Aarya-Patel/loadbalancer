package backend

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
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
		Status:       Initial,
		Mux:          sync.Mutex{},
		ReverseProxy: httputil.NewSingleHostReverseProxy(targetURL),
	}

	return &bknd, nil
}

func InitBackends(backends []*Backend) {
	for _, bknd := range backends {
		server := bknd.Server
		log.Print("Starting backend on ", bknd.URL.String())
		go server.ListenAndServe()
		bknd.UpdateBackendStatus(Healthy)
	}
}

func (bknd *Backend) IsHealthy() bool {
	bknd.Mux.Lock()
	defer bknd.Mux.Unlock()
	return bknd.Status.String() == Healthy.String()
}

func (bknd *Backend) UpdateBackendStatus(newStatus Status) {
	bknd.Mux.Lock()
	defer bknd.Mux.Unlock()

	if bknd.Status.String() != newStatus.String() {
		log.Printf("Backend on %s changed status from %s -> %s",
			bknd.URL.String(),
			bknd.Status.String(),
			newStatus.String(),
		)
		bknd.Status = newStatus
	}
}

func generateBackendHandler(url *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, fmt.Sprintf("Hello from %s", url.Host))
	}
}
