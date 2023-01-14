package loadbalancer

import (
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type RoundRobinLoadBalancer struct {
	abstractLoadBalancer

	index int
	mux   sync.Mutex
}

func NewRoundRobinLoadBalancer(name string, serverURL string) (*RoundRobinLoadBalancer, error) {
	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return &RoundRobinLoadBalancer{}, err
	}

	lb := RoundRobinLoadBalancer{
		index:                0,
		mux:                  sync.Mutex{},
		abstractLoadBalancer: abstractLoadBalancer{},
	}

	abstractLB := abstractLoadBalancer{
		Name: name,
		URL:  targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(lb.GenerateLBServeHTTP()),
		},
		ServerURLs: []string{},
		Mapping:    make(map[string]*backend.Backend),
	}

	// Set the abstractLB after we GenerateLBServeHTTP
	lb.abstractLoadBalancer = abstractLB

	return &lb, nil
}

func (lb *RoundRobinLoadBalancer) GenerateLBServeHTTP() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(lb.Mapping) == 0 {
			io.WriteString(rw, "There are no backends available to process your request!")
			return
		}

		lb.mux.Lock()
		defer lb.mux.Unlock()

		wrapAroundIndex := lb.index + len(lb.ServerURLs)
		for ; lb.index < wrapAroundIndex; lb.index++ {
			proxyServerURL := lb.ServerURLs[lb.index%len(lb.Mapping)]
			bknd := lb.Mapping[proxyServerURL]

			if bknd.IsHealthy() {
				bknd.ReverseProxy.ServeHTTP(rw, req)
				lb.index++
				return
			}
		}
		io.WriteString(rw, "There are no healthy backends available to process your request!")
	}
}
