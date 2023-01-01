package loadbalancer

import (
	"io"
	"net/http"
	"net/url"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type RoundRobinLoadBalancer struct {
	abstractLoadBalancer
	index int
}

func NewRoundRobinLoadBalancer(name string, serverURL string) (*RoundRobinLoadBalancer, error) {
	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return &RoundRobinLoadBalancer{}, err
	}

	lb := RoundRobinLoadBalancer{
		index:                0,
		abstractLoadBalancer: abstractLoadBalancer{},
	}

	abstractLB := abstractLoadBalancer{
		Name: name,
		URL:  targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(lb.GenerateLBServeHTTP()),
		},
		Mapping: make(map[*url.URL]*backend.Backend),
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

		serverURLs := []*url.URL{}
		for url := range lb.Mapping {
			serverURLs = append(serverURLs, url)
		}

		proxyServerURL := serverURLs[lb.index%len(lb.Mapping)]
		bknd := lb.Mapping[proxyServerURL]
		bknd.ReverseProxy.ServeHTTP(rw, req)

		lb.index++
	}
}
