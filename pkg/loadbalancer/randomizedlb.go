package loadbalancer

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

const MaxRetryAttempts int = 5

type RandomizedLoadBalancer struct {
	abstractLoadBalancer
}

func NewRandomizedLoadBalancer(name string, serverURL string) (*RandomizedLoadBalancer, error) {
	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return &RandomizedLoadBalancer{}, err
	}

	lb := RandomizedLoadBalancer{}

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

func (lb *RandomizedLoadBalancer) GenerateLBServeHTTP() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(lb.Mapping) == 0 {
			io.WriteString(rw, "There are no backends available to process your request!")
			return
		}

		rand.Seed(time.Now().UnixNano())
		for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
			index := rand.Intn(len(lb.Mapping))
			proxyServerURL := lb.ServerURLs[index]
			bknd := lb.Mapping[proxyServerURL]

			if bknd.IsHealthy() {
				bknd.ReverseProxy.ServeHTTP(rw, req)
				return
			}
		}
		io.WriteString(rw, "All retry attempts exhausted")
	}
}
