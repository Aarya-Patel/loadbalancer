package loadbalancer

import (
	"io"
	"net/http"
	"net/url"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type RoundRobinLoadBalancer struct {
	*abstractLoadBalancer
	// Maintain a consistent order of server URL.String()'s
	ServerURLs []string
	index      int
}

func NewRoundRobinLoadBalancer(name string, serverURL string) (*RoundRobinLoadBalancer, error) {
	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return &RoundRobinLoadBalancer{}, err
	}

	lb := RoundRobinLoadBalancer{
		index:                0,
		ServerURLs:           []string{},
		abstractLoadBalancer: nil,
	}

	abstractLB := abstractLoadBalancer{
		Name: name,
		URL:  targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(lb.GenerateLBServeHTTP()),
		},
		Mapping: make(map[string]*backend.Backend),
	}

	// Set the abstractLB after we GenerateLBServeHTTP
	lb.abstractLoadBalancer = &abstractLB

	return &lb, nil
}

func (lb *RoundRobinLoadBalancer) GenerateLBServeHTTP() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(lb.Mapping) == 0 {
			io.WriteString(rw, "There are no backends available to process your request!")
			return
		}

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

func (lb *RoundRobinLoadBalancer) InsertBackend(serverURL string) error {
	err := lb.abstractLoadBalancer.InsertBackend(serverURL)
	if err != nil {
		return err
	}

	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return err
	}

	lb.ServerURLs = append(lb.ServerURLs, targetURL.String())
	return nil
}
